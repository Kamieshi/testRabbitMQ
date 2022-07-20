package consumer

import (
	"fmt"
	"testing"
)

func TestGetMessages(t *testing.T) {
	prod, err := NewConsumer("amqp://user:pass@localhost:5672/", "test_queue")
	if err != nil {
		t.Fatal(err)
	}
	messages, err := prod.GetMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(messages)
}
