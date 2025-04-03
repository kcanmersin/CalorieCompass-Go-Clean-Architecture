package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type (
	Config struct {
		App       `yaml:"app"`
		HTTP      `yaml:"http"`
		Log       `yaml:"logger"`
		Postgres  `yaml:"postgres"`
		JWT       `yaml:"jwt"`
		FatSecret `yaml:"fatsecret"`
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

	FatSecret struct {
		ClientID       string `yaml:"client_id"`
		ClientSecret   string `yaml:"client_secret"`
		ConsumerKey    string `yaml:"consumer_key"`
		ConsumerSecret string `yaml:"consumer_secret"`
	}
)

func substituteEnvVars(s string) string {
	if !strings.Contains(s, "${") {
		return s
	}

	return os.Expand(s, func(key string) string {
		return os.Getenv(key)
	})
}

func NewConfig(configPath string) (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

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

	// Override with environment variables if they exist
	if connStr := os.Getenv("CONNECTION_STRING"); connStr != "" {
		config.Postgres.URL = connStr
	}

	// Load FatSecret API credentials from environment variables
	if clientID := os.Getenv("FATSECRET_CLIENT_ID"); clientID != "" {
		config.FatSecret.ClientID = clientID
	}
	if clientSecret := os.Getenv("FATSECRET_CLIENT_SECRET"); clientSecret != "" {
		config.FatSecret.ClientSecret = clientSecret
	}
	if consumerKey := os.Getenv("FATSECRET_CONSUMER_KEY"); consumerKey != "" {
		config.FatSecret.ConsumerKey = consumerKey
	}
	if consumerSecret := os.Getenv("FATSECRET_CONSUMER_SECRET"); consumerSecret != "" {
		config.FatSecret.ConsumerSecret = consumerSecret
	}

	return config, nil
}
