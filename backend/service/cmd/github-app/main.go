package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/err/db"
	github_bounty "github.com/err/github"
	"github.com/err/kafka"
	"github.com/err/protoc/bounty"
	"github.com/golang/protobuf/proto"
	"github.com/gregjones/httpcache"
	"github.com/jackc/pgx/v5"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func trackingMiddleware(next http.Handler, logger zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info().Msgf("Hitting webhook endpoint: %s", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func initGithubApp(ghConfig *githubapp.Config) (githubapp.ClientCreator, error) {
	metricsRegistry := metrics.DefaultRegistry

	return githubapp.NewDefaultCachingClientCreator(
		*ghConfig,
		githubapp.WithClientUserAgent("github-app"),
		githubapp.WithClientTimeout(10*time.Second),
		githubapp.WithClientCaching(true, func() httpcache.Cache {
			return httpcache.NewMemoryCache()
		}),
		githubapp.WithClientMiddleware(githubapp.ClientMetrics(metricsRegistry), githubapp.ClientLogging(zerolog.DebugLevel)),
	)

}

func githubWebhook(serverAddr string, port string, logger zerolog.Logger, bountyORM *db.BountyORM, kafkaClient *kafka.BountyKafkaClient, clientCreator githubapp.ClientCreator, ghConfig *githubapp.Config, ghWhErr chan error) {

	prCommentHandler := &PRCommentHandler{
		ClientCreator: clientCreator,
		preamble:      "Sandblizzard",
		bountyOrm:     bountyORM,
		kafkaClient:   kafkaClient,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(*ghConfig, prCommentHandler)
	h := alice.New().Append(hlog.NewHandler(logger)).Then(
		webhookHandler,
	)

	http.Handle(githubapp.DefaultWebhookRoute, h)

	addr := fmt.Sprintf("%s:%s", serverAddr, port)

	logger.Info().Msgf("Starting server on %s:%s", serverAddr, port)
	ghWhErr <- http.ListenAndServe(addr, nil)
}

// github_app is a service that listens for github events
func main() {
	ctx := context.Background()
	serverAddr := "0.0.0.0"
	port := "8080"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	kafkaMessage := make(chan kafka.KafkaMessage)
	ghWhErr := make(chan error)
	ghConfig := new(githubapp.Config)
	ghConfig.SetValuesFromEnv("")
	cc, err := initGithubApp(ghConfig)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize github app")
		return
	}
	// init postgres connection
	psqlConnStr := "postgres://user:pwd@ghapp-psql-svc:5432/user?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), psqlConnStr)
	if err != nil {
		logger.Print("db open error: ", err)
		return
	}
	defer conn.Close(context.Background())

	bountyORM := db.InitBountyOrm(conn)

	kafkaBountyClient := &kafka.BountyKafkaClient{
		Logger: &logger,
	}

	go githubWebhook(serverAddr, port, logger, &bountyORM, kafkaBountyClient, cc, ghConfig, ghWhErr)
	go kafkaBountyClient.GenerateKafkaConsumer(ctx, "bounty", kafkaMessage)

	logger.Info().Msg("Waiting for requests")
	select {
	case err := <-ghWhErr:
		logger.Error().Err(err).Msg("github webhook error")
		panic(err)
	case msg := <-kafkaMessage:
		ctx := context.Background()
		var bountyMessage bounty.BountyMessage
		err := proto.Unmarshal(msg.Msg, &bountyMessage)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to decode bounty message %v", msg.Msg)
			return
		}

		logger.Info().Msgf("Received bounty message %v", &bountyMessage)

		// check bounty message
		client, err := cc.NewInstallationClient(bountyMessage.InstallationId)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to create client for installation %d", bountyMessage.InstallationId)
			return
		}

		githubBountyClient := github_bounty.NewBountyGithubClientWithLogger(client, "Sandblizzard", &bountyORM, kafkaBountyClient, logger)
		bountyHandler := &BountyHandler{
			bountyMessage:      &bountyMessage,
			githubBountyClient: githubBountyClient,
		}
		if err := bountyHandler.Handle(ctx); err != nil {
			logger.Error().Err(err).Msgf("Failed to handle bounty message %v", bountyMessage)
			return
		}

	}
}
