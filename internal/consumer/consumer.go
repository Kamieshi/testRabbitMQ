package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v4"

	"testRabitMQ/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"

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

func (c *consumerRabbit) GetMessage(ctx context.Context, chErr chan error, chOut chan *models.Message) {
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
		chErr <- err
		return
	}
	select {
	case <-ctx.Done():
		chErr <- ctx.Err()
	case d := <-msgs:
		var message models.Message
		err = json.Unmarshal(d.Body, &message)
		if err != nil {
			chErr <- err
			return
		}
		chOut <- &message
	}
}

func (c *consumerRabbit) StreamListening(ctx context.Context, chErr chan error, chOut chan *models.Message) {
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
		chErr <- err
		return
	}
	for {
		select {
		case <-ctx.Done():
			chErr <- ctx.Err()
		case d := <-msgs:
			var message models.Message
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				chErr <- err
				return
			}
			chOut <- &message
		}
	}
}

func (c *consumerRabbit) StreamListeningAndWriteToDB(ctx context.Context, chErr chan error, rep *repository.MessageRepositoryPostgres) {
	msgs, err := c.chanel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	var batchPgx pgx.Batch
	if err != nil {
		chErr <- err
		return
	}
	var message models.Message
	nextTime := time.Now().Add(10 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			chErr <- ctx.Err()
			if batchPgx.Len() != 0 {
				err = rep.BatchQuery(ctx, &batchPgx)
				if err != nil {
					chErr <- err
				}
			}
			return
		case d := <-msgs:

			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				chErr <- err
				continue
			}
			batchPgx.Queue("INSERT INTO messages(id,pay_load) values ($1,$2)", message.ID, message.PayLoad)
			if batchPgx.Len() > 1000 || nextTime.Before(time.Now()) {
				err = rep.BatchQuery(ctx, &batchPgx)
				if err != nil {
					chErr <- err
					continue
				}
				batchPgx = pgx.Batch{}
				nextTime = time.Now().Add(10 * time.Millisecond)
			}
		default:
			if batchPgx.Len() != 0 && nextTime.Before(time.Now()) {
				err = rep.BatchQuery(ctx, &batchPgx)
				if err != nil {
					chErr <- err
				}
				batchPgx = pgx.Batch{}
				nextTime = time.Now().Add(10 * time.Millisecond)
			}
		}
	}
}
