package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int) (*amqp.Channel, amqp.Queue, error) {
    channel, err := conn.Channel()
    if err != nil {
        return nil, amqp.Queue{}, err
    }
    
    queue, err := channel.QueueDeclare(queueName, simpleQueueType == int(amqp.Persistent), simpleQueueType == int(amqp.Transient), simpleQueueType == int(amqp.Transient), false, nil)
    if err != nil {
        return nil, amqp.Queue{}, err
    }

    err = channel.QueueBind(queueName, key, exchange, false, nil)
    if err != nil {
        return nil, amqp.Queue{}, err
    }

    return channel, queue, nil
}
