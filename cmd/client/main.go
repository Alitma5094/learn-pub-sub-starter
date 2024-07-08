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
    
    _, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, fmt.Sprint(routing.PauseKey, ".", username), routing.PauseKey, int(amqp.Transient))
    if err != nil {
        fmt.Println(err)
        return 
    }

    state := gamelogic.NewGameState(username)

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
           _, err = state.CommandMove(input)
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
