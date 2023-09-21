package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	numberOfExperiment = 2000
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
	for i := 0; i < numberOfExperiment; i++ {
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
	}
	elapsedTime := time.Since(startTime)
	fmt.Println("Total time for RabbitMQ:", elapsedTime)
	ch1 <- "Done"
}

func measureKafka(ch2 chan string) {
	fmt.Println("Testing Kafka")
	apiURL1 := "http://localhost:8081/produce"
	// apiURL2 := "http://localhost:8081/consume"

	requestBody1 := []byte(`{
		"config": {
		  "retries": 0
		},
		"message": "Hello World!",
		"topic": "Italian"
	}`)
	// requestBody2 := []byte(`{
	// 	"config": {
	// 	  "auto.offset.reset": "earliest"
	// 	},
	// 	"group_id": "Kitchen 1",
	// 	"topic": "Italian"
	// }`)

	client2 := &http.Client{}
	startTime2 := time.Now()
	for i := 0; i < numberOfExperiment; i++ {
		req, err := http.NewRequest("POST", apiURL1, bytes.NewBuffer(requestBody1)) // Produce
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client2.Do(req)
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

		// req2, err := http.NewRequest("POST", apiURL2, bytes.NewBuffer(requestBody2)) // Consume
		// if err != nil {
		// 	fmt.Println("Error creating request:", err)
		// 	return
		// }
		// req2.Header.Set("Content-Type", "application/json")

		// resp2, err := client2.Do(req2)
		// if err != nil {
		// 	fmt.Println("Error making request:", err)
		// 	return
		// }
		// defer resp2.Body.Close()

		// body2, err := io.ReadAll(resp2.Body)
		// if err != nil {
		// 	fmt.Println("Error reading response body:", err)
		// 	return
		// }
		// _ = body2
	}
	elapsedTime2 := time.Since(startTime2)
	fmt.Println("Total time for Kafka:", elapsedTime2)
	ch2 <- "Done"
}