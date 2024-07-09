package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int
const (
    Ack AckType = iota
    NackRequeue
    NackDiscard
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: jsonVal})
	if err != nil {
		return err
	}

	return nil
}

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int, handler func(T) AckType) error {
    channel, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
    if err != nil {
        return err
    }

    delivery, err := channel.Consume(queueName, "", false, false, false, false, nil)
    if err != nil {
        return err
    }

    go func() {
        for {
            val, ok := <-delivery
            if !ok {
                return
            }

            var valT T
            err := json.Unmarshal(val.Body, &valT)
            if err != nil {
                return
            }
            
            ackType := handler(valT)

            if ackType == Ack {
                val.Ack(false)
            } else if ackType == NackRequeue {
                val.Nack(false, true)
            } else {
                val.Nack(false, false) 
            }

        }
    }()
    return nil
}
    
