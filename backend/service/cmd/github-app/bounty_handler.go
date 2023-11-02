package main

import (
	"context"
	"fmt"

	github_bounty "github.com/err/github"
	"github.com/err/protoc/bounty"
)

// bounty_handler handles the bounty message

type GithubHandler interface {
	Handle(ctx context.Context) error
}

type BountyHandler struct {
	bountyMessage      *bounty.BountyMessage
	githubBountyClient github_bounty.BountyGithubI
}

func (b BountyHandler) GenerateSignedMessage() string {
	return fmt.Sprintf("Bounty has been activated by %s ", b.bountyMessage.CreatorAddress)
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
		if err := b.githubBountyClient.UpdateAndCommentIssue(ctx, int(b.bountyMessage.Bountyid), bounty.BountySignStatus_SIGNED.String(), b.GenerateSignedMessage()); err != nil {
			return err
		}
	case bounty.BountySignStatus_FAILED_TO_SIGN:
		//todo pu lock on the
		if err := b.githubBountyClient.CommentIssue(ctx, int(b.bountyMessage.Bountyid), b.GenerateFailedToSignMessage()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown bounty status %s", b.bountyMessage.BountySignStatus.String())
	}

	return nil
}
