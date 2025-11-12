package queue

import (
	"context"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConnection is a mock for amqp.Connection
type MockConnection struct {
	mock.Mock
}

func (m *MockConnection) Channel() (*amqp.Channel, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*amqp.Channel), args.Error(1)
}

func (m *MockConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewQueue_Success(t *testing.T) {
	// This test would require a real RabbitMQ connection or more sophisticated mocking
	// For now, we'll just test the error cases
	t.Skip("Requires RabbitMQ connection or advanced mocking")
}

func TestRabbitMQ_Publish_ClosedChannel(t *testing.T) {
	r := &RabbitMQ{
		conn:    nil,
		channel: nil,
	}

	err := r.Publish(context.Background(), "", "test-key", []byte("test"))
	assert.ErrorIs(t, err, ErrQueueClosed)
}

func TestRabbitMQ_Consume_ClosedChannel(t *testing.T) {
	r := &RabbitMQ{
		conn:    nil,
		channel: nil,
	}

	handler := func([]byte) error {
		return nil
	}

	err := r.Consume("test-queue", handler)
	assert.ErrorIs(t, err, ErrQueueClosed)
}

func TestRabbitMQ_DeclareQueue_ClosedChannel(t *testing.T) {
	r := &RabbitMQ{
		conn:    nil,
		channel: nil,
	}

	err := r.DeclareQueue("test-queue")
	assert.ErrorIs(t, err, ErrQueueClosed)
}

func TestRabbitMQ_Close_NilConnAndChannel(t *testing.T) {
	r := &RabbitMQ{
		conn:    nil,
		channel: nil,
	}

	err := r.Close()
	assert.NoError(t, err)
}
