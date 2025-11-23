package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	ch, _, err := DecalreAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		return err
	}

	delivery, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range delivery {
			body := new(T)
			err := json.Unmarshal(msg.Body, body)
			if err != nil {
				log.Printf("Failed to unmarshel msg: %s\n", err)
			}

			msgAck := handler(*body)

			switch msgAck {
			case Ack:
				log.Printf("Message Ack\n")
				msg.Ack(true)
			case NackDiscard:
				log.Printf("Message NackDiscard\n")
				msg.Nack(false, false)
			case NackRequeue:
				log.Printf("Message NackRequeue\n")
				msg.Nack(false, true)
			}

		}
	}()

	return nil
}
