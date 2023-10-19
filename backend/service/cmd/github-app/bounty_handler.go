package main

import (
	"context"
	"fmt"

	github_bounty "github.com/err/github"
	"github.com/err/protoc/bounty"
)

// bounty_handler handles the bounty message

type BountyHandler struct {
	bountyMessage      *bounty.BountyMessage
	githubBountyClient *github_bounty.BountyGithub
}

func (b *BountyHandler) GenerateSignedMessage() string {
	return fmt.Sprintf("Bounty has been activated by %s ", b.bountyMessage.CreatorAddress)
}

// handle handles the bounty message
// it will based on the status update the bounty
func (b *BountyHandler) Handle(ctx context.Context) error {
	switch b.bountyMessage.BountySignStatus {
	case bounty.BountySignStatus_SIGNED:
		// update bounty status to signed
		if err := b.githubBountyClient.UpdateAndCommentIssue(ctx, int(b.bountyMessage.Bountyid), bounty.BountySignStatus_SIGNED.String(), b.GenerateSignedMessage()); err != nil {
			return err
		}
	}

	return nil
}
