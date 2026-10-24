package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	ConnectionString string `yaml:"connection_string"`
	SchemaFile       string `yaml:"schema_file"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config yaml: %w", err)
	}

	// Allow environment variable override for Docker / CI environments
	if envConn := os.Getenv("DB_CONNECTION_STRING"); envConn != "" {
		cfg.Database.ConnectionString = envConn
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	return &cfg, nil
}
