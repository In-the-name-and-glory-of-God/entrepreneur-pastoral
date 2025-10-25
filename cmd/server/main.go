package main

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/cache"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/database"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/logger"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Application.Env)

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	defer db.Close()

	log.Info("database connection established")

	rdb, err := cache.New(cfg.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis", err)
	}
	defer rdb.Close()

	log.Info("redis connection established")

	log.Info("starting server")
}
