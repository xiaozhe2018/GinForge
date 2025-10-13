package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"goweb/pkg/config"
	"goweb/pkg/logger"
)

// Client Gateway 客户端
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     logger.Logger
	timeout    time.Duration
}

// NewClient 创建 Gateway 客户端
func NewClient(cfg *config.Config, log logger.Logger) *Client {
	timeout := cfg.GetDuration("gateway.timeout")
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		baseURL: cfg.GetString("gateway.base_url"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger:  log,
		timeout: timeout,
	}
}

// Request 请求结构
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    interface{}
	Query   map[string]string
}

// Response 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"trace_id"`
}

// Call 调用 Gateway 服务
func (c *Client) Call(ctx context.Context, req *Request) (*Response, error) {
	// 构建完整URL
	url := c.baseURL + req.Path
	if req.Query != nil && len(req.Query) > 0 {
		url += "?"
		first := true
		for k, v := range req.Query {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}

	// 准备请求体
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// 添加追踪ID
	if traceID := ctx.Value("trace_id"); traceID != nil {
		httpReq.Header.Set("X-Request-Id", traceID.(string))
	}

	// 添加自定义请求头
	if req.Headers != nil {
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}
	}

	// 发送请求
	c.logger.Debug("calling gateway", "method", req.Method, "url", url)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// 解析响应
	var response Response
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode >= 400 {
		return &response, fmt.Errorf("http error: %d, message: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, path string, query map[string]string, headers map[string]string) (*Response, error) {
	return c.Call(ctx, &Request{
		Method:  "GET",
		Path:    path,
		Query:   query,
		Headers: headers,
	})
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Call(ctx, &Request{
		Method:  "POST",
		Path:    path,
		Body:    body,
		Headers: headers,
	})
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Call(ctx, &Request{
		Method:  "PUT",
		Path:    path,
		Body:    body,
		Headers: headers,
	})
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, path string, headers map[string]string) (*Response, error) {
	return c.Call(ctx, &Request{
		Method:  "DELETE",
		Path:    path,
		Headers: headers,
	})
}

// 业务方法封装

// GetUser 获取用户信息
func (c *Client) GetUser(ctx context.Context, userID string) (*Response, error) {
	return c.Get(ctx, fmt.Sprintf("/api/v1/user/%s", userID), nil, nil)
}

// CreateUser 创建用户
func (c *Client) CreateUser(ctx context.Context, userData map[string]interface{}) (*Response, error) {
	return c.Post(ctx, "/api/v1/user", userData, nil)
}

// UpdateUser 更新用户
func (c *Client) UpdateUser(ctx context.Context, userID string, userData map[string]interface{}) (*Response, error) {
	return c.Put(ctx, fmt.Sprintf("/api/v1/user/%s", userID), userData, nil)
}

// GetMerchant 获取商户信息
func (c *Client) GetMerchant(ctx context.Context, merchantID string) (*Response, error) {
	return c.Get(ctx, fmt.Sprintf("/api/v1/merchant/%s", merchantID), nil, nil)
}

// GetProduct 获取商品信息
func (c *Client) GetProduct(ctx context.Context, productID string) (*Response, error) {
	return c.Get(ctx, fmt.Sprintf("/api/v1/product/%s", productID), nil, nil)
}

// GetOrder 获取订单信息
func (c *Client) GetOrder(ctx context.Context, orderID string) (*Response, error) {
	return c.Get(ctx, fmt.Sprintf("/api/v1/order/%s", orderID), nil, nil)
}

// CreateOrder 创建订单
func (c *Client) CreateOrder(ctx context.Context, orderData map[string]interface{}) (*Response, error) {
	return c.Post(ctx, "/api/v1/order", orderData, nil)
}
