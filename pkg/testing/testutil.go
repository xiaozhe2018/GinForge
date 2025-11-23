package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/redis"
)

// TestConfig 测试配置
type TestConfig struct {
	*config.Config
	TestDB    string
	TestRedis string
}

// NewTestConfig 创建测试配置
func NewTestConfig() *TestConfig {
	cfg := config.New()

	// 设置测试环境
	cfg.Set("app.env", "test")
	cfg.Set("log.level", "debug")
	cfg.Set("log.format", "console")

	// 测试数据库配置
	cfg.Set("database.driver", "sqlite")
	cfg.Set("database.database", ":memory:")
	cfg.Set("database.log_level", "silent")

	// 测试Redis配置
	cfg.Set("redis.enabled", true)
	cfg.Set("redis.host", "localhost")
	cfg.Set("redis.port", 6379)
	cfg.Set("redis.database", 15) // 使用测试数据库

	return &TestConfig{
		Config:    cfg,
		TestDB:    ":memory:",
		TestRedis: "localhost:6379",
	}
}

// TestServer 测试服务器
type TestServer struct {
	*httptest.Server
	Config *TestConfig
	Logger logger.Logger
	Redis  *redis.Manager
	URL    string
	T      *testing.T
}

// NewTestServer 创建测试服务器
func NewTestServer(t *testing.T, setupRoutes func(*gin.Engine)) *TestServer {
	cfg := NewTestConfig()
	log := logger.New("test", "debug", "stdout", "")

	// 初始化Redis
	redisCfg := &config.RedisConfig{
		Host:               cfg.GetString("redis.host"),
		Port:               cfg.GetInt("redis.port"),
		Password:           cfg.GetString("redis.password"),
		Database:           cfg.GetInt("redis.database"),
		PoolSize:           cfg.GetInt("redis.pool_size"),
		MinIdleConns:       cfg.GetInt("redis.min_idle_conns"),
		MaxRetries:         cfg.GetInt("redis.max_retries"),
		DialTimeout:        cfg.GetDuration("redis.dial_timeout"),
		ReadTimeout:        cfg.GetDuration("redis.read_timeout"),
		WriteTimeout:       cfg.GetDuration("redis.write_timeout"),
		IdleTimeout:        cfg.GetDuration("redis.idle_timeout"),
		IdleCheckFrequency: cfg.GetDuration("redis.idle_check_frequency"),
	}
	redisMgr := redis.NewManager(redisCfg, log)

	// 创建Gin引擎
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 设置路由
	setupRoutes(r)

	// 创建测试服务器
	server := httptest.NewServer(r)

	return &TestServer{
		Server: server,
		Config: cfg,
		Logger: log,
		Redis:  redisMgr,
		URL:    server.URL,
		T:      t,
	}
}

// Close 关闭测试服务器
func (ts *TestServer) Close() {
	ts.Server.Close()
	if ts.Redis != nil {
		ts.Redis.Close()
	}
}

// Get 发送GET请求
func (ts *TestServer) Get(path string) *TestResponse {
	return ts.Request("GET", path, nil)
}

// Post 发送POST请求
func (ts *TestServer) Post(path string, body interface{}) *TestResponse {
	return ts.Request("POST", path, body)
}

// Put 发送PUT请求
func (ts *TestServer) Put(path string, body interface{}) *TestResponse {
	return ts.Request("PUT", path, body)
}

// Delete 发送DELETE请求
func (ts *TestServer) Delete(path string) *TestResponse {
	return ts.Request("DELETE", path, nil)
}

// Request 发送HTTP请求
func (ts *TestServer) Request(method, path string, body interface{}) *TestResponse {
	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(method, ts.URL+path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(ts.T, err)

	return &TestResponse{Response: resp}
}

// TestResponse 测试响应
type TestResponse struct {
	*http.Response
}

// AssertStatus 断言状态码
func (tr *TestResponse) AssertStatus(t *testing.T, expected int) *TestResponse {
	assert.Equal(t, expected, tr.Response.StatusCode)
	return tr
}

// AssertJSON 断言JSON响应
func (tr *TestResponse) AssertJSON(t *testing.T, expected interface{}) *TestResponse {
	var actual interface{}
	err := json.NewDecoder(tr.Response.Body).Decode(&actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
	return tr
}

// AssertContains 断言响应包含特定内容
func (tr *TestResponse) AssertContains(t *testing.T, substring string) *TestResponse {
	body, _ := io.ReadAll(tr.Response.Body)
	assert.Contains(t, string(body), substring)
	return tr
}

// TestDatabase 测试数据库
type TestDatabase struct {
	*config.Config
	Logger logger.Logger
}

// NewTestDatabase 创建测试数据库
func NewTestDatabase(t *testing.T) *TestDatabase {
	cfg := NewTestConfig()
	log := logger.New("test-db", "debug", "stdout", "")

	return &TestDatabase{
		Config: cfg.Config,
		Logger: log,
	}
}

// TestRedis 测试Redis
type TestRedis struct {
	*redis.Manager
	Logger logger.Logger
}

// NewTestRedis 创建测试Redis
func NewTestRedis(t *testing.T) *TestRedis {
	cfg := NewTestConfig()
	log := logger.New("test-redis", "debug", "stdout", "")

	redisCfg := &config.RedisConfig{
		Host:               cfg.GetString("redis.host"),
		Port:               cfg.GetInt("redis.port"),
		Password:           cfg.GetString("redis.password"),
		Database:           cfg.GetInt("redis.database"),
		PoolSize:           cfg.GetInt("redis.pool_size"),
		MinIdleConns:       cfg.GetInt("redis.min_idle_conns"),
		MaxRetries:         cfg.GetInt("redis.max_retries"),
		DialTimeout:        cfg.GetDuration("redis.dial_timeout"),
		ReadTimeout:        cfg.GetDuration("redis.read_timeout"),
		WriteTimeout:       cfg.GetDuration("redis.write_timeout"),
		IdleTimeout:        cfg.GetDuration("redis.idle_timeout"),
		IdleCheckFrequency: cfg.GetDuration("redis.idle_check_frequency"),
	}
	redisMgr := redis.NewManager(redisCfg, log)

	return &TestRedis{
		Manager: redisMgr,
		Logger:  log,
	}
}

// Cleanup 清理测试数据
func (tr *TestRedis) Cleanup() {
	tr.Manager.GetCache().Clear(context.Background())
}

// TestSuite 测试套件
type TestSuite struct {
	*testing.T
	Config *TestConfig
	Logger logger.Logger
	Redis  *TestRedis
}

// NewTestSuite 创建测试套件
func NewTestSuite(t *testing.T) *TestSuite {
	cfg := NewTestConfig()
	log := logger.New("test-suite", "debug", "stdout", "")

	redisCfg := &config.RedisConfig{
		Host:               cfg.GetString("redis.host"),
		Port:               cfg.GetInt("redis.port"),
		Password:           cfg.GetString("redis.password"),
		Database:           cfg.GetInt("redis.database"),
		PoolSize:           cfg.GetInt("redis.pool_size"),
		MinIdleConns:       cfg.GetInt("redis.min_idle_conns"),
		MaxRetries:         cfg.GetInt("redis.max_retries"),
		DialTimeout:        cfg.GetDuration("redis.dial_timeout"),
		ReadTimeout:        cfg.GetDuration("redis.read_timeout"),
		WriteTimeout:       cfg.GetDuration("redis.write_timeout"),
		IdleTimeout:        cfg.GetDuration("redis.idle_timeout"),
		IdleCheckFrequency: cfg.GetDuration("redis.idle_check_frequency"),
	}
	redisMgr := redis.NewManager(redisCfg, log)

	return &TestSuite{
		T:      t,
		Config: cfg,
		Logger: log,
		Redis: &TestRedis{
			Manager: redisMgr,
			Logger:  log,
		},
	}
}

// Setup 设置测试环境
func (ts *TestSuite) Setup() {
	// 清理Redis
	ts.Redis.Cleanup()
}

// Teardown 清理测试环境
func (ts *TestSuite) Teardown() {
	ts.Redis.Cleanup()
	ts.Redis.GetClient().Close()
}

// TestCase 测试用例
type TestCase struct {
	Name     string
	Setup    func()
	Test     func()
	Teardown func()
}

// RunTestCase 运行测试用例
func (ts *TestSuite) RunTestCase(tc TestCase) {
	ts.T.Run(tc.Name, func(t *testing.T) {
		if tc.Setup != nil {
			tc.Setup()
		}

		if tc.Teardown != nil {
			defer tc.Teardown()
		}

		tc.Test()
	})
}

// MockHTTPClient 模拟HTTP客户端
type MockHTTPClient struct {
	Responses map[string]*http.Response
}

// NewMockHTTPClient 创建模拟HTTP客户端
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		Responses: make(map[string]*http.Response),
	}
}

// SetResponse 设置模拟响应
func (m *MockHTTPClient) SetResponse(url string, resp *http.Response) {
	m.Responses[url] = resp
}

// Do 模拟HTTP请求
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if resp, exists := m.Responses[req.URL.String()]; exists {
		return resp, nil
	}

	// 返回默认404响应
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewBufferString("Not Found")),
	}, nil
}

// TestHelper 测试辅助函数
type TestHelper struct {
	*testing.T
}

// NewTestHelper 创建测试辅助函数
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{T: t}
}

// AssertNoError 断言无错误
func (th *TestHelper) AssertNoError(err error) {
	require.NoError(th.T, err)
}

// AssertEqual 断言相等
func (th *TestHelper) AssertEqual(expected, actual interface{}) {
	assert.Equal(th.T, expected, actual)
}

// AssertTrue 断言为真
func (th *TestHelper) AssertTrue(condition bool) {
	assert.True(th.T, condition)
}

// AssertFalse 断言为假
func (th *TestHelper) AssertFalse(condition bool) {
	assert.False(th.T, condition)
}

// AssertContains 断言包含
func (th *TestHelper) AssertContains(s, substr string) {
	assert.Contains(th.T, s, substr)
}

// AssertNotNil 断言非空
func (th *TestHelper) AssertNotNil(obj interface{}) {
	assert.NotNil(th.T, obj)
}

// AssertNil 断言为空
func (th *TestHelper) AssertNil(obj interface{}) {
	assert.Nil(th.T, obj)
}

// WaitFor 等待条件满足
func (th *TestHelper) WaitFor(condition func() bool, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			th.T.Fatalf("等待超时: %v", timeout)
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}

// TestData 测试数据
type TestData struct {
	Users    []map[string]interface{}
	Products []map[string]interface{}
	Orders   []map[string]interface{}
}

// NewTestData 创建测试数据
func NewTestData() *TestData {
	return &TestData{
		Users: []map[string]interface{}{
			{"id": 1, "name": "测试用户1", "email": "user1@test.com"},
			{"id": 2, "name": "测试用户2", "email": "user2@test.com"},
		},
		Products: []map[string]interface{}{
			{"id": 1, "name": "测试商品1", "price": 100.0},
			{"id": 2, "name": "测试商品2", "price": 200.0},
		},
		Orders: []map[string]interface{}{
			{"id": 1, "user_id": 1, "product_id": 1, "quantity": 2},
			{"id": 2, "user_id": 2, "product_id": 2, "quantity": 1},
		},
	}
}

// LoadTestData 加载测试数据到Redis
func (td *TestData) LoadTestData(redisMgr *redis.Manager) error {
	ctx := context.Background()

	// 加载用户数据
	for _, user := range td.Users {
		key := fmt.Sprintf("test:user:%v", user["id"])
		data, _ := json.Marshal(user)
		if err := redisMgr.GetCache().Set(ctx, key, data, time.Hour); err != nil {
			return err
		}
	}

	// 加载商品数据
	for _, product := range td.Products {
		key := fmt.Sprintf("test:product:%v", product["id"])
		data, _ := json.Marshal(product)
		if err := redisMgr.GetCache().Set(ctx, key, data, time.Hour); err != nil {
			return err
		}
	}

	// 加载订单数据
	for _, order := range td.Orders {
		key := fmt.Sprintf("test:order:%v", order["id"])
		data, _ := json.Marshal(order)
		if err := redisMgr.GetCache().Set(ctx, key, data, time.Hour); err != nil {
			return err
		}
	}

	return nil
}

// CleanTestData 清理测试数据
func (td *TestData) CleanTestData(redisMgr *redis.Manager) error {
	ctx := context.Background()

	// 清理所有测试数据
	pattern := "test:*"
	keys, err := redisMgr.GetClient().GetClient().Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return redisMgr.GetClient().GetClient().Del(ctx, keys...).Err()
	}

	return nil
}
