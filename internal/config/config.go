package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	RocketMQ  RocketMQConfig  `mapstructure:"rocketmq"`
	Auth      AuthConfig      `mapstructure:"auth"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Metrics   MetricsConfig   `mapstructure:"metrics"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Name, d.Charset)
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type RocketMQConfig struct {
	NameServer    string `mapstructure:"nameserver"`
	GroupID       string `mapstructure:"group_id"`
	Topic         string `mapstructure:"topic"`
	InstanceName  string `mapstructure:"instance_name"`
}

type AuthConfig struct {
	SessionKey       string `mapstructure:"session_key"`
	SessionExpireHours int   `mapstructure:"session_expire_hours"`
	CookieName       string `mapstructure:"cookie_name"`
}

type RateLimitConfig struct {
	QPS   int `mapstructure:"qps"`
	Burst int `mapstructure:"burst"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	Path   string `mapstructure:"path"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

var cfg *Config

// Init 初始化配置
func Init(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Get 获取全局配置
func Get() *Config {
	return cfg
}

// GetAppConfig 获取应用配置
func GetAppConfig() *AppConfig {
	return &cfg.App
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() *DatabaseConfig {
	return &cfg.Database
}

// GetRedisConfig 获取 Redis 配置
func GetRedisConfig() *RedisConfig {
	return &cfg.Redis
}

// GetRocketMQConfig 获取 RocketMQ 配置
func GetRocketMQConfig() *RocketMQConfig {
	return &cfg.RocketMQ
}

// GetAuthConfig 获取认证配置
func GetAuthConfig() *AuthConfig {
	return &cfg.Auth
}

// GetRateLimitConfig 获取限流配置
func GetRateLimitConfig() *RateLimitConfig {
	return &cfg.RateLimit
}

// GetMetricsConfig 获取监控配置
func GetMetricsConfig() *MetricsConfig {
	return &cfg.Metrics
}

// GetLoggingConfig 获取日志配置
func GetLoggingConfig() *LoggingConfig {
	return &cfg.Logging
}
