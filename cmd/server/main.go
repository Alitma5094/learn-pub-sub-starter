package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

const CONNECT_URL = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	connection, err := amqp.Dial(CONNECT_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	fmt.Println("Successful connected to RabbitMQ server")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Exiting...")

}
