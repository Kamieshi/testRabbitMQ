package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"testRabitMQ/internal/conf"
	"testRabitMQ/internal/models"
	"testRabitMQ/internal/repository"
)

func TestGetMessages(t *testing.T) {
	prod, err := NewConsumer("amqp://user:pass@localhost:5672/", "test_queue")
	if err != nil {
		t.Fatal(err)
	}
	chErr := make(chan error)
	chMess := make(chan *models.Message)
	ctxWT, _ := context.WithTimeout(context.Background(), 100*time.Second)
	go prod.GetMessage(ctxWT, chErr, chMess)
	select {
	case err = <-chErr:
		t.Fatal(err)
	case mess := <-chMess:
		fmt.Println(mess)
	}
}

func TestStreamListeningAndWriteToDB(t *testing.T) {
	prod, err := NewConsumer("amqp://user:pass@localhost:5672/", "test_queue")
	if err != nil {
		t.Fatal(err)
	}
	chErr := make(chan error)

	config, err := conf.NewConfig()
	if err != nil {
		panic(fmt.Errorf("repository.MainTest: %v", err))
	}

	pool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v", config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort, config.PostgresDB))
	if err != nil {
		panic(fmt.Errorf("repository.MainTest poolconnection: %v", err))
	}
	repMessage := repository.NewMessageRepositoryPostgres(pool)

	ctxWT, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	go prod.StreamListeningAndWriteToDB(ctxWT, repMessage)

	for err = range chErr {
		if err == ctxWT.Err() {
			break
		}
		fmt.Println(err)
	}

}
