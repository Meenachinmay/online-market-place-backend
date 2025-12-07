package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	Tracing  TracingConfig  `mapstructure:"tracing"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`
	Debug       bool   `mapstructure:"debug"`
}

type ServerConfig struct {
	GRPC    GRPCConfig    `mapstructure:"grpc"`
	Timeout TimeoutConfig `mapstructure:"timeout"`
}

type GRPCConfig struct {
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	MaxConnectionIdle     time.Duration `mapstructure:"max_connection_idle"`
	MaxConnectionAge      time.Duration `mapstructure:"max_connection_age"`
	MaxConnectionAgeGrace time.Duration `mapstructure:"max_connection_age_grace"`
	KeepAliveTime         time.Duration `mapstructure:"keep_alive_time"`
	KeepAliveTimeout      time.Duration `mapstructure:"keep_alive_timeout"`
}

type TimeoutConfig struct {
	Read    time.Duration `mapstructure:"read"`
	Write   time.Duration `mapstructure:"write"`
	Idle    time.Duration `mapstructure:"idle"`
	Request time.Duration `mapstructure:"request"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	LogLevel        string        `mapstructure:"log_level"`
}

type LoggingConfig struct {
	Level             string `mapstructure:"level"`
	Format            string `mapstructure:"format"`
	Output            string `mapstructure:"output"`
	IncludeCaller     bool   `mapstructure:"include_caller"`
	IncludeStacktrace bool   `mapstructure:"include_stacktrace"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

type TracingConfig struct {
	Enabled    bool    `mapstructure:"enabled"`
	Provider   string  `mapstructure:"provider"`
	Endpoint   string  `mapstructure:"endpoint"`
	SampleRate float64 `mapstructure:"sample_rate"`
}

func Load() (*Config, error) {
	v := viper.New()

	environment := os.Getenv("APP_ENVIRONMENT")
	if environment == "" {
		environment = "local"
	}

	configPath := "./foundation/config"
	configFile := fmt.Sprintf("%s.yaml", environment)

	v.SetConfigName(environment)
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	v.AddConfigPath("./config")             // Alternative path for flexibility
	v.AddConfigPath("../config")            // For tests running from subdirectories
	v.AddConfigPath("../../config")         // For tests running from deeper subdirectories
	v.AddConfigPath("../foundation/config") // For service main.go files

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	v.AutomaticEnv()

	v.SetEnvPrefix("APP")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

func LoadWithPath(configPath, environment string) (*Config, error) {
	v := viper.New()

	if environment == "" {
		environment = "local"
	}

	v.SetConfigName(environment)
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file from %s: %w", configPath, err)
	}

	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if cfg.App.Environment == "" {
		return fmt.Errorf("app.environment is required")
	}

	if cfg.Server.GRPC.Port <= 0 || cfg.Server.GRPC.Port > 65535 {
		return fmt.Errorf("server.grpc.port must be between 1 and 65535")
	}
	if cfg.Server.GRPC.Host == "" {
		return fmt.Errorf("server.grpc.host is required")
	}

	if cfg.Database.Postgres.Host == "" {
		return fmt.Errorf("database.postgres.host is required")
	}
	if cfg.Database.Postgres.Port <= 0 || cfg.Database.Postgres.Port > 65535 {
		return fmt.Errorf("database.postgres.port must be between 1 and 65535")
	}
	if cfg.Database.Postgres.User == "" {
		return fmt.Errorf("database.postgres.user is required")
	}
	if cfg.Database.Postgres.Name == "" {
		return fmt.Errorf("database.postgres.name is required")
	}

	if cfg.Database.Postgres.MaxOpenConns < 1 {
		return fmt.Errorf("database.postgres.max_open_conns must be at least 1")
	}
	if cfg.Database.Postgres.MaxIdleConns < 0 {
		return fmt.Errorf("database.postgres.max_idle_conns must be non-negative")
	}
	if cfg.Database.Postgres.MaxIdleConns > cfg.Database.Postgres.MaxOpenConns {
		return fmt.Errorf("database.postgres.max_idle_conns cannot exceed max_open_conns")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[cfg.Logging.Level] {
		return fmt.Errorf("logging.level must be one of: debug, info, warn, error, fatal")
	}

	if cfg.Metrics.Enabled {
		if cfg.Metrics.Port <= 0 || cfg.Metrics.Port > 65535 {
			return fmt.Errorf("metrics.port must be between 1 and 65535")
		}
		if cfg.Metrics.Path == "" {
			return fmt.Errorf("metrics.path is required when metrics are enabled")
		}
	}

	if cfg.Tracing.Enabled {
		if cfg.Tracing.Provider == "" {
			return fmt.Errorf("tracing.provider is required when tracing is enabled")
		}
		if cfg.Tracing.Endpoint == "" {
			return fmt.Errorf("tracing.endpoint is required when tracing is enabled")
		}
		if cfg.Tracing.SampleRate < 0.0 || cfg.Tracing.SampleRate > 1.0 {
			return fmt.Errorf("tracing.sample_rate must be between 0.0 and 1.0")
		}
	}

	return nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Postgres.Host,
		c.Database.Postgres.Port,
		c.Database.Postgres.User,
		c.Database.Postgres.Password,
		c.Database.Postgres.Name,
		c.Database.Postgres.SSLMode,
	)
}

func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.GRPC.Host, c.Server.GRPC.Port)
}

func (c *Config) GetMetricsAddress() string {
	if !c.Metrics.Enabled {
		return ""
	}
	return fmt.Sprintf(":%d", c.Metrics.Port)
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production" || c.App.Environment == "prod"
}

func (c *Config) IsLocal() bool {
	return c.App.Environment == "local"
}

func (c *Config) IsTesting() bool {
	return c.App.Environment == "test"
}