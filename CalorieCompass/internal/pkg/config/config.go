package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		Postgres `yaml:"postgres"`
		JWT      `yaml:"jwt"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	HTTP struct {
		Port string `yaml:"port"`
	}

	Log struct {
		Level string `yaml:"log_level"`
	}

	Postgres struct {
		PoolMax int    `yaml:"pool_max"`
		URL     string `yaml:"url"`
	}

	JWT struct {
		Secret         string `yaml:"secret"`
		ExpirationHour int    `yaml:"expiration_hours"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("open config file error: %w", err)
	}
	defer configFile.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("decode config error: %w", err)
	}

	return config, nil
}