package consumer

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"testRabitMQ/internal/models"
)

type consumerRabbit struct {
	conn   *amqp.Connection
	queue  *amqp.Queue
	chanel *amqp.Channel
}

// NewConsumer constructor
func NewConsumer(urlCon, queueName string) (*consumerRabbit, error) {
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

	return &consumerRabbit{
		conn:   conn,
		queue:  &q,
		chanel: ch,
	}, err
}

func (c *consumerRabbit) GetMessage() (*models.Message, error) {
	msgs, err := c.chanel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		return nil, err
	}

	var messages []*models.Message
	var message models.Message
	for d := range msgs {
		if messages == nil {
			messages = make([]*models.Message, 0, len(msgs))
		}

		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			logrus.Error(err)
			continue
		}
		break
	}
	return &message, nil
}
