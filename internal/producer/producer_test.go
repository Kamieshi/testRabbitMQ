package producer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"testRabitMQ/internal/models"
)

func TestSendMessage(t *testing.T) {
	prod, err := NewProducer("amqp://user:pass@localhost:5672/", "test_queue")
	if err != nil {
		t.Fatal(err)
	}
	mes := models.Message{
		ID:      uuid.UUID{},
		PayLoad: "Some payLoad",
	}
	for i := 0; ; i++ {
		mes.ID = uuid.New()
		err = prod.SendMessage(context.Background(), &mes)
		if err != nil {
			t.Fatal(t)
		}
		if i%10000 == 0 {
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}
	}
}
