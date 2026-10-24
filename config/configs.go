package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Port     int            `yaml:"port"`
}

type DatabaseConfig struct {
	ConnectionString string `yaml:"connectionString"`
}

var (
	configInstance *Config
)

func GetConfig() *Config {
	return configInstance
}

func LoadConfig() error {
	data, err := os.ReadFile("resources/application.yml")
	if err != nil {
		return err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	configInstance = &cfg

	return nil
}
