package queue

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrQueueClosed = errors.New("queue connection closed")
)

type Queue interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	Consume(queueName string, handler func([]byte) error) error
	DeclareQueue(queueName string) error
	Close() error
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewQueue(conn *amqp.Connection) (Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	if r.channel == nil {
		return ErrQueueClosed
	}

	return r.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Consume(queueName string, handler func([]byte) error) error {
	if r.channel == nil {
		return ErrQueueClosed
	}

	msgs, err := r.channel.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				// Negative acknowledgment - message will be requeued
				msg.Nack(false, true)
			} else {
				// Positive acknowledgment - message will be removed from queue
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) DeclareQueue(queueName string) error {
	if r.channel == nil {
		return ErrQueueClosed
	}

	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	return err
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return err
		}
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}
