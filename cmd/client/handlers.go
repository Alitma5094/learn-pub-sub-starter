package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
   
    return func(ps routing.PlayingState) pubsub.AckType {
        defer fmt.Print("> ")
        gs.HandlePause(ps)
        return pubsub.Ack
    }
}

func handlerMove(gs *gamelogic.GameState, channel *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {

    return func(move gamelogic.ArmyMove) pubsub.AckType {
        defer fmt.Print("> ")
        outcome := gs.HandleMove(move)

        if outcome == gamelogic.MoveOutComeSafe{
            fmt.Println("DEBUG: ack")
            return pubsub.Ack
        }

        if outcome == gamelogic.MoveOutcomeMakeWar {
            fmt.Println("DEBUG: ack")
            err := pubsub.PublishJSON(channel, routing.ExchangePerilTopic, fmt.Sprint(routing.WarRecognitionsPrefix, ".", gs.Player.Username), gamelogic.RecognitionOfWar{Attacker: move.Player, Defender: gs.Player})
            if err != nil {
                return pubsub.NackRequeue
            }
            return pubsub.Ack
        }
        
        fmt.Println("DEBUG: nack disq")
        return pubsub.NackDiscard
    }
}


func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
   
    return func(war gamelogic.RecognitionOfWar) pubsub.AckType {
        defer fmt.Print("> ")
        outcome, _, _ := gs.HandleWar(war)

        switch outcome {
        case gamelogic.WarOutcomeNotInvolved:
            return pubsub.NackRequeue

        case gamelogic.WarOutcomeNoUnits:
            return pubsub.NackDiscard

        case gamelogic.WarOutcomeOpponentWon:
        case gamelogic.WarOutcomeYouWon:
        case gamelogic.WarOutcomeDraw:
            return pubsub.Ack

        default:
            fmt.Println("ERROR: Invalid war outcome")
            return pubsub.NackDiscard
        }
        return pubsub.NackDiscard
    }
}
