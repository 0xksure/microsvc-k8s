package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

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

// github_app is a service that listens for github events
func main() {
	serverAddr := "0.0.0.0"
	port := "8080"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ghConfig := new(githubapp.Config)
	ghConfig.SetValuesFromEnv("")

	cc, err := initGithubApp(ghConfig)
	if err != nil {
		panic(err)
	}

	// init postgres connection
	psqlConnStr := "postgres://user:pwd@ghapp-psql-svc:5432/user?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), psqlConnStr)
	if err != nil {
		logger.Print("db open error: ", err)
		return
	}
	defer conn.Close(context.Background())

	BountyORM := &BountyORM{
		db: conn,
	}
	prCommentHandler := &PRCommentHandler{
		ClientCreator: cc,
		preamble:      "Sandblizzard",
		bountyOrm:     BountyORM,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(*ghConfig, prCommentHandler)
	h := alice.New().Append(hlog.NewHandler(logger)).Then(
		webhookHandler,
	)

	http.Handle(githubapp.DefaultWebhookRoute, h)

	addr := fmt.Sprintf("%s:%s", serverAddr, port)

	logger.Info().Msgf("Starting server on %s:%s", serverAddr, port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
