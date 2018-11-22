package main

import (
	"fmt"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(e error, msg string) {
	if e != nil {
		fmt.Printf("%v: %v", msg, e)
		os.Exit(0)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://wk:112851@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World %v!"
	for {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf(body, time.Now())),
			})
		failOnError(err, "Publish")
		fmt.Println("Ping!")
		time.Sleep(time.Minute * 10)
	}
}
