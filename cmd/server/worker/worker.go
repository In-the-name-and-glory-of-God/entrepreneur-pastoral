package worker

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/worker/email"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/worker/notification"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"go.uber.org/zap"
)

type Worker struct {
	Queue  storage.QueueStorage
	Config config.Config
	Logger *zap.SugaredLogger
}

func NewWorker(queue storage.QueueStorage, cfg config.Config, logger *zap.SugaredLogger) *Worker {
	return &Worker{
		Queue:  queue,
		Config: cfg,
		Logger: logger,
	}
}

func (w *Worker) Start() {
	emailService := email.NewSMTPService(w.Config.SMTP)
	notificationConsumer := notification.NewNotificationConsumer(w.Queue, emailService, w.Logger)

	if err := notificationConsumer.Start(); err != nil {
		w.Logger.Fatal("Failed to start notification consumer", "error", err)
	}

	w.Logger.Info("Worker started")
}
