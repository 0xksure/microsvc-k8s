package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// PR handler defines rules for when
//
// This example is taken and adjusted from https://github.com/palantir/go-githubapp/blob/develop/example/issue_comment.go

type PRCommentHandler struct {
	githubapp.ClientCreator
	preamble  string
	bountyOrm *BountyORM
}

func (h *PRCommentHandler) Handles() []string {
	return []string{"issue_comment", "issues"}
}

func CreateSigningLink(token, amount string) string {
	return fmt.Sprintf("https://app.bounties.network/bounty/%s/%s", token, amount)
}

func (h *PRCommentHandler) CommentIssue(ctx context.Context, event github.IssueCommentEvent, msg string) error {
	repo := event.GetRepo()
	prNum := event.GetIssue().GetNumber()
	instId := githubapp.GetInstallationIDFromEvent(&event)

	ctx, logger := githubapp.PreparePRContext(ctx, instId, repo, prNum)
	logger.Info().Msgf("Event action is %s", event.GetAction())

	client, err := h.NewInstallationClient(instId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create installation client")
		return err
	}

	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	author := event.GetComment().GetUser().GetLogin()

	if strings.HasSuffix(author, "[bot]") {
		logger.Info().Msg("Issue comment was created by a bot")
		return nil
	}

	logger.Info().Msgf("Echoing comment on %s/%s#%d by %s", repoOwner, repoName, prNum, author)
	prComment := github.IssueComment{
		Body: &msg,
	}

	userId := event.GetComment().GetUser().GetID()
	userName := event.GetComment().GetUser().GetLogin()
	issueUrl := event.GetIssue().GetURL()
	issueId := event.GetIssue().GetID()
	repoId := repo.GetID()
	ownerId := repo.GetOwner().GetID()

	// create bounty in db
	_, err = h.bountyOrm.createBounty(ctx, int(userId), userName, issueUrl, int(issueId), int(repoId), repoName, int(ownerId), "open")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create bounty")
		return err
	}

	// send bounty message to github
	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &prComment); err != nil {
		logger.Error().Err(err).Msg("Failed to comment on pull request")
	}

	return nil
}

// CreateBounty extracts the bounty from the comment and
func (h *PRCommentHandler) GetBounty(ctx context.Context, event github.IssueCommentEvent) (string, error) {
	issueText := event.GetIssue().GetBody()
	author := event.GetComment().GetUser().GetLogin()
	// check if bounty is in text
	// if not, return false
	// if yes, create bount
	r := regexp.MustCompile(`\$(\w+:\d+)\$`)
	bounty := r.FindString(issueText)
	if bounty == "" {
		return "", errors.New("No bounty found")
	}
	bountyParts := strings.Split(strings.Trim(bounty, "$"), ":")
	if len(bountyParts) != 2 {
		return "", errors.Errorf("Expected bounty to be two values. Got %v", bounty)
	}
	// token is a string literal e.g. USDC
	token := bountyParts[0]
	// assume amount is in decimals e.g. 100.00
	amount := bountyParts[1]

	// generate signing link
	signingLink := CreateSigningLink(token, amount)
	msg := fmt.Sprintf("In order for the bounty for %s %s to be activated %s please open \n \n %s \n\n to sign the transaction", amount, token, author, signingLink)
	return msg, nil
}

// Handle for PRCommentHandler handles the incoming data when a comment has
// been posted to a PR.
//
// It will echo the comment back to the PR.
func (h *PRCommentHandler) Handle(ctx context.Context, eventType, deliveryId string, payload []byte) error {
	zerolog.Ctx(ctx).Info().Msg("Handling issue comment event")
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to parse the incoming data into an issue comment event")
		return errors.Wrap(err, "failed to parse the incoming data into an issue comment event")
	}

	if event.GetAction() == "opened" {
		zerolog.Ctx(ctx).Info().Msg("Issue comment event action is opened")
		msg, err := h.GetBounty(ctx, event)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("No bounty found")
			return nil
		}
		err = h.CommentIssue(ctx, event, msg)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to comment on issue")
			return err
		}
		return nil
	}

	zerolog.Ctx(ctx).Info().Msg("No action to be made on issue comment")
	return nil
}
