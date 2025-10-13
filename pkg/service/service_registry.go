package service

import (
	"context"
	"fmt"
	"goweb/pkg/config"
	"goweb/pkg/gateway"
	"goweb/pkg/logger"
	"sync"
)

// ServiceRegistry 服务注册表
type ServiceRegistry struct {
	services map[string]ServiceInfo
	mutex    sync.RWMutex
	gateway  *gateway.Client
	logger   logger.Logger
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Version string `json:"version"`
	Status  string `json:"status"` // active, inactive, maintenance
}

// NewServiceRegistry 创建服务注册表
func NewServiceRegistry(cfg *config.Config, log logger.Logger) *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]ServiceInfo),
		gateway:  gateway.NewClient(cfg, log),
		logger:   log,
	}
}

// Register 注册服务
func (sr *ServiceRegistry) Register(service ServiceInfo) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	sr.services[service.Name] = service
	sr.logger.Info("service registered", "name", service.Name, "host", service.Host, "port", service.Port)
}

// Unregister 注销服务
func (sr *ServiceRegistry) Unregister(serviceName string) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	delete(sr.services, serviceName)
	sr.logger.Info("service unregistered", "name", serviceName)
}

// GetService 获取服务信息
func (sr *ServiceRegistry) GetService(serviceName string) (ServiceInfo, bool) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	service, exists := sr.services[serviceName]
	return service, exists
}

// ListServices 列出所有服务
func (sr *ServiceRegistry) ListServices() map[string]ServiceInfo {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	services := make(map[string]ServiceInfo)
	for k, v := range sr.services {
		services[k] = v
	}
	return services
}

// CallService 调用服务
func (sr *ServiceRegistry) CallService(ctx context.Context, serviceName, method, path string, body interface{}) (*gateway.Response, error) {
	service, exists := sr.GetService(serviceName)
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	if service.Status != "active" {
		return nil, fmt.Errorf("service %s is not active", serviceName)
	}

	// 构建完整路径
	fullPath := fmt.Sprintf("/%s%s", serviceName, path)

	// 通过 Gateway 调用服务
	switch method {
	case "GET":
		return sr.gateway.Get(ctx, fullPath, nil, nil)
	case "POST":
		return sr.gateway.Post(ctx, fullPath, body, nil)
	case "PUT":
		return sr.gateway.Put(ctx, fullPath, body, nil)
	case "DELETE":
		return sr.gateway.Delete(ctx, fullPath, nil)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

// HealthCheck 健康检查
func (sr *ServiceRegistry) HealthCheck(ctx context.Context, serviceName string) error {
	_, err := sr.CallService(ctx, serviceName, "GET", "/healthz", nil)
	return err
}

// ServiceClient 服务客户端
type ServiceClient struct {
	registry    *ServiceRegistry
	serviceName string
}

// NewServiceClient 创建服务客户端
func NewServiceClient(registry *ServiceRegistry, serviceName string) *ServiceClient {
	return &ServiceClient{
		registry:    registry,
		serviceName: serviceName,
	}
}

// Call 调用服务方法
func (sc *ServiceClient) Call(ctx context.Context, method, path string, body interface{}) (*gateway.Response, error) {
	return sc.registry.CallService(ctx, sc.serviceName, method, path, body)
}

// Get 发送GET请求
func (sc *ServiceClient) Get(ctx context.Context, path string) (*gateway.Response, error) {
	return sc.Call(ctx, "GET", path, nil)
}

// Post 发送POST请求
func (sc *ServiceClient) Post(ctx context.Context, path string, body interface{}) (*gateway.Response, error) {
	return sc.Call(ctx, "POST", path, body)
}

// Put 发送PUT请求
func (sc *ServiceClient) Put(ctx context.Context, path string, body interface{}) (*gateway.Response, error) {
	return sc.Call(ctx, "PUT", path, body)
}

// Delete 发送DELETE请求
func (sc *ServiceClient) Delete(ctx context.Context, path string) (*gateway.Response, error) {
	return sc.Call(ctx, "DELETE", path, nil)
}
