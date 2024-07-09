package main

import (
	"fmt"

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

    channel, _, err := pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, fmt.Sprint(routing.PauseKey, ".", username), routing.PauseKey, int(amqp.Transient))
	if err != nil {
		fmt.Println(err)
		return
	}

	state := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, fmt.Sprint(routing.PauseKey, ".", username), routing.PauseKey, int(amqp.Transient), handlerPause(state))
    if err != nil {
        fmt.Println(err)
    }
    
    err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, fmt.Sprint(routing.ArmyMovesPrefix, ".", username), fmt.Sprint(routing.ArmyMovesPrefix, ".*"), int(amqp.Transient), handlerMove(state, channel))
    if err != nil {
        fmt.Println(err)
    }

    
    err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, fmt.Sprint(routing.WarRecognitionsPrefix, ".*"), int(amqp.Persistent), handlerWar(state))
    if err != nil {
        fmt.Println(err)
    }

REPL:
	for {
		input := gamelogic.GetInput()
		if input == nil {
			continue
		}

		switch input[0] {
		case "spawn":
			err = state.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
			continue
		case "move":
            move, err := state.CommandMove(input)
			if err != nil {
				fmt.Println(err)
			}
            err = pubsub.PublishJSON(channel, routing.ExchangePerilTopic, fmt.Sprint(routing.ArmyMovesPrefix, ".", username), move)
            if err != nil {
                fmt.Println(err)
            }
			continue
		case "status":
			state.CommandStatus()
			continue
		case "help":
			gamelogic.PrintClientHelp()
			continue
		case "spam":
			fmt.Println("Spamming is not allowed yet!")
			continue
		case "quit":
			gamelogic.PrintQuit()
			break REPL
		default:
			fmt.Println("Invalid command")
		}
	}
}
