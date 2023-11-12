package main

import (
	"context"
	"fmt"

	github_bounty "github.com/err/github"
	"github.com/google/go-github/v55/github"
)

type Comenter struct {
	prefix             *string
	githubBountyClient github_bounty.BountyGithubI
	comment            *github.IssueComment
	event              *github.IssueCommentEvent
}

func NewComenter(prefix string, event *github.IssueCommentEvent, client github_bounty.BountyGithubI) *Comenter {
	return &Comenter{
		prefix:             &prefix,
		event:              event,
		githubBountyClient: client,
	}
}

func NewBaseCommenter(prefix string) *Comenter {
	return &Comenter{
		prefix: nil,
	}
}

func (c *Comenter) SetComment(comment *github.IssueComment) {
	c.comment = comment
}

func (c *Comenter) SetEvent(event *github.IssueCommentEvent) {
	c.event = event
}

func (c *Comenter) CommentOrUpdate(msg string) (github.IssueComment, error) {
	if c.comment != nil {
		return c.githubBountyClient.UpdateComment(context.Background(), *c.event, c.comment, msg)
	} else {
		return c.githubBountyClient.CommentOnEvent(context.Background(), *c.event, msg)
	}
}

func (c *Comenter) IssueOpened(ctx context.Context) error {
	comment, err := c.CommentOrUpdate(":white_check_mark: Issue comment event action is opened")
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) IssueClosed(ctx context.Context) error {
	comment, err := c.CommentOrUpdate("Issue comment event action is closed")
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) BountyStored(ctx context.Context, postMsg *string) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s \n\n :x: Hmm. We were not able to find the issue with ID=%d that you were looking for.", msg, c.event.GetIssue().GetID())
	if postMsg != nil {
		msg = fmt.Sprintf("%s \n %s", msg, *postMsg)
	}
	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) FailedToFindBounty(ctx context.Context) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s :x: Hmm. We were not able to find the issue with ID=%d that you were looking for.", msg, c.event.GetIssue().GetID())
	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) FailedToCompleteBountyAsRelayer(ctx context.Context) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s :x: Wut?. We're sorry we are not able to close the on chain bounty at this point for issue with ID=%d.", msg, c.event.GetIssue().GetID())
	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) BountyCompletedOnChain(ctx context.Context) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s \n :white_check_mark: Bounty has been closed on chain. \n\n :white_check_mark: %d.", msg, c.event.GetIssue().GetID())
	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) FailedToCreateBounty(ctx context.Context, postText *string) error {

	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s \n Failed to create bounty input from event. Please try again with a new issue.", msg)

	if postText != nil {
		msg = fmt.Sprintf("%s \n %s", msg, *postText)
	}
	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) BountyIsClosed(ctx context.Context) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s :x: %s", msg, "Bounty is already closed and the bounty has been distributed according to the rules.")

	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) BountyIsNotClosed(ctx context.Context) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s :x: %s", msg, "Bounty is not closed yet. Please wait for the owner to close the issue.")

	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) BountyClosable(ctx context.Context) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s \n\n Yes! Lets try to close this bounty and reward some open source contributors! \n\n :white_check_mark: %s", msg, c.event.GetComment().GetBody())

	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}

func (c *Comenter) CommentWithMessage(ctx context.Context, postMsg string) error {
	var msg string
	if c.prefix != nil {
		msg = *c.prefix
	}
	msg = fmt.Sprintf("%s \n\n Yes! Lets try to close this bounty and reward some open source contributors! \n\n :white_check_mark: %s", msg, c.event.GetComment().GetBody())

	msg = fmt.Sprintf("%s \n %s", msg, postMsg)
	comment, err := c.CommentOrUpdate(msg)
	if err != nil {
		return err
	}
	c.comment = &comment
	return nil
}
