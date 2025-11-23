package pubsub

import "github.com/rabbitmq/amqp091-go"

type SimpleQueueType int

const (
	DurableQueue SimpleQueueType = iota
	TransientQueue
)

func DecalreAndBind(
	conn *amqp091.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp091.Channel, amqp091.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueType == DurableQueue, // durable
		queueType != DurableQueue, // delete when unused
		queueType != DurableQueue, // exclusive
		false,
		amqp091.Table{
			"x-dead-letter-exchange": "peril_dlx",
		})

	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	err = ch.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	return ch, queue, nil
}
