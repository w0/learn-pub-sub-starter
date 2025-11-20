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
	fmt.Println("Starting Peril client...")

	conStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conStr)
	failOnError(err, "Failed to connect to rabbit mq")

	userName, err := gamelogic.ClientWelcome()
	failOnError(err, "Failed to get user name")

	_, _, err = pubsub.DecalreAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+userName, routing.PauseKey, pubsub.TransientQueue)
	failOnError(err, "Failed to declare and bind queue")

	gs := gamelogic.NewGameState(userName)

	for {
		input := gamelogic.GetInput()

		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "spawn":
			err = gs.CommandSpawn(input)
			if err != nil {
				fmt.Printf("Failed to spawn unit: %s\n", err)
			}
		case "move":
			armyMove, err := gs.CommandMove(input)
			if err != nil {
				fmt.Printf("Failed to move unit: %s\n", err)
			} else {
				fmt.Printf("Successfully moved to %s\n", armyMove.ToLocation)
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		}
	}

}
