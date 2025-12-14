package config

import (
	"sso/internal/logger"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

type AuthConfig struct {
	Secret string `yaml:"secret"`
}

type ServerConfig struct {
	GrpcPort int `yaml:"grpc_port"`
	HttpPort int `yaml:"http_port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	SslMode  string `yaml:"sslmode"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		logger.Logger().Error("config don't load", zap.Error(err))
		return nil, err
	}
	return &cfg, nil
}
