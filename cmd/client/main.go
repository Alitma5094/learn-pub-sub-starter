package main

import (
	"fmt"
    "os"
    "os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)


const CONNECT_URL = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril client...")

	connection, err := amqp.Dial(CONNECT_URL)
	if err != nil {
		fmt.Printf("Connection: %s\n", err)
		return
	}
	defer connection.Close()

    username, err := gamelogic.ClientWelcome()
    if err != nil {
        fmt.Println(err)
        return 
    }
    
    _, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, fmt.Sprint(routing.PauseKey, ".", username), routing.PauseKey, int(amqp.Transient))
    if err != nil {
        fmt.Println(err)
        return 
    }

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

}
