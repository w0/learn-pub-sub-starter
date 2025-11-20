package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conStr)
	failOnError(err, "Failed to connect to rabbit mq")
	defer conn.Close()

	log.Println("Connected to the Peril game server")

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()

		switch input[0] {
		case "pause":
			fmt.Println("Sending pause message.")
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect), string(routing.PauseKey),
				routing.PlayingState{
					IsPaused: true,
				})
			failOnError(err, "Failed to publish message")
		case "resume":
			fmt.Println("Sending resume message")
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect), string(routing.PauseKey),
				routing.PlayingState{
					IsPaused: false,
				})
			failOnError(err, "Failed to publish message")
		case "quit":
			fmt.Println("Server shutting down.")
			return
		default:
			fmt.Println("not impl")
		}
	}

}
