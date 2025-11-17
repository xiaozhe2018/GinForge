package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 配置管理器
type Config struct {
	viper *viper.Viper
	env   string
}

// New 创建配置管理器
// 配置组织方式：按环境目录组织
// - 基础配置：configs/{env}/base.yaml
// - 服务配置：configs/{env}/{service}.yaml
//
// 配置优先级（从高到低）：
// 1. 环境变量（GOEASE_* 前缀）
// 2. .env 文件中的环境变量
// 3. YAML 配置文件中的值
// 4. 代码中的默认值
func New() *Config {
	// 自动加载 .env 文件（如果存在）
	// 忽略错误，允许不提供 .env 文件
	_ = godotenv.Load()

	v := viper.New()

	// 设置环境变量前缀
	v.SetEnvPrefix("GOEASE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 设置默认值
	setDefaults(v)

	// 优先从环境变量获取环境名称，如果没有则使用默认值
	env := os.Getenv("GOEASE_APP_ENV")
	if env == "" {
		env = "dev"
	}

	// 获取服务名称（可选）
	serviceName := os.Getenv("SERVICE_NAME")
	// 如果未设置，尝试从程序名推断（例如：admin-api -> admin-api）
	if serviceName == "" {
		if programName := os.Args[0]; programName != "" {
			// 从程序路径中提取服务名
			baseName := filepath.Base(programName)
			// 移除可能的扩展名
			serviceName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		}
	}

	// 按环境目录组织方式加载配置
	envDir := filepath.Join("./configs", env)
	baseConfigPath := filepath.Join(envDir, "base.yaml")

	// 检查环境目录是否存在
	if _, err := os.Stat(baseConfigPath); err != nil {
		panic(fmt.Errorf("config file not found: %s, please set GOEASE_APP_ENV environment variable (dev/test/prod)", baseConfigPath))
	}

	// 1. 读取基础配置（base.yaml 必须存在且最先加载）
	v.SetConfigName("base")
	v.SetConfigType("yaml")
	v.AddConfigPath(envDir)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading base config: %w", err))
	}

	// 2. 自动加载环境目录下的所有其他 YAML 配置文件
	// 这样添加新配置文件时，无需修改代码，只需创建文件即可
	if err := loadConfigFiles(v, envDir, serviceName); err != nil {
		panic(fmt.Errorf("fatal error loading config files: %w", err))
	}

	// 更新环境变量（从配置文件读取）
	if v.IsSet("app.env") {
		env = v.GetString("app.env")
	}

	return &Config{
		viper: v,
		env:   env,
	}
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	// 应用配置
	v.SetDefault("app.name", "GinForge Framework")
	v.SetDefault("app.version", "0.1.0")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.debug", true)
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.read_timeout", "10s")
	v.SetDefault("app.write_timeout", "10s")
	v.SetDefault("app.idle_timeout", "60s")

	// 日志配置
	v.SetDefault("log.level", "debug")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("log.file_path", "logs/app.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.max_backups", 10)
	v.SetDefault("log.compress", true)

	// 数据库配置
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.database", "goweb.db")
	v.SetDefault("database.username", "")
	v.SetDefault("database.password", "")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.timezone", "Asia/Shanghai")
	v.SetDefault("database.table_prefix", "gf_") // 默认表前缀
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.log_level", "warn")

	// Redis配置
	v.SetDefault("redis.enabled", false)
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.database", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")
	v.SetDefault("redis.idle_timeout", "5m")
	v.SetDefault("redis.idle_check_frequency", "1m")

	// 缓存配置
	v.SetDefault("cache.default_ttl", "5m")
	v.SetDefault("cache.max_size", 1000)
	v.SetDefault("cache.cleanup_interval", "10m")

	// JWT配置
	v.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	v.SetDefault("jwt.expire", "24h")
	v.SetDefault("jwt.issuer", "GinForge")
	v.SetDefault("jwt.audience", "GinForge-Users")

	// 限流配置
	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.rps", 100)
	v.SetDefault("rate_limit.burst", 200)
	v.SetDefault("rate_limit.window", "1m")

	// CORS配置
	v.SetDefault("cors.enabled", true)
	v.SetDefault("cors.origins", []string{"*"})
	v.SetDefault("cors.methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.headers", []string{"Content-Type", "Authorization", "X-Requested-With"})
	v.SetDefault("cors.credentials", true)
	v.SetDefault("cors.max_age", "12h")

	// 监控配置
	v.SetDefault("monitor.enabled", true)
	v.SetDefault("monitor.metrics_path", "/metrics")
	v.SetDefault("monitor.health_path", "/healthz")
	v.SetDefault("monitor.ready_path", "/readyz")

	// Gateway 配置
	v.SetDefault("gateway.base_url", "http://localhost:8080")
	v.SetDefault("gateway.timeout", "30s")
	v.SetDefault("gateway.retry_count", 3)
	v.SetDefault("gateway.retry_delay", "1s")
}

// loadConfigFiles 自动加载环境目录下的所有 YAML 配置文件
// 加载顺序：
// 1. base.yaml（已加载）
// 2. 其他 .yaml 文件按文件名排序（确保一致性）
//
// 这样添加新配置文件时，只需创建文件即可，无需修改代码
func loadConfigFiles(v *viper.Viper, envDir string, serviceName string) error {
	// 读取目录下的所有文件
	entries, err := os.ReadDir(envDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	// 收集所有 .yaml 文件（排除已加载的 base.yaml）
	var configFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// 跳过 base.yaml（已加载）
		if name == "base.yaml" {
			continue
		}

		// 只处理 .yaml 文件
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			configFiles = append(configFiles, name)
		}
	}

	// 按文件名排序，确保加载顺序一致
	// 这允许通过文件名前缀控制加载顺序，例如：01-cache.yaml, 02-service.yaml
	sort.Strings(configFiles)

	// 如果指定了服务名称，优先加载对应的服务配置文件
	if serviceName != "" {
		serviceFileName := fmt.Sprintf("%s.yaml", serviceName)
		// 查找服务配置文件的位置
		for i, file := range configFiles {
			if file == serviceFileName {
				// 将服务配置文件移到最前面（紧跟在 base.yaml 之后）
				configFiles = append([]string{serviceFileName}, append(configFiles[:i], configFiles[i+1:]...)...)
				break
			}
		}
	}

	// 加载所有配置文件（按顺序合并）
	for _, fileName := range configFiles {
		configName := strings.TrimSuffix(fileName, ".yaml")
		configName = strings.TrimSuffix(configName, ".yml")

		v.SetConfigName(configName)
		if err := v.MergeInConfig(); err != nil {
			// 配置文件不存在或格式错误时记录警告，但不中断
			// 这允许配置文件是可选的
			continue
		}
	}

	return nil
}

// GetString 获取字符串配置
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetInt64 获取64位整数配置
func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetDuration 获取时间配置
func (c *Config) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

// GetStringSlice 获取字符串切片配置
func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

// Get 获取任意类型配置
func (c *Config) Get(key string) interface{} {
	return c.viper.Get(key)
}

// Set 设置配置
func (c *Config) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

// IsSet 检查配置是否设置
func (c *Config) IsSet(key string) bool {
	return c.viper.IsSet(key)
}

// AllSettings 获取所有配置
func (c *Config) AllSettings() map[string]interface{} {
	return c.viper.AllSettings()
}

// GetEnv 获取环境
func (c *Config) GetEnv() string {
	return c.env
}

// IsDevelopment 是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.env == "development"
}

// IsProduction 是否为生产环境
func (c *Config) IsProduction() bool {
	return c.env == "production"
}

// IsTest 是否为测试环境
func (c *Config) IsTest() bool {
	return c.env == "test"
}

// GetConfigFile 获取配置文件路径
func (c *Config) GetConfigFile() string {
	return c.viper.ConfigFileUsed()
}

// WatchConfig 监听配置文件变化
func (c *Config) WatchConfig() {
	c.viper.WatchConfig()
}

// OnConfigChange 配置文件变化回调
func (c *Config) OnConfigChange(fn func()) {
	c.viper.OnConfigChange(func(e fsnotify.Event) {
		fn()
	})
}

// SaveConfig 保存配置到文件
func (c *Config) SaveConfig() error {
	configFile := c.GetConfigFile()
	if configFile == "" {
		configFile = "configs/config.yaml"
	}

	// 确保目录存在
	dir := filepath.Dir(configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return c.viper.WriteConfigAs(configFile)
}

// Unmarshal 解析配置到结构体
func (c *Config) Unmarshal(key string, rawVal interface{}) error {
	return c.viper.UnmarshalKey(key, rawVal)
}

// UnmarshalAll 解析所有配置到结构体
func (c *Config) UnmarshalAll(rawVal interface{}) error {
	return c.viper.Unmarshal(rawVal)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `yaml:"driver" json:"driver"`
	Host            string        `yaml:"host" json:"host"`
	Port            int           `yaml:"port" json:"port"`
	Database        string        `yaml:"database" json:"database"`
	Username        string        `yaml:"username" json:"username"`
	Password        string        `yaml:"password" json:"password"`
	Charset         string        `yaml:"charset" json:"charset"`
	Timezone        string        `yaml:"timezone" json:"timezone"`
	TablePrefix     string        `yaml:"table_prefix" json:"table_prefix"` // 数据表前缀
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
	LogLevel        string        `yaml:"log_level" json:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Enabled            bool          `yaml:"enabled" json:"enabled"`
	Host               string        `yaml:"host" json:"host"`
	Port               int           `yaml:"port" json:"port"`
	Password           string        `yaml:"password" json:"password"`
	Database           int           `yaml:"database" json:"database"`
	PoolSize           int           `yaml:"pool_size" json:"pool_size"`
	MinIdleConns       int           `yaml:"min_idle_conns" json:"min_idle_conns"`
	MaxRetries         int           `yaml:"max_retries" json:"max_retries"`
	DialTimeout        time.Duration `yaml:"dial_timeout" json:"dial_timeout"`
	ReadTimeout        time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout" json:"write_timeout"`
	IdleTimeout        time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
	IdleCheckFrequency time.Duration `yaml:"idle_check_frequency" json:"idle_check_frequency"`
}

// GetDatabaseConfig 获取数据库配置
func (c *Config) GetDatabaseConfig() DatabaseConfig {
	var config DatabaseConfig
	c.Unmarshal("database", &config)

	// 如果表前缀为空，使用默认值
	if config.TablePrefix == "" {
		config.TablePrefix = c.GetString("database.table_prefix")
		if config.TablePrefix == "" {
			config.TablePrefix = "gf_" // 最终默认值
		}
	}

	return config
}

// GetTablePrefix 获取数据表前缀
func (c *Config) GetTablePrefix() string {
	prefix := c.GetString("database.table_prefix")
	if prefix == "" {
		return "gf_" // 默认表前缀
	}
	return prefix
}

// GetRedisConfig 获取Redis配置
func (c *Config) GetRedisConfig() RedisConfig {
	var config RedisConfig
	c.Unmarshal("redis", &config)
	return config
}
