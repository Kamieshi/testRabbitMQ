package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"testRabitMQ/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"

	"testRabitMQ/internal/models"
)

type Rabbit struct {
	conn   *amqp.Connection
	queue  *amqp.Queue
	chanel *amqp.Channel
}

const waitTimeMlSec = 10

// NewConsumer constructor
func NewConsumer(urlCon, queueName string) (*Rabbit, error) {
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

	return &Rabbit{
		conn:   conn,
		queue:  &q,
		chanel: ch,
	}, err
}

// GetMessage Get one message
func (c *Rabbit) GetMessage(ctx context.Context, chErr chan error, chOut chan *models.Message) {
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

// StreamListening Get more messages
func (c *Rabbit) StreamListening(ctx context.Context, chOut chan *models.Message) {
	msgCh, err := c.chanel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		log.WithError(err).Error()
		return
	}
	for {
		select {
		case <-ctx.Done():
			log.WithError(err).Info()
		case d := <-msgCh:
			var message models.Message
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				log.WithError(err).Error()
				return
			}
			chOut <- &message
		}
	}
}

// StreamListeningAndWriteToDB Write all queue in db
func (c *Rabbit) StreamListeningAndWriteToDB(ctx context.Context, rep *repository.MessageRepositoryPostgres) {
	msgCh, err := c.chanel.Consume(
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
		log.WithError(err).Error()
		return
	}
	var message models.Message
	nextTime := time.Now().Add(waitTimeMlSec * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			log.WithError(err).Info()
			if batchPgx.Len() != 0 {
				err = rep.BatchQuery(ctx, &batchPgx)
				if err != nil {
					log.WithError(err).Error()
				}
			}
			return
		case d := <-msgCh:

			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				log.WithError(err).Error()
				continue
			}
			batchPgx.Queue("INSERT INTO messages(id,pay_load) values ($1,$2)", message.ID, message.PayLoad)
			if batchPgx.Len() > 1000 || nextTime.Before(time.Now()) {
				err = rep.BatchQuery(ctx, &batchPgx)
				if err != nil {
					log.WithError(err).Error()
					continue
				}
				batchPgx = pgx.Batch{}
				nextTime = time.Now().Add(waitTimeMlSec * time.Millisecond)
			}
		default:
			if batchPgx.Len() != 0 && nextTime.Before(time.Now()) {
				err = rep.BatchQuery(ctx, &batchPgx)
				if err != nil {
					log.WithError(err).Error()
				}
				batchPgx = pgx.Batch{}
				nextTime = time.Now().Add(waitTimeMlSec * time.Millisecond)
			}
		}
	}
}
