package github_bounty

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/err/common"
	"github.com/err/identity"
	"github.com/err/kafka"
	"github.com/err/tokens"
	"github.com/gagliardetto/solana-go"
	"github.com/google/go-github/v55/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type BountyGithubI interface {
	CommentOnEvent(ctx context.Context, event github.IssueCommentEvent, msg string) (github.IssueComment, error)
	UpdateComment(ctx context.Context, event github.IssueCommentEvent, comment *github.IssueComment, msg string) (github.IssueComment, error)
	CommentEvent(ctx context.Context, repoOwner, repoName, msg string, prNum int, logger zerolog.Logger) error
	GetNewBountyMessage(ctx context.Context, event github.IssueCommentEvent) (string, error)
	GetCloseBountyMessage(ctx context.Context, event github.IssueCommentEvent, solverIdentities []identity.Identity, signature solana.Signature) (string, error)
}

type BountyGithub struct {
	client      *github.Client
	preamble    string
	kafkaClient kafka.KafkaClient
	logger      zerolog.Logger
	rpcUrl      string
	network     tokens.Network
}

func NewBountyGithubClient(client *github.Client, preamble string, kafkaClient kafka.KafkaClient, rpcUrl string, network tokens.Network) BountyGithubI {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return NewBountyGithubClientWithLogger(client, preamble, kafkaClient, logger, rpcUrl, network)
}

func NewBountyGithubClientWithLogger(client *github.Client, preamble string, kafkaClient kafka.KafkaClient, logger zerolog.Logger, rpcUrl string, network tokens.Network) BountyGithubI {
	return BountyGithub{
		client:      client,
		preamble:    preamble,
		kafkaClient: kafkaClient,
		logger:      logger,
		rpcUrl:      rpcUrl,
		network:     network,
	}
}

func (b BountyGithub) CommentOnEvent(ctx context.Context, event github.IssueCommentEvent, msg string) (github.IssueComment, error) {
	comment, response, err := b.client.Issues.CreateComment(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), event.GetIssue().GetNumber(), &github.IssueComment{
		Body: &msg,
	})
	if err != nil {
		return *comment, err
	}

	if response.StatusCode != 201 {
		return *comment, errors.Errorf("Expected status code 201. Got %v", response.StatusCode)
	}
	return *comment, err
}

func (b BountyGithub) UpdateComment(ctx context.Context, event github.IssueCommentEvent, comment *github.IssueComment, msg string) (github.IssueComment, error) {
	newComment := fmt.Sprintf("%s \n\n :white_check_mark: %s", comment.GetBody(), msg)
	comment.Body = &newComment
	comment, response, err := b.client.Issues.EditComment(ctx, event.GetRepo().GetOwner().GetLogin(), event.GetRepo().GetName(), comment.GetID(), comment)
	if err != nil {
		return *comment, err
	}
	if response.StatusCode != 200 {
		return *comment, errors.Errorf("Expected status code 200. Got %v", response.StatusCode)
	}
	return *comment, err
}
func (b BountyGithub) CommentEvent(ctx context.Context, repoOwner, repoName, msg string, prNum int, logger zerolog.Logger) error {
	prComment := github.IssueComment{
		Body: &msg,
	}

	if _, _, err := b.client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &prComment); err != nil {
		logger.Error().Err(err).Msg("Failed to comment on pull request")
		return err
	}
	return nil
}

// CreateBounty extracts the bounty from the comment and creates a
// bounty message
func (b BountyGithub) GetNewBountyMessage(ctx context.Context, event github.IssueCommentEvent) (string, error) {
	issueText := event.GetIssue().GetBody()
	issueId := event.GetIssue().GetID()
	issueUrl := event.GetIssue().GetURL()
	author := event.GetSender().GetLogin()
	userId := event.GetSender().GetID()
	organizationName := event.GetRepo().GetOwner().GetLogin()
	repoName := event.GetRepo().GetName()
	instId := githubapp.GetInstallationIDFromEvent(&event)

	githubBounty, err := common.ParseBountyMessage(issueText, b.network)
	if err != nil {
		return "", err
	}

	// assume amount is in decimals e.g. 100.00
	if userId < 1 {
		return "", errors.Errorf("Expected userId to be greater than 0. Got %v", userId)
	}
	userIdu64 := uint64(userId)

	creator, err := identity.GetIdentity(b.rpcUrl, "github", userIdu64)
	if err != nil {
		return "", err
	}

	// generate signing link
	signingLink := common.CreateSigningLink(issueId, instId, githubBounty.Token.Address, strconv.Itoa(int(githubBounty.Amount)), creator.Address.String(), issueUrl, organizationName, repoName, "issues")
	msg := fmt.Sprintf("In order for the bounty for %s %s to be activated @%s please open \n \n :coin: [the bounty link](%s) :coin: \n\n and sign the transaction", githubBounty.AmountUI, githubBounty.Token.Address, author, signingLink)
	return msg, nil
}

// GetCloseBountyMessage uses the github event to create a close message
func (b BountyGithub) GetCloseBountyMessage(ctx context.Context, event github.IssueCommentEvent, solverIdentities []identity.Identity, signature solana.Signature) (string, error) {
	prettySolvers := ""
	for _, identity := range solverIdentities {
		prettySolvers += fmt.Sprintf("%s ", identity.Address)
	}
	explorerLink := common.GetExplorerLink(b.network, signature)
	msg := fmt.Sprintf("Bounty has been closed by %s \n :globe_with_meridians: %s", prettySolvers, explorerLink)

	return msg, nil
}
