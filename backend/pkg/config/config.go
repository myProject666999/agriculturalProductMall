package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Charset  string
	ParseTime bool `mapstructure:"parse_time"`
	Loc      string
}

type JWTConfig struct {
	Secret string
	Expire int
}

type EmailConfig struct {
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
	Username string
	Password string
	From     string
}

type UploadConfig struct {
	Path    string
	MaxSize int `mapstructure:"max_size"`
}

var App *Config

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("/config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	App = &Config{}
	if err := viper.Unmarshal(App); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
