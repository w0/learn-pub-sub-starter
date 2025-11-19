package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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
	fmt.Println("Starting Peril client...")

	conStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conStr)
	failOnError(err, "Failed to connect to rabbit mq")

	userName, err := gamelogic.ClientWelcome()
	failOnError(err, "Failed to get user name")

	_, _, err = pubsub.DecalreAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+userName, routing.PauseKey, pubsub.TransientQueue)
	failOnError(err, "Failed to declare and bind queue")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	log.Println("os interrupt recieved. shutting down..")
}
