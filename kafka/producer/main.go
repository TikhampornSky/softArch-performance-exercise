package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const numMessages = 200000

func main() {
	startTime := time.Now()

	// to produce messages
	topic := "my-topic"
	partition := 0

	message := "Some Order are place..."

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9093", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9093"},
		Topic:   topic,
	})

	messages := []kafka.Message{}
	for i := 0; i < numMessages; i++ {
		messages = append(messages, kafka.Message{Value: []byte(message + " " + strconv.Itoa(i))})
	}

	err = w.WriteMessages(context.Background(), messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	elapsedTime := time.Since(startTime)
	log.Println("Send message to kafka: ", numMessages, " messages")
	log.Printf("Execution time: %v seconds\n", elapsedTime.Seconds())
}
