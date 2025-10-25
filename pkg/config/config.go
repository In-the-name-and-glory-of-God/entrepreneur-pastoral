package config

import (
	"fmt"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/env"
)

type (
	Config struct {
		Application Application
		Database    Database
		Redis       Redis
		RabbitMQ    RabbitMQ
	}

	Application struct {
		Env string
	}

	Database struct {
		Host            string
		Port            int
		User            string
		Password        string
		Name            string
		Params          string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
	}

	Redis struct {
		Protocol int
		Host     string
		Port     int
		Password string
		DB       int
	}

	RabbitMQ struct {
		Host     string
		Port     int
		User     string
		Password string
	}
)

func Load() Config {
	return Config{
		Application: Application{
			Env: env.GetString("APP_ENV", "development"),
		},
		Database: Database{
			Host:            env.GetString("DB_HOST", "localhost"),
			Port:            env.GetInt("DB_PORT", 5432),
			User:            env.GetString("DB_USER", "user"),
			Password:        env.GetString("DB_PASSWORD", "password"),
			Name:            env.GetString("DB_NAME", "entrepreneur-pastoral"),
			Params:          env.GetString("DB_PARAMS", "sslmode=disable"),
			MaxOpenConns:    env.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    env.GetInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: env.GetDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: Redis{
			Protocol: env.GetInt("REDIS_PROTOCOL", 2),
			Host:     env.GetString("REDIS_HOST", "localhost"),
			Port:     env.GetInt("REDIS_PORT", 6379),
			Password: env.GetString("REDIS_PASSWORD", ""),
			DB:       env.GetInt("REDIS_DB", 0),
		},
		RabbitMQ: RabbitMQ{
			Host:     env.GetString("RABBITMQ_HOST", "localhost"),
			Port:     env.GetInt("RABBITMQ_PORT", 5672),
			User:     env.GetString("RABBITMQ_USER", "guest"),
			Password: env.GetString("RABBITMQ_PASSWORD", "guest"),
		},
	}
}

func (d Database) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", d.User, d.Password, d.Host, d.Port, d.Name, d.Params)
}

func (r Redis) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r RabbitMQ) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
}
