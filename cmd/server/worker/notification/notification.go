package notification

import (
	"encoding/json"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/worker/email"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/constants"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"go.uber.org/zap"
)

type NotificationPayload struct {
	From         string   `json:"from"`
	To           []string `json:"to"`
	Subject      string   `json:"subject"`
	TemplateName string   `json:"template_name"`
	Data         any      `json:"data"`
}

type NotificationConsumer struct {
	queue        storage.QueueStorage
	emailService *email.SMTPService
	logger       *zap.SugaredLogger
}

func NewNotificationConsumer(queue storage.QueueStorage, emailService *email.SMTPService, logger *zap.SugaredLogger) *NotificationConsumer {
	return &NotificationConsumer{
		queue:        queue,
		emailService: emailService,
		logger:       logger,
	}
}

func (c *NotificationConsumer) Start() error {
	if err := c.queue.DeclareQueue(constants.QUEUE_NOTIFICATIONS); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", constants.QUEUE_NOTIFICATIONS, err)
	}

	return c.queue.Consume(constants.QUEUE_NOTIFICATIONS, c.handleMessage)
}

func (c *NotificationConsumer) handleMessage(body []byte) error {
	var payload NotificationPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.logger.Error("Failed to unmarshal notification payload", "error", err)
		return nil // Return nil to ack message as it is malformed and cannot be processed
	}

	if err := c.emailService.Send(payload.From, payload.To, payload.Subject, payload.TemplateName, payload.Data); err != nil {
		c.logger.Error("Failed to send email", "error", err)
		return err
	}

	c.logger.Info("Email sent successfully", "to", payload.To, "template", payload.TemplateName)
	return nil
}
