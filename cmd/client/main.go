package main

import (
	"fmt"
	"log"
	"time"

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

	publishCh, err := conn.Channel()
	failOnError(err, "Failed to create channel")

	userName, err := gamelogic.ClientWelcome()
	failOnError(err, "Failed to get user name")

	gs := gamelogic.NewGameState(userName)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.TransientQueue,
		handlerMove(gs, publishCh),
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		string(routing.ExchangePerilTopic),
		string(routing.WarRecognitionsPrefix),
		string(routing.WarRecognitionsPrefix)+".*",
		pubsub.DurableQueue,
		handlerWar(gs, publishCh),
	)

	failOnError(err, "failed subscribing to war outcomes")

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.TransientQueue,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

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
				continue
			}

			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+userName,
				armyMove,
			)

			if err != nil {
				fmt.Printf("error publishing army move: %s\n", err)
			}

			fmt.Printf("Moved %d units to %s", len(armyMove.Units), armyMove.ToLocation)

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

func publishGameLog(publishCh *amqp.Channel, username, msg string) error {
	return pubsub.PublishGob(
		publishCh,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			Username:    username,
			CurrentTime: time.Now(),
			Message:     msg,
		},
	)
}
