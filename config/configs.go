package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database   DatabaseConfig          `yaml:"database"`
	Port       int                     `yaml:"port"`
	PointTypes map[int]PointTypeConfig `yaml:"point_types"`
}

type PointTypeConfig string

const (
	PointTypeCredit PointTypeConfig = "credit"
	PointTypeDebit  PointTypeConfig = "debit"
)

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

	if err := validatePointTypes(cfg.PointTypes); err != nil {
		return err
	}

	configInstance = &cfg

	return nil
}

func validatePointTypes(pointTypes map[int]PointTypeConfig) error {
	if len(pointTypes) == 0 {
		return fmt.Errorf("point_types must contain at least one point type")
	}

	for id, direction := range pointTypes {
		if id <= 0 {
			return fmt.Errorf("point_types id must be positive: %d", id)
		}
		if direction != PointTypeCredit && direction != PointTypeDebit {
			return fmt.Errorf("point_types[%d] must be %q or %q, got %q", id, PointTypeCredit, PointTypeDebit, direction)
		}
	}

	return nil
}
