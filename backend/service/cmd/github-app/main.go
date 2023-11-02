package main

import (
	"context"
	"os"

	"github.com/err/kafka"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// github_app is a service that listens for github events
func main() {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// create channels
	kafkaBountyMessage := make(chan kafka.KafkaMessage)
	ghWhErr := make(chan error)

	// init postgres connection
	api, err := InitApi(logger)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize api")
		return
	}
	defer api.bountyOrm.Close()

	go api.githubWebhook(ghWhErr)
	go api.kafkaClient.GenerateKafkaConsumer(ctx, "bounty", kafkaBountyMessage)

	logger.Info().Msg("Waiting for requests")
	select {
	case err := <-ghWhErr:
		logger.Error().Err(err).Msg("github webhook error")
		panic(err)
	case msg := <-kafkaBountyMessage:
		if err := api.BountyKafkaHandler(&msg); err != nil {
			logger.Error().Err(err).Msg("Failed to handle bounty kafka message")
		}
	}
}
