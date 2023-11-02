package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	github_bounty "github.com/err/github"
	"github.com/err/protoc/bounty"
	"github.com/google/go-github/v55/github"
	"github.com/rs/zerolog"
)

var (
	Client http.Client
)

type ClientMock struct{}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MBountyGithub struct{}

func (m MBountyGithub) UpdateAndCommentIssue(ctx context.Context, issueId int, status, message string) error {
	return nil
}

func (m MBountyGithub) CommentIssue(ctx context.Context, issueId int, message string) error {
	return nil
}
func (m MBountyGithub) CloseAndCommentIssue(ctx context.Context, event github.IssueCommentEvent, msg string) error {
	return nil
}

func (m MBountyGithub) CreateAndCommentIssue(ctx context.Context, event github.IssueCommentEvent, msg string) error {
	return nil
}

// CommentEvent will always return nil
func (m MBountyGithub) CommentEvent(ctx context.Context, repoOwner, repoName, msg string, prNum int, logger zerolog.Logger) error {
	return nil
}
func (m MBountyGithub) GetNewBountyMessage(ctx context.Context, event github.IssueCommentEvent) (string, error) {
	return "", nil
}
func (m MBountyGithub) GetCloseBountyMessage(ctx context.Context, event github.IssueCommentEvent) (string, error) {
	return "", nil
}

func generateMockGithubClient(testServer *httptest.Server) *github.Client {
	testClient := testServer.Client()
	url, err := url.Parse(testServer.URL + "/")
	if err != nil {
		panic(err)
	}
	githubClient := github.NewClient(testClient)
	githubClient.BaseURL = url
	return githubClient
}

func generateBountyHandlerMock(bountyMessage *bounty.BountyMessage, testServer *httptest.Server) BountyHandler {
	bountyOrmMock := MBountyOrm{}
	kafkaClientMock := MKafkaClient{}
	githubClient := generateMockGithubClient(testServer)
	bountyGithub := github_bounty.NewBountyGithubClient(githubClient, "", bountyOrmMock, kafkaClientMock, "", "mainnet")

	return BountyHandler{
		bountyMessage:      bountyMessage,
		githubBountyClient: bountyGithub,
	}
}

func TestBountyHandler(t *testing.T) {
	t.Log("Testing bounty handler")

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("New request: ", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1}`))
	}))
	defer func() { testServer.Close() }()

	t.Run("Test Bounty handler with SIGNED", func(t *testing.T) {
		t.Log("Testing handler")
		bountyhandler := generateBountyHandlerMock(&bounty.BountyMessage{
			Bountyid:         1,
			BountySignStatus: bounty.BountySignStatus_SIGNED,
		}, testServer)
		err := bountyhandler.Handle(context.Background())
		if err != nil {
			t.Errorf("Expected nil, got %s", err)
		}
	})

	t.Run("Test Bounty handler with FAILED_TO_SIGN -> should succeed", func(t *testing.T) {
		bountyhandler := generateBountyHandlerMock(&bounty.BountyMessage{
			Bountyid:         1,
			BountySignStatus: bounty.BountySignStatus_FAILED_TO_SIGN,
		}, testServer)
		err := bountyhandler.Handle(context.Background())
		if err != nil {
			t.Errorf("Expected nil, got %s", err)
		}
	})

	t.Run("Test Bounty handler with not supported status -> should succeed", func(t *testing.T) {
		bountyhandler := generateBountyHandlerMock(&bounty.BountyMessage{
			Bountyid:         1,
			BountySignStatus: bounty.BountySignStatus_CREATED,
		}, testServer)
		err := bountyhandler.Handle(context.Background())
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
