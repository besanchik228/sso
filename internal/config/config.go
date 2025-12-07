package config

import (
    "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
    GrpcPort int `yaml:"grpc_port" env:"GRPC_PORT" env-default:"50051"`
    HttpPort int `yaml:"http_port" env:"HTTP_PORT" env-default:"8080"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
    Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
    User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
    Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
    DbName   string `yaml:"dbname" env:"DB_NAME" env-default:"sso_db"`
    SslMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

func LoadConfig(path string) (*Config, error) {
    var cfg Config
    if err := cleanenv.ReadConfig(path, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
