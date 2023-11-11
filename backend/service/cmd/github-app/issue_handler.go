package main

import (
	"context"
	"encoding/json"
	"fmt"

	bounty_program "github.com/err/bounty"
	"github.com/err/common"
	"github.com/err/db"
	github_bounty "github.com/err/github"
	"github.com/err/kafka"
	"github.com/err/tokens"
	"github.com/gagliardetto/solana-go"
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
	bountyOrm     db.BountyOrm
	kafkaClient   kafka.KafkaClient
	rpcUrl        string
	network       tokens.Network
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
	githubBountyClient := github_bounty.NewBountyGithubClientWithLogger(client, h.preamble, h.kafkaClient, *logger, h.rpcUrl, h.network)

	// when issue is opened
	if event.GetAction() == "opened" {
		logger.Info().Msg("Issue comment event action is opened")
		msg, err := githubBountyClient.GetNewBountyMessage(ctx, event)
		if err != nil {
			logger.Err(err).Msg("No bounty found")
			return nil
		}

		bountyInput, err := db.CreateBountyInputFromEvent(ctx, event, h.network, h.rpcUrl)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create bounty input from event")
			return err
		}
		_, err = h.bountyOrm.CreateBounty(ctx, bountyInput)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to comment on issue")
			return err
		}

		err = githubBountyClient.CommentOnEvent(ctx, event, msg)
		return nil
	}

	// when issue is closed
	if event.GetAction() == "created" {
		var msg string
		logger.Info().Msg("Issue comment event action is closed")

		// check if closed
		if event.GetIssue().GetState() != "closed" {
			logger.Info().Msg("Issue is not closed")
			return nil
		}

		isClosed, err := h.bountyOrm.IsBountyClosed(ctx, int(event.GetIssue().GetID()))
		if err != nil {
			logger.Error().Err(err).Msg("Failed to check if bounty is closed")
			return err
		}

		if isClosed {
			logger.Info().Msg("Bounty is already closed")
			return nil
		}

		solverIdentities, err := common.FindIdentitiesFromAtName(event.GetComment().GetBody(), h.rpcUrl, client)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to find identities from at name")
			return err
		}

		if len(solverIdentities) == 0 {
			msg = fmt.Sprintf("No identities found in \n %s", event.GetComment().GetBody())
			githubBountyClient.CommentOnEvent(ctx, event, msg)
			logger.Info().Msg("No identities found")
			return nil
		}

		bountyInfo, err := h.bountyOrm.GetBountyOnIssueId(ctx, int(event.GetIssue().GetID()))
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get bounty on issue id")
			return err
		}

		// try to complete bounty
		bountyMint, err := solana.PublicKeyFromBase58(bountyInfo.BountyToken)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get bounty mint")
			return err
		}
		if err := bounty_program.CompleteBountyAsRelayer(h.rpcUrl, uint64(bountyInfo.IssueId), solverIdentities.GetAddresses(), bountyMint); err != nil {
			logger.Error().Err(err).Msg("Failed to complete bounty as relayer")
			return err
		}

		// update status to completed
		if err := h.bountyOrm.UpdateBountyStatus(ctx, bountyInfo.IssueId, "complete"); err != nil {
			logger.Error().Err(err).Msg("Failed to update bounty status")
			return err
		}

		// Send message to issue
		msg, err = githubBountyClient.GetCloseBountyMessage(ctx, event, solverIdentities)
		if err != nil {
			msg = fmt.Sprintf("Failed to close bounty %s", err.Error())
			githubBountyClient.CommentOnEvent(ctx, event, msg)
			logger.Err(err).Msg("No bounty found")
			return nil
		}
		if err := githubBountyClient.CommentOnEvent(ctx, event, msg); err != nil {
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
