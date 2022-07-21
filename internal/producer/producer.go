// Package producer RabbitMQ
package producer

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"testRabitMQ/internal/models"
)

// Rabbit Producer rabbit
type Rabbit struct {
	conn   *amqp.Connection
	queue  *amqp.Queue
	chanel *amqp.Channel
}

// NewProducer constructor
func NewProducer(urlCon, queueName string) (*Rabbit, error) {
	conn, err := amqp.Dial(urlCon)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("producer/NewProducer : %v", err)
	}
	return &Rabbit{
		conn:   conn,
		queue:  &q,
		chanel: ch,
	}, nil
}

// SendMessage send message to queue
func (p *Rabbit) SendMessage(ctx context.Context, message *models.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = p.chanel.PublishWithContext(
		ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	if err != nil {
		return fmt.Errorf("producer/SendMessage : %v", err)
	}
	return nil
}
