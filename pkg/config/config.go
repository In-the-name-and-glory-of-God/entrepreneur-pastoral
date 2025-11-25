package config

import (
	"fmt"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/constants"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/env"
)

type (
	Config struct {
		Application Application
		API         API
		Database    Database
		Redis       Redis
		RabbitMQ    RabbitMQ
		SMTP        SMTP
	}

	Application struct {
		Secret string
		Name   string
		Env    string
	}

	API struct {
		Host         string
		Port         int
		WriteTimeout time.Duration
		ReadTimeout  time.Duration
		IdleTimeout  time.Duration
		CORS         CORS
		RateLimiter  RateLimiter
	}

	CORS struct {
		AllowedOrigins   []string
		AllowedMethods   []string
		AllowedHeaders   []string
		ExposedHeaders   []string
		AllowCredentials bool
		MaxAge           int
	}

	RateLimiter struct {
		RequestLimit int
		WindowLength time.Duration
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

	SMTP struct {
		Host     string
		Port     int
		User     string
		Password string
		From     string
	}
)

func Load() Config {
	return Config{
		Application: Application{
			Secret: env.GetString("APP_SECRET", "my-supa-dupa-app-secret-yes-it-is-okay"),
			Name:   env.GetString("APP_NAME", "entrepreneur-pastoral"),
			Env:    env.GetString("APP_ENV", "development"),
		},
		API: API{
			Host:         env.GetString("API_HOST", "localhost"),
			Port:         env.GetInt("API_PORT", 8080),
			WriteTimeout: env.GetDuration("API_WRITE_TIMEOUT", 30*time.Second),
			ReadTimeout:  env.GetDuration("API_READ_TIMEOUT", 10*time.Second),
			IdleTimeout:  env.GetDuration("API_IDLE_TIMEOUT", 1*time.Minute),
			CORS: CORS{
				AllowedOrigins:   env.GetStringSlice("API_CORS_ALLOWED_ORIGINS", []string{"*"}),
				AllowedMethods:   env.GetStringSlice("API_CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
				AllowedHeaders:   env.GetStringSlice("API_CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
				ExposedHeaders:   env.GetStringSlice("API_CORS_EXPOSED_HEADERS", []string{"Link"}),
				AllowCredentials: env.GetBool("API_CORS_ALLOW_CREDENTIALS", false),
				MaxAge:           env.GetInt("API_CORS_MAX_AGE", 300),
			},
			RateLimiter: RateLimiter{
				RequestLimit: env.GetInt("API_RATE_LIMITER_REQUEST_LIMIT", 10),
				WindowLength: env.GetDuration("API_RATE_LIMITER_WINDOW_LENGTH", 1*time.Minute),
			},
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
		SMTP: SMTP{
			Host:     env.GetString("SMTP_HOST", "smtp.gmail.com"),
			Port:     env.GetInt("SMTP_PORT", 587),
			User:     env.GetString("SMTP_USER", "user@example.com"),
			Password: env.GetString("SMTP_PASSWORD", "password"),
			From:     env.GetString("SMTP_FROM", "user@example.com"),
		},
	}
}

func (a Application) IsDevelopment() bool {
	return a.Env == constants.DEVELOPMENT
}

func (a Application) IsProduction() bool {
	return a.Env == constants.PRODUCTION
}

func (a API) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
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

func (e SMTP) DSN() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}
