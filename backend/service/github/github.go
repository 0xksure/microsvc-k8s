package github_bounty

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/err/db"
	"github.com/err/identity"
	"github.com/err/kafka"
	"github.com/err/tokens"
	"github.com/google/go-github/v55/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/reiver/go-cast"
	"github.com/rs/zerolog"
)

type BountyGithub struct {
	client      *github.Client
	preamble    string
	bountyOrm   *db.BountyORM
	kafkaClient *kafka.BountyKafkaClient
	logger      zerolog.Logger
	rpcUrl      string
	network     tokens.Network
}

func NewBountyGithubClient(client *github.Client, preamble string, bountyOrm *db.BountyORM, kafkaClient *kafka.BountyKafkaClient, rpcUrl string, network tokens.Network) *BountyGithub {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return NewBountyGithubClientWithLogger(client, preamble, bountyOrm, kafkaClient, logger, rpcUrl, network)
}

func NewBountyGithubClientWithLogger(client *github.Client, preamble string, bountyOrm *db.BountyORM, kafkaClient *kafka.BountyKafkaClient, logger zerolog.Logger, rpcUrl string, network tokens.Network) *BountyGithub {
	return &BountyGithub{
		client:      client,
		preamble:    preamble,
		bountyOrm:   bountyOrm,
		kafkaClient: kafkaClient,
		logger:      logger,
		rpcUrl:      rpcUrl,
		network:     network,
	}
}

func CreateSigningLink(bountytId, installationId int64, tokenAddress, bountyUIAmount, creatorAddress, issueUrl, organization, team, domainType string) string {
	return fmt.Sprintf("https://localhost:3030/bounty?bountyId=%d&tokenAddress=%s&bountyUIAmount=%s&creatorAddress=%s&installationId=%d&referrer=%s&platform=%s&organization=%s&team=%s&domainType=%s", bountytId, tokenAddress, bountyUIAmount, creatorAddress, installationId, issueUrl, "github", organization, team, domainType)
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

func (b *BountyGithub) CommentIssue(ctx context.Context, issueId int, msg string) error {
	bounty, err := b.bountyOrm.GetBountyOnIssueId(ctx, issueId)
	if err != nil {
		return err
	}

	// comment event
	if err = b.CommentEvent(ctx, bounty.RepoOwner, bounty.RepoName, msg, bounty.IssueNumber, b.logger); err != nil {
		b.logger.Error().Err(err).Msg("Failed to comment on pull request")
		return github.ErrBranchNotProtected
	}
	return nil
}

func (b *BountyGithub) CloseAndCommentIssue(ctx context.Context, event github.IssueCommentEvent, msg string) error {
	issueId := event.GetIssue().GetID()
	if err := b.UpdateAndCommentIssue(ctx, int(issueId), "closed", msg); err != nil {
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
func (b *BountyGithub) GetNewBountyMessage(ctx context.Context, event github.IssueCommentEvent) (string, error) {
	issueText := event.GetIssue().GetBody()
	issueId := event.GetIssue().GetID()
	author := event.GetComment().GetUser().GetLogin()
	userId := event.GetComment().GetUser().GetID()
	instId := githubapp.GetInstallationIDFromEvent(&event)
	// check if bounty is in text
	// if not, return false
	// if yes, create bount
	r := regexp.MustCompile(`\$(\w+:\d+)\$`)
	bounty := r.FindString(issueText)
	if bounty == "" {
		return "", errors.New("No bounty found in issueText")
	}
	bountyParts := strings.Split(strings.Trim(bounty, "$"), ":")
	if len(bountyParts) != 2 {
		return "", errors.Errorf("Expected bounty to be two values. Got %v", bounty)
	}
	// token is a string literal e.g. USDC
	tokenSymbol := bountyParts[0]
	token, err := tokens.GetTokenFromSymbol(tokenSymbol, b.network)
	if err != nil {
		return "", err
	}

	// assume amount is in decimals e.g. 100.00
	amount := bountyParts[1]
	userIdu64, err := cast.Uint64(userId)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to cast userId %d to uint64", userId)
	}
	creator, err := identity.GetIdentity(b.rpcUrl, "github", userIdu64)
	if err != nil {
		return "", err
	}

	// generate signing link
	signingLink := CreateSigningLink(issueId, instId, token.Address, amount, creator.Address.String(), *event.GetIssue().URL, *event.GetOrganization().Name, event.Repo.GetFullName(), "issues")
	msg := fmt.Sprintf("In order for the bounty for %s %s to be activated %s please open \n \n :coin: [the bounty link](%s) :coin: \n\n and sign the transaction", amount, token.Address, author, signingLink)
	return msg, nil
}

// GetCloseBountyMessage uses the github event to create a close message
func (b *BountyGithub) GetCloseBountyMessage(ctx context.Context, event github.IssueCommentEvent) (string, error) {
	msg := fmt.Sprintf("Bounty has been closed by %s ", event.GetComment().GetUser().GetLogin())
	return msg, nil
}
