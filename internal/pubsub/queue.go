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

	durable, autoDelete, exclusive := getQueueParameters(queueType)
	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	err = ch.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	return ch, queue, nil
}

func getQueueParameters(queueType SimpleQueueType) (bool, bool, bool) {
	switch queueType {
	case DurableQueue:
		return true, false, false
	case TransientQueue:
		return false, true, true
	default:
		return false, false, false
	}
}
