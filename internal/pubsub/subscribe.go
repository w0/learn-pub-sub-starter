package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
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

			handler(*body)

			msg.Ack(false)

		}
	}()

	return nil
}
