package main

import (
	"context"
	"encoding/json"

	"github.com/err/db"
	github_bounty "github.com/err/github"
	"github.com/err/kafka"
	"github.com/google/go-github/v55/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// PR handler defines rules for when
//
// This example is taken and adjusted from https://github.com/palantir/go-githubapp/blob/develop/example/issue_comment.go

type PRCommentHandler struct {
	ClientCreator githubapp.ClientCreator
	preamble      string
	bountyOrm     *db.BountyORM
	kafkaClient   *kafka.BountyKafkaClient
}

func (h *PRCommentHandler) Handles() []string {
	return []string{"issue_comment", "issues"}
}

// Handle for PRCommentHandler handles the incoming data when a comment has
// been posted to a PR.
//
// It will echo the comment back to the PR.
func (h *PRCommentHandler) Handle(ctx context.Context, eventType, deliveryId string, payload []byte) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Handling issue comment event")
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		logger.Error().Err(err).Msg("Failed to parse the incoming data into an issue comment event")
		return errors.Wrap(err, "failed to parse the incoming data into an issue comment event")
	}

	instId := githubapp.GetInstallationIDFromEvent(&event)
	client, err := h.ClientCreator.NewInstallationClient(instId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create installation client")
		return errors.Wrap(err, "failed to create installation client")
	}
	githubBountyClient := github_bounty.NewBountyGithubClientWithLogger(client, h.preamble, h.bountyOrm, h.kafkaClient, *logger)

	// when issue is opened
	if event.GetAction() == "opened" {
		logger.Info().Msg("Issue comment event action is opened")
		msg, err := githubBountyClient.GetNewBountyMessage(ctx, event)
		if err != nil {
			logger.Err(err).Msg("No bounty found")
			return nil
		}
		err = githubBountyClient.CreateAndCommentIssue(ctx, event, msg)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to comment on issue")
			return err
		}
		return nil
	}

	// when issue is closed
	if event.GetAction() == "closed" {
		logger.Info().Msg("Issue comment event action is closed")

		msg, err := githubBountyClient.GetCloseBountyMessage(ctx, event)
		if err != nil {
			logger.Err(err).Msg("No bounty found")
			return nil
		}

		if err := githubBountyClient.CloseAndCommentIssue(ctx, event, msg); err != nil {
			logger.Error().Err(err).Msg("Failed to comment on issue")
			return err
		}
		return nil
	}

	// when issue is commented
	if event.GetAction() == "created" {
		logger.Info().Msg("Issue is commented")
	}

	logger.Info().Msg("No action to be made on issue comment")
	return nil
}
