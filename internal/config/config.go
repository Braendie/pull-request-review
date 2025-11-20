package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server Server
}

type Server struct {
	Address string `env:"ADDRESS" envDefault:"localhost"`
	Port    int `env:"PORT" envDefault:"8080"`
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func Load() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, fmt.Errorf("Failed to load from .env file: %v", err)
	}

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse environment variables:%v", err)
	}
	
	return &cfg, nil
}
