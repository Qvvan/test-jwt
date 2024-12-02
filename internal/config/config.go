package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Debug        string `env:"DEBUG" env-required:"true"`
	LogLevel     string `env:"LOG_LEVEL" env-default:"local"`
	PgDSN        string `env:"PG_DSN" env-required:"true"`
	JwtSecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
	HttpServer   `envPrefix:"HTTP_SERVER_"`
}

type HttpServer struct {
	Address     string        `env:"ADDRESS" env-default:"0.0.0.0:8881"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("error reading environment variables: %s", err.Error())
	}
	return &cfg
}
