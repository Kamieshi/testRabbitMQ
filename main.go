package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://user:pass@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	s := sync.Mutex{}
	body := "Hello World!"
	go func() {
		for i := 0; ; i++ {
			s.Lock()
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprint(body, " ", i)),
				})
			s.Unlock()
			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		for {
			s.Lock()
			msgs, _ := ch.Consume(
				q.Name, // queue
				"",     // consumer
				true,   // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
			)
			s.Unlock()

			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)

			}
		}
	}()

	time.Sleep(30 * time.Second)
}
