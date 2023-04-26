package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level      string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
		Structured bool   `yaml:"structured"   env:"STRUCTURED" env-default:"false"`
	}
)

// New creates new config based on config.yml and env
func New() (*Config, error) {
	cfg := Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = cleanenv.ReadConfig(dir+"/config.yml", &cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return &cfg, nil
}
