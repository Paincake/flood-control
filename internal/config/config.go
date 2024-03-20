package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env             string `env:"ENV" env_default:"local"`
	Port            string `env:"DB_PORT" env_default:"5432"`
	Host            string `env:"DB_HOST" env_default:"localhost"`
	Name            string `env:"DB_NAME" env_default:"postgres"`
	User            string `env:"DB_USER" env_default:"user"`
	Password        string `env:"DB_PASSWORD" env_default:"password"`
	KCallLimit      int    `env:"KCallLimit"`
	NSecondInterval int    `env:"NSecondInterval"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH env variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
