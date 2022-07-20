package producer

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"testRabitMQ/internal/models"
)

func TestSendMessage(t *testing.T) {
	prod, err := NewProducer("amqp://user:pass@localhost:5672/", "test_queue")
	if err != nil {
		t.Fatal(err)
	}
	mes := models.Message{
		ID:      uuid.New(),
		PayLoad: "Some payLoad",
	}
	err = prod.SendMessage(context.Background(), &mes)
	if err != nil {
		t.Fatal(t)
	}
}
