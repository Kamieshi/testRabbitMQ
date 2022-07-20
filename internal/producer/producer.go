// Package producer RabbitMQ
package producer

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"testRabitMQ/internal/models"
)

type producerRabbit struct {
	conn   *amqp.Connection
	queue  *amqp.Queue
	chanel *amqp.Channel
}

// NewProducer constructor
func NewProducer(urlCon, queueName string) (*producerRabbit, error) {
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

	return &producerRabbit{
		conn:   conn,
		queue:  &q,
		chanel: ch,
	}, err
}

// SendMessage send message to queue
func (p *producerRabbit) SendMessage(ctx context.Context, message *models.Message) error {
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
	return err
}
