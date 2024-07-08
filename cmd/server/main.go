package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const CONNECT_URL = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	connection, err := amqp.Dial(CONNECT_URL)
	if err != nil {
		fmt.Printf("Connection: %s\n", err)
		return
	}
	defer connection.Close()

	fmt.Println("Successful connected to RabbitMQ server")

	channel, err := connection.Channel()
	if err != nil {
		fmt.Printf("Channel: %s\n", err)
		return
	}

	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		fmt.Println(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Exiting...")

}
