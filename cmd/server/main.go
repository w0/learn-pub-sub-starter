package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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
	fmt.Println("Starting Peril server...")

	conStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conStr)
	failOnError(err, "Failed to connect to rabbit mq")
	defer conn.Close()

	log.Println("connection successful")

	ch, err := conn.Channel()
	failOnError(err, "Failed to create channel")

	err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect), string(routing.PauseKey),
		routing.PlayingState{
			IsPaused: true,
		})

	failOnError(err, "Failed to publish message")

	log.Println("message published")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	log.Println("os interrupt recieved. shutting down..")

}
