package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"

	_ "github.com/lib/pq"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Integration tests for

func initDocker() *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("could not connect to Docker: %s", err)
	}

	pool.MaxWait = 120 * time.Second
	return pool
}

func createPostgresContainer(dockerPool *dockertest.Pool) *dockertest.Resource {
	container, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine3.18",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=user",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("could not start postgres: %s", err)
	}

	container.Expire(120)
	return container
}

func waitPostgresContainerToBeReady(url string) func() error {
	return func() error {
		var err error
		db, err := sql.Open("postgres", url)
		if err != nil {
			return err
		}

		return db.Ping()
	}
}

func startMigration(databaseUrl string) error {

	_, err := exec.Command("mkdir", "-p", "tmp").Output()
	if err != nil {
		log.Fatalf("could not create tmp: %s", err)
	}

	migrate, err := migrate.New(
		"file://./tmp",
		databaseUrl)
	if err != nil {
		log.Fatalf("could not apply the migration: %s", err)
	}

	return migrate.Up()
}

func TestServerIntegration(t *testing.T) {
	ctx := context.Background()
	dockerPool := initDocker()

	postgresContainer := createPostgresContainer(dockerPool)
	defer dockerPool.Purge(postgresContainer)
	hostAndPort := postgresContainer.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/user?sslmode=disable", hostAndPort)

	if err := dockerPool.Retry(waitPostgresContainerToBeReady(databaseUrl)); err != nil {
		log.Fatalf("postgres container not intialized: %s", err)
	}
	db, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("could not open db connection: %s", err)
	}

	err = startMigration(databaseUrl)
	if err != nil {
		log.Fatalf("could not apply the migration: %s", err)
	}
	bountyOrm := InitBountyOrm(db)
	issueId := 1
	bountyInput := BountyInput{
		EntityId:    1,
		Url:         "url",
		IssueId:     issueId,
		IssueNumber: 2,
		RepoId:      3,
		RepoName:    "repoName",
		RepoOwner:   "repoOwner",
		OwnerId:     4,
		Status:      "open",
		EntityName:  "entityName",
	}
	issueIdOut, err := bountyOrm.CreateBounty(ctx, bountyInput)
	if err != nil {
		t.Errorf("error creating bounty: %s", err)
	}
	if issueId != issueIdOut {
		t.Errorf("issueId not returned: %d", issueIdOut)
	}

	// try to get the bounty
	bounty, err := bountyOrm.GetBountyOnIssueId(ctx, issueId)
	if err != nil {
		t.Errorf("error getting bounty: %s", err)
	}
	if bounty.Status != "open" {
		t.Errorf("status not open: %s", bounty.Status)
	}

	// try to change the status of the bounty
	err = bountyOrm.UpdateBountyStatus(ctx, issueId, "closed")
	if err != nil {
		t.Errorf("error updating bounty status: %s", err)
	}
	bounty, err = bountyOrm.GetBountyOnIssueId(ctx, issueId)
	if err != nil {
		t.Errorf("error getting bounty: %s", err)
	}
	if bounty.Status != "closed" {
		t.Errorf("status not closed: %s", bounty.Status)
	}

}
