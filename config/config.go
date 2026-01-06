package config

import (
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName      string `mapstructure:"APP_NAME" env:"APP_NAME"`
	AppEnv       string `mapstructure:"APP_ENV" env:"APP_ENV"`
	AppVersion   string `mapstructure:"APP_VERSION" env:"APP_VERSION"`
	AppCommit    string `mapstructure:"APP_COMMIT" env:"APP_COMMIT"`
	AppBuildTime string `mapstructure:"APP_BUILD_TIME" env:"APP_BUILD_TIME"`
	AppURL       string `mapstructure:"APP_URL" env:"APP_URL"`
	AppClientURL string `mapstructure:"APP_CLIENT_URL" env:"APP_CLIENT_URL"`

	ServerHost string `mapstructure:"SERVER_HOST" env:"SERVER_HOST"`
	ServerPort string `mapstructure:"SERVER_PORT" env:"SERVER_PORT"`

	DBHost     string `mapstructure:"DB_HOST" env:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT" env:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER" env:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD" env:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_DB" env:"DB_DB"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE" env:"DB_SSL_MODE"`

	RedisHost     string `mapstructure:"REDIS_HOST" env:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT" env:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" env:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB" env:"REDIS_DB"`

	SMTPHost     string `mapstructure:"SMTP_HOST" env:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT" env:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER" env:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD" env:"SMTP_PASSWORD"`

	LoggerFile       string        `mapstructure:"LOGGER_FILE" env:"LOGGER_FILE"`
	LoggerLevel      string        `mapstructure:"LOGGER_LEVEL" env:"LOGGER_LEVEL"`
	LoggerMaxSize    int           `mapstructure:"LOGGER_MAX_SIZE" env:"LOGGER_MAX_SIZE"`
	LoggerMaxBackups int           `mapstructure:"LOGGER_MAX_BACKUPS" env:"LOGGER_MAX_BACKUPS"`
	LoggerMaxAge     int           `mapstructure:"LOGGER_MAX_AGE" env:"LOGGER_MAX_AGE"`
	LoggerCompress   bool          `mapstructure:"LOGGER_COMPRESS" env:"LOGGER_COMPRESS"`
	LoggerOutput     string        `mapstructure:"LOGGER_OUTPUT" env:"LOGGER_OUTPUT"`
	LoggerTick       time.Duration `mapstructure:"LOGGER_TICK" env:"LOGGER_TICK"`
	LoggerThreshold  int           `mapstructure:"LOGGER_THRESHOLD" env:"LOGGER_THRESHOLD"`
	LoggerRate       float64       `mapstructure:"LOGGER_RATE" env:"LOGGER_RATE"`

	AuthJWTSecret          string        `mapstructure:"AUTH_JWT_SECRET" env:"AUTH_JWT_SECRET"`
	AuthAccessTokenExpiry  time.Duration `mapstructure:"AUTH_ACCESS_TOKEN_EXPIRY" env:"AUTH_ACCESS_TOKEN_EXPIRY"`
	AuthRefreshTokenExpiry time.Duration `mapstructure:"AUTH_REFRESH_TOKEN_EXPIRY" env:"AUTH_REFRESH_TOKEN_EXPIRY"`

	// OpenTelemetry configuration
	OTLPEndpoint      string  `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT" env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTLPEnableTracing bool    `mapstructure:"OTEL_ENABLE_TRACING" env:"OTEL_ENABLE_TRACING"`
	OTLPEnableMetrics bool    `mapstructure:"OTEL_ENABLE_METRICS" env:"OTEL_ENABLE_METRICS"`
	OTLPSampleRate    float64 `mapstructure:"OTEL_SAMPLE_RATE" env:"OTEL_SAMPLE_RATE"`

	// Web Push VAPID configuration
	VapidPublicKey  string `mapstructure:"VAPID_PUBLIC_KEY" env:"VAPID_PUBLIC_KEY"`
	VapidPrivateKey string `mapstructure:"VAPID_PRIVATE_KEY" env:"VAPID_PRIVATE_KEY"`
	VapidSubject    string `mapstructure:"VAPID_SUBJECT" env:"VAPID_SUBJECT"`

	// Google OAuth configuration
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID" env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET" env:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL  string `mapstructure:"GOOGLE_CALLBACK_URL" env:"GOOGLE_CALLBACK_URL"`

	// Event (Canonical Log Lines) configuration
	EventSampleRate     float64 `mapstructure:"EVENT_SAMPLE_RATE" env:"EVENT_SAMPLE_RATE"`
	EventP99ThresholdMs int64   `mapstructure:"EVENT_P99_THRESHOLD_MS" env:"EVENT_P99_THRESHOLD_MS"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func (c *Config) RedisDSN() string {
	return net.JoinHostPort(c.RedisHost, fmt.Sprintf("%d", c.RedisPort))
}

// Validate checks if all required configuration fields are properly set
func (c *Config) Validate() error {
	var errors []string

	// Required database fields
	if c.DBHost == "" {
		errors = append(errors, "DB_HOST is required")
	}
	if c.DBPort == 0 {
		errors = append(errors, "DB_PORT is required")
	}
	if c.DBUser == "" {
		errors = append(errors, "DB_USER is required")
	}
	if c.DBName == "" {
		errors = append(errors, "DB_DB is required")
	}

	// Required Redis fields
	if c.RedisHost == "" {
		errors = append(errors, "REDIS_HOST is required")
	}
	if c.RedisPort == 0 {
		errors = append(errors, "REDIS_PORT is required")
	}

	// Required auth fields
	if c.AuthJWTSecret == "" {
		errors = append(errors, "AUTH_JWT_SECRET is required")
	}
	if len(c.AuthJWTSecret) < 32 {
		errors = append(errors, "AUTH_JWT_SECRET must be at least 32 characters")
	}
	if c.AuthAccessTokenExpiry == 0 {
		errors = append(errors, "AUTH_ACCESS_TOKEN_EXPIRY is required")
	}
	if c.AuthRefreshTokenExpiry == 0 {
		errors = append(errors, "AUTH_REFRESH_TOKEN_EXPIRY is required")
	}

	// Validate server config
	if c.ServerPort == "" {
		errors = append(errors, "SERVER_PORT is required")
	}

	// Validate VAPID config
	if c.VapidPublicKey == "" {
		errors = append(errors, "VAPID_PUBLIC_KEY is required")
	}
	if c.VapidPrivateKey == "" {
		errors = append(errors, "VAPID_PRIVATE_KEY is required")
	}
	if c.VapidSubject == "" {
		errors = append(errors, "VAPID_SUBJECT is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// setDefaults sets sensible defaults for optional configuration fields
func (c *Config) setDefaults() {
	// Application defaults
	if c.AppName == "" {
		c.AppName = "ethos-go"
	}
	if c.AppEnv == "" {
		c.AppEnv = "development"
	}

	// Server defaults
	if c.ServerHost == "" {
		c.ServerHost = "0.0.0.0"
	}

	// Database defaults
	if c.DBSSLMode == "" {
		c.DBSSLMode = "disable"
	}

	// Logger defaults
	if c.LoggerLevel == "" {
		c.LoggerLevel = "info"
	}
	if c.LoggerOutput == "" {
		c.LoggerOutput = "stdout"
	}
	if c.LoggerMaxSize == 0 {
		c.LoggerMaxSize = 100 // 100 MB
	}
	if c.LoggerMaxBackups == 0 {
		c.LoggerMaxBackups = 3
	}
	if c.LoggerMaxAge == 0 {
		c.LoggerMaxAge = 28 // 28 days
	}

	// Event defaults
	if c.EventSampleRate == 0 {
		c.EventSampleRate = 0.05 // 5% sampling for normal requests
	}
	if c.EventP99ThresholdMs == 0 {
		c.EventP99ThresholdMs = 2000 // 2 seconds
	}
}

/*
        +------------------+
        |   Environment    |   ← Highest Priority
        |   Variables      |
        +--------+---------+
                 |
                 v
        +------------------+
        |     .env File    |
        |  (v.ReadInConfig)
        +--------+---------+
                 |
                 v
        +------------------+
        |   Default Values |
        | (cfg.setDefaults)
        +--------+---------+
                 |
                 v
        +------------------+
        |   Final Config   |
        |    (cfg struct)  |
        +------------------+

Priority Resolution Rule:
ENV > .env > default

Explanation:
- v.ReadInConfig() loads .env
- v.AutomaticEnv() overrides any matching key from .env
- setDefaults() fills only missing values
*/

func Load() (*Config, error) {
	v := viper.New()

	// Set configuration file details
	v.SetConfigFile(".env")
	v.AddConfigPath(".")

	// Read .env file, if it exists
	if err := v.ReadInConfig(); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			slog.Warn("No .env file found, relying on environment variables")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Enable reading from environment variables
	v.AutomaticEnv()

	// Explicitly bind environment variables based on mapstructure tags
	// This is necessary because Viper's Unmarshal doesn't automatically look for env vars
	// corresponding to struct fields unless those keys are explicitly known/bound.
	t := reflect.TypeOf(Config{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			if err := v.BindEnv(tag); err != nil {
				slog.Error("Failed to bind env var", "key", tag, "error", err)
			}
		}
	}

	// Unmarshal configuration into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults for optional fields
	cfg.setDefaults()

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
