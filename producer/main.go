package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"testRabitMQ/internal/conf"
	"testRabitMQ/internal/models"
	"testRabitMQ/internal/producer"
)

func main() {
	config, err := conf.NewConfig()
	if err != nil {
		log.WithError(err).Fatal()
	}

	prod, err := producer.NewProducer(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			config.RabbitUser,
			config.RabbitPassword,
			config.RabbitHost,
			config.RabbitPort,
		),
		config.RabbitQueueName)
	if err != nil {
		log.WithError(err).Fatal()
	}

	mes := models.Message{
		ID:      uuid.UUID{},
		PayLoad: "Some payLoad",
	}
	for i := 0; ; i++ {
		mes.ID = uuid.New()
		err = prod.SendMessage(context.Background(), &mes)
		if err != nil {
			log.WithError(err).Fatal()
		}
		if i%10000 == 0 {
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}
	}
}
