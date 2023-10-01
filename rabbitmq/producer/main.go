package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var (
	routing_key = "Italian"
	numMessages = 200000
)

func main() {
	startTime := time.Now()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"my_exchange", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Italian food is the best!"

	for i := 0; i < numMessages; i++ {
		err = ch.PublishWithContext(ctx,
			"my_exchange", // exchange
			routing_key,   // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")	
	}
	
	elapsedTime := time.Since(startTime)
	log.Printf("Sent %d messages", numMessages)
	log.Printf("Execution time: %v seconds\n", elapsedTime.Seconds())
}
