package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	duration = 1 * time.Minute
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go measureRabbitMQ(ch1)
	go measureKafka(ch2)

	<-ch1
	<-ch2
}

func measureRabbitMQ(ch1 chan string) {
	fmt.Println("Testing RabbitMQ")
	apiURL := "http://localhost:3000/placeorder"

	requestBody := []byte(`{
		"id": "a68b823c-7ca6-44bc-b721-fb4d5312cafc",
		"name": "KFC",
		"quantity": 999
	  }`)

	client := &http.Client{}
	startTime := time.Now()
	iteration := 0
	for {
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		_ = body

		iteration++
		if time.Since(startTime) >= duration {
			break
		}
	}
	fmt.Println("Number of task RabbitMQ can produce:", iteration)
	ch1 <- "Done"
}

func measureKafka(ch2 chan string) {
	fmt.Println("Testing Kafka")
	apiURL := "http://localhost:8081/produce"

	requestBody := []byte(`{
		"config": {
		  "retries": 0
		},
		"message": "Hello World!",
		"topic": "Italian"
	}`)

	client := &http.Client{}
	startTime := time.Now()
	iteration := 0
	for {
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		_ = body

		iteration++
		if time.Since(startTime) >= duration {
			break
		}
	}
	fmt.Println("Number of task Kafka can produce:", iteration)
	ch2 <- "Done"
}