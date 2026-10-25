package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestApplicationConfigPointTypes(t *testing.T) {
	data, err := os.ReadFile("../resources/application.yml")
	if err != nil {
		t.Fatalf("read application config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse application config: %v", err)
	}
	if err := validatePointTypes(cfg.PointTypes); err != nil {
		t.Fatalf("validate point types: %v", err)
	}
}

func TestValidatePointTypesRejectsInvalidDirections(t *testing.T) {
	err := validatePointTypes(map[int]PointTypeConfig{
		5: "other",
	})
	if err == nil {
		t.Fatal("expected unsupported point type direction to fail validation")
	}
}
