package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"testRabitMQ/internal/conf"
	"testRabitMQ/internal/consumer"
	"testRabitMQ/internal/repository"
)

func main() {
	config, err := conf.NewConfig()
	if err != nil {
		log.WithError(err).Fatal()
	}

	prod, err := consumer.NewConsumer(
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

	pool, err := pgxpool.Connect(
		context.Background(),
		fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v",
			config.PostgresUser,
			config.PostgresPassword,
			config.PostgresHost,
			config.PostgresPort,
			config.PostgresDB,
		))
	if err != nil {
		log.WithError(err).Fatal()
	}
	repMessage := repository.NewMessageRepositoryPostgres(pool)

	ctxWT, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	go prod.StreamListeningAndWriteToDB(ctxWT, repMessage)
	log.Info("Start listening")
	<-ctxWT.Done()
	log.Info("Stop listening")
}
