package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/orchestrator"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/router"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/database"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/logger"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/queue"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"go.uber.org/zap"
)

var (
	cfg config.Config
	log *zap.SugaredLogger
)

func main() {
	cfg = config.Load()
	log = logger.New(cfg.Application)

	db, err := database.NewPostgresConn(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	defer db.Close()
	log.Info("database connection established")

	client, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis", err)
	}
	defer client.Close()
	cache := storage.NewCacheStorage(client)
	log.Info("redis connection established")

	rabbitConn, err := database.NewRabbitMQConn(cfg.RabbitMQ)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq", err)
	}
	defer rabbitConn.Close()
	queueService, err := queue.NewQueue(rabbitConn)
	if err != nil {
		log.Fatal("failed to create queue service", err)
	}
	defer queueService.Close()
	log.Info("rabbitmq connection established")

	tokenManager := auth.NewTokenManager(cfg.Application.Secret)
	orchestrator := orchestrator.New(cfg, log, db, cache, queueService, tokenManager)
	symphony := orchestrator.Compose()

	router := router.NewServerRouter(cfg, symphony)

	log.Fatal(run(router.Mount(client)))
}

func run(handler http.Handler) error {
	srv := &http.Server{
		Addr:         cfg.API.Addr(),
		Handler:      handler,
		WriteTimeout: cfg.API.WriteTimeout,
		ReadTimeout:  cfg.API.ReadTimeout,
		IdleTimeout:  cfg.API.IdleTimeout,
	}

	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Infow("Signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	log.Infow("Server has started", "addr", cfg.API.Addr(), "env", cfg.Application.Env)

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	log.Infow("Server has stopped", "addr", cfg.API.Addr(), "env", cfg.Application.Env)

	return nil
}
