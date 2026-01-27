package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Logger   Logger   `mapstructure:"logger"`
	JWT      JWT      `mapstructure:"jwt"`
	Postgres Postgres `mapstructure:"postgres"`
	Redis    Redis    `mapstructure:"redis"`
	Minio    Minio    `mapstructure:"minio"`
}

type App struct {
	Name             string `mapstructure:"name"`
	Port             int    `mapstructure:"port"`
	Host             string `mapstructure:"host"`
	Env              string `mapstructure:"env"`
	TelegramBotToken string `mapstructure:"telegram_bot_token"`
	TelegramChatID   string `mapstructure:"telegram_chat_id"`

	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
}

func (a *App) GetDSN() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *App) IsDev() bool {
	return a.Env == "development"
}

type Logger struct {
	Level       string `mapstructure:"level"`
	LogDir      string `mapstructure:"directory"`
	Filename    string `mapstructure:"filename"`
	MaxSize     int    `mapstructure:"max_size"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      int    `mapstructure:"max_age"`
	Compress    bool   `mapstructure:"compress"`
	Console     bool   `mapstructure:"console"`
	RotateDaily bool   `mapstructure:"rotate_daily"`
}

type JWT struct {
	SecretKey            string        `mapstructure:"secret_key"`
	AccessExpireMinutes  time.Duration `mapstructure:"access_expire_minutes"`
	RefreshExpireMinutes time.Duration `mapstructure:"refresh_expire_minutes"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`

	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

func (p *Postgres) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		p.Host,
		p.User,
		p.Password,
		p.DBName,
		p.Port,
		p.SSLMode,
		p.TimeZone,
	)
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`

	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxRetries      int           `mapstructure:"max_retries"`
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	DialTimeout     time.Duration `mapstructure:"dial_timeout"`
}

func (r *Redis) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type Minio struct {
	Host         string         `mapstructure:"host"`
	RootUser     string         `mapstructure:"root_user"`
	RootPassword string         `mapstructure:"root_password"`
	APIPort      int            `mapstructure:"api_port"`
	ConsolePort  int            `mapstructure:"console_port"`
	Buckets      []BucketConfig `mapstructure:"buckets"`
}

type BucketConfig struct {
	Name   string `mapstructure:"name"`
	Public bool   `mapstructure:"public"`
}

func (m *Minio) GetAddr() string {
	return fmt.Sprintf("%s:%d", m.Host, m.APIPort)
}

func Load(path string) (*Config, error) {
	_ = gotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config.yml: %w", err)
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
