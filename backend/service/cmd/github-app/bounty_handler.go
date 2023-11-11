package main

import (
	"context"
	"fmt"

	"github.com/err/db"
	github_bounty "github.com/err/github"
	"github.com/err/protoc/bounty"
	"github.com/rs/zerolog"
)

// bounty_handler handles the bounty message

type GithubHandler interface {
	Handle(ctx context.Context) error
}

type BountyHandler struct {
	logger             zerolog.Logger
	bountyMessage      *bounty.BountyMessage
	githubBountyClient github_bounty.BountyGithubI
	db                 db.BountyOrm
}

func (b BountyHandler) GenerateSignedMessage() string {
	return fmt.Sprintf(":white_check_mark: Bounty has been activated by %s \n\n :point_right: When the owner has closes the issue the rewards will be distributed amongst the solvers \n\n :bulb: Remember to tag the users who will receive parts of the reward.", b.bountyMessage.CreatorAddress)
}

func (b BountyHandler) GenerateFailedToSignMessage() string {
	return fmt.Sprintf(":x: There was an attempt by %s to create the bounty but it failed \n Please retry the transaction by following the link above :point_up:  ", b.bountyMessage.CreatorAddress)
}

// handle handles the bounty message
// it will based on the status update the bounty
func (b BountyHandler) Handle(ctx context.Context) error {
	switch b.bountyMessage.BountySignStatus {
	case bounty.BountySignStatus_SIGNED:
		// update bounty status to signed
		bountyState, err := b.db.GetBountyOnIssueId(ctx, int(b.bountyMessage.Bountyid))
		if err != nil {
			return err
		}

		// update bounty status to signed
		if err := b.db.UpdateBountyStatus(ctx, bountyState.Id, bounty.BountySignStatus_SIGNED.String()); err != nil {
			return err
		}

		// send message
		msg := b.GenerateSignedMessage()
		if err := b.githubBountyClient.CommentEvent(ctx, bountyState.RepoOwner, bountyState.RepoName, msg, bountyState.IssueNumber, b.logger); err != nil {
			return err
		}
	case bounty.BountySignStatus_FAILED_TO_SIGN:
		bountyState, err := b.db.GetBountyOnIssueId(ctx, int(b.bountyMessage.Bountyid))
		if err != nil {
			return err
		}
		//todo pu lock on the
		msg := b.GenerateFailedToSignMessage()
		if err := b.githubBountyClient.CommentEvent(ctx, bountyState.RepoOwner, bountyState.RepoName, msg, bountyState.IssueNumber, b.logger); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown bounty status %s", b.bountyMessage.BountySignStatus.String())
	}

	return nil
}
