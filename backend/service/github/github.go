package github_bounty

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/err/db"
	"github.com/err/kafka"
	"github.com/google/go-github/v55/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type BountyGithub struct {
	client      *github.Client
	preamble    string
	bountyOrm   *db.BountyORM
	kafkaClient *kafka.BountyKafkaClient
	logger      zerolog.Logger
}

func NewBountyGithubClient(client *github.Client, preamble string, bountyOrm *db.BountyORM, kafkaClient *kafka.BountyKafkaClient) *BountyGithub {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return NewBountyGithubClientWithLogger(client, preamble, bountyOrm, kafkaClient, logger)
}

func NewBountyGithubClientWithLogger(client *github.Client, preamble string, bountyOrm *db.BountyORM, kafkaClient *kafka.BountyKafkaClient, logger zerolog.Logger) *BountyGithub {
	return &BountyGithub{
		client:      client,
		preamble:    preamble,
		bountyOrm:   bountyOrm,
		kafkaClient: kafkaClient,
		logger:      logger,
	}
}

func CreateSigningLink(bountytId, installationId int64, tokenAddress, bountyUIAmount, creatorAddress string) string {
	return fmt.Sprintf("https://app.bounties.network/bounty?bountyId=%d&tokenAddress=%s&bountyUIAmount=%s&creatorAddress=%s&installationId=%d", bountytId, tokenAddress, bountyUIAmount, creatorAddress, installationId)
}

func (b *BountyGithub) UpdateAndCommentIssue(ctx context.Context, issueId int, status, msg string) error {
	bounty, err := b.bountyOrm.GetBountyOnIssueId(ctx, issueId)
	if err != nil {
		return err
	}

	// comment event
	if err = b.CommentEvent(ctx, bounty.RepoOwner, bounty.RepoName, msg, bounty.IssueNumber, b.logger); err != nil {
		b.logger.Error().Err(err).Msg("Failed to comment on pull request")
		return github.ErrBranchNotProtected
	}
	// update bounty status
	if err := b.bountyOrm.UpdateBountyStatus(ctx, bounty.Id, status); err != nil {
		return err
	}
	return nil
}

func (b *BountyGithub) CreateAndCommentIssue(ctx context.Context, event github.IssueCommentEvent, msg string) error {
	repo := event.GetRepo()
	prNum := event.GetIssue().GetNumber()
	instId := githubapp.GetInstallationIDFromEvent(&event)

	ctx, logger := githubapp.PreparePRContext(ctx, instId, repo, prNum)
	logger.Info().Msgf("Event action is %s", event.GetAction())

	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	author := event.GetComment().GetUser().GetLogin()

	if strings.HasSuffix(author, "[bot]") {
		logger.Info().Msg("Issue comment was created by a bot")
		return nil
	}

	logger.Info().Msgf("Echoing comment on %s/%s#%d by %s", repoOwner, repoName, prNum, author)

	userId := event.GetComment().GetUser().GetID()
	userName := event.GetComment().GetUser().GetLogin()
	issueUrl := event.GetIssue().GetURL()
	issueId := event.GetIssue().GetID()
	repoId := repo.GetID()
	ownerId := repo.GetOwner().GetID()

	// try to get wallet address from userId

	// create bounty in db
	_, err := b.bountyOrm.CreateBounty(ctx, db.BountyInput{
		EntityId:    int(userId),
		Url:         issueUrl,
		IssueId:     int(issueId),
		IssueNumber: prNum,
		RepoId:      int(repoId),
		RepoName:    repoName,
		RepoOwner:   repoOwner,
		OwnerId:     int(ownerId),
		Status:      "open",
	}, userName)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create bounty")
		return err
	}

	// send bounty message to github
	if err = b.CommentEvent(ctx, repoOwner, repoName, msg, prNum, logger); err != nil {
		logger.Error().Err(err).Msg("Failed to comment on pull request")
	}

	return nil
}

func (b *BountyGithub) CommentEvent(ctx context.Context, repoOwner, repoName, msg string, prNum int, logger zerolog.Logger) error {
	prComment := github.IssueComment{
		Body: &msg,
	}
	// send bounty message to github
	if _, _, err := b.client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &prComment); err != nil {
		logger.Error().Err(err).Msg("Failed to comment on pull request")
		return err
	}
	return nil
}

// CreateBounty extracts the bounty from the comment and creates a
// bounty message
func (h *BountyGithub) GetNewBountyMessage(ctx context.Context, event github.IssueCommentEvent) (string, error) {
	issueText := event.GetIssue().GetBody()
	issueId := event.GetIssue().GetID()
	author := event.GetComment().GetUser().GetLogin()
	instId := githubapp.GetInstallationIDFromEvent(&event)
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
	signingLink := CreateSigningLink(issueId, instId, "0xaljkdhjkls", amount, "0xkjfksla")
	msg := fmt.Sprintf("In order for the bounty for %s %s to be activated %s please open \n \n %s \n\n to sign the transaction", amount, token, author, signingLink)
	return msg, nil
}
