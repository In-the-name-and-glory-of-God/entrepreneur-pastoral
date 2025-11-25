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
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/worker"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/database"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/logger"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"go.uber.org/zap"
)

var (
	cfg config.Config
	log *zap.SugaredLogger
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func main() {
	cfg = config.Load()
	log = logger.New(cfg.Application)

	db, err := database.NewPostgresConn(cfg.Database)
	failOnError(err, "failed to connect to database")
	defer db.Close()
	log.Info("database connection established")

	client, err := database.NewRedisClient(cfg.Redis)
	failOnError(err, "failed to connect to redis")
	defer client.Close()
	cache := storage.NewCacheStorage(client)
	log.Info("redis connection established")

	conn, err := database.NewRabbitMQConn(cfg.RabbitMQ)
	failOnError(err, "failed to connect to rabbitmq")
	defer conn.Close()
	queue, err := storage.NewQueueStorage(conn, log)
	failOnError(err, "failed to create queue service")
	defer queue.Close()
	log.Info("rabbitmq connection established")

	w := worker.NewWorker(queue, cfg, log)
	go w.Start()

	tokenManager := auth.NewTokenManager(cfg.Application.Secret)
	orchestrator := orchestrator.New(cfg, log, db, cache, queue, tokenManager)
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
