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
	"github.com/err/tokens"
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

type API struct {
	serverAddr    string
	port          string
	logger        zerolog.Logger
	kafkaClient   *kafka.BountyKafkaClient
	clientCreator githubapp.ClientCreator
	ghConfig      *githubapp.Config
	rpcUrl        string
	network       tokens.Network
	bountyOrm     *db.BountyORM
	psqlConn      *pgx.Conn
}

func convertNetworkName(networkName string) (tokens.Network, error) {
	if networkName == "mainnet" {
		return tokens.Mainnet, nil
	} else if networkName == "devnet" {
		return tokens.Devnet, nil
	}
	return tokens.Mainnet, fmt.Errorf("Invalid network name %s", networkName)
}

func InitApi() (API, error) {
	serverAddr := "0.0.0.0"
	port := "8080"
	rpcUrl := os.Getenv("RPC_URL")
	network, err := convertNetworkName(os.Getenv("NETWORK_NAME"))
	if err != nil {
		return API{}, err
	}
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	psqlConnStr := "postgres://user:pwd@ghapp-psql-svc:5432/user?sslmode=disable"
	psqlConn, err := pgx.Connect(context.Background(), psqlConnStr)
	if err != nil {
		logger.Print("db open error: ", err)
		return API{}, err
	}
	bountyOrm := db.InitBountyOrm(psqlConn)

	ghConfig := new(githubapp.Config)
	ghConfig.SetValuesFromEnv("")
	clientCreator, err := initGithubApp(ghConfig)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize github app")
		return API{}, err
	}

	kafkaClient := &kafka.BountyKafkaClient{
		Logger: &logger,
	}
	return API{
		serverAddr:    serverAddr,
		port:          port,
		logger:        logger,
		kafkaClient:   kafkaClient,
		clientCreator: clientCreator,
		ghConfig:      ghConfig,
		rpcUrl:        rpcUrl,
		network:       network,
		bountyOrm:     &bountyOrm,
		psqlConn:      psqlConn,
	}, nil
}

func (api *API) githubWebhook(ghWhErr chan error) {

	prCommentHandler := &PRCommentHandler{
		ClientCreator: api.clientCreator,
		preamble:      "Sandblizzard",
		bountyOrm:     api.bountyOrm,
		kafkaClient:   api.kafkaClient,
		rpcUrl:        api.rpcUrl,
		network:       api.network,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(*api.ghConfig, prCommentHandler)
	h := alice.New().Append(hlog.NewHandler(api.logger)).Then(
		webhookHandler,
	)

	http.Handle(githubapp.DefaultWebhookRoute, h)

	addr := fmt.Sprintf("%s:%s", api.serverAddr, api.port)

	api.logger.Info().Msgf("Starting server on %s:%s", api.serverAddr, api.port)
	ghWhErr <- http.ListenAndServe(addr, nil)
}

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

// github_app is a service that listens for github events
func main() {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	kafkaMessage := make(chan kafka.KafkaMessage)
	ghWhErr := make(chan error)

	// init postgres connection
	api, err := InitApi()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize api")
		return
	}
	defer api.psqlConn.Close(ctx)

	go api.githubWebhook(ghWhErr)
	go api.kafkaClient.GenerateKafkaConsumer(ctx, "bounty", kafkaMessage)

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
		client, err := api.clientCreator.NewInstallationClient(bountyMessage.InstallationId)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to create client for installation %d", bountyMessage.InstallationId)
			return
		}

		githubBountyClient := github_bounty.NewBountyGithubClientWithLogger(client, "Sandblizzard", api.bountyOrm, api.kafkaClient, logger, api.rpcUrl, api.network)
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
