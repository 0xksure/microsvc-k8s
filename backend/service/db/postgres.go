// siple package for interacting with the
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/err/common"
	"github.com/err/tokens"
	"github.com/google/go-github/v55/github"
	"github.com/jackc/pgx/v5"
)

type BountyOrm interface {
	CreateBountyCreator(entityId int, username, entityType string) error
	CreateBounty(ctx context.Context, bountyInput BountyInput) (int, error)
	GetBountyOnIssueId(ctx context.Context, issueId int) (Bounty, error)
	UpdateBountyStatus(ctx context.Context, issueId int, status string) error
	CloseIssue(ctx context.Context, issueId int) error
	IsBountyClosed(ctx context.Context, issueId int) (bool, error)
	Close()
}

type BountyORM struct {
	db *pgx.Conn
}

type BountyInput struct {
	Id           int
	EntityId     int
	Url          string
	IssueId      int
	IssueNumber  int
	RepoId       int
	RepoName     string
	RepoOwner    string
	OwnerId      int
	Status       string
	EntityName   string
	BountyAmount int
	BountyToken  string
}

func CreateBountyInputFromEvent(ctx context.Context, event github.IssueCommentEvent, network tokens.Network, rpcUrl string) (BountyInput, error) {
	var bountyInput BountyInput
	githubBounty, err := common.ParseBountyMessage(event.GetIssue().GetBody(), network)
	if err != nil {
		return bountyInput, err
	}
	repo := event.GetRepo()
	prNum := event.GetIssue().GetNumber()
	userId := event.GetSender().GetID()
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	issueUrl := event.GetIssue().GetURL()
	issueId := event.GetIssue().GetID()
	repoId := repo.GetID()
	ownerId := repo.GetOwner().GetID()
	entityName := event.GetSender().GetLogin()

	// validate bountyToken
	if !tokens.IsValidAccount(ctx, githubBounty.Token.Address, rpcUrl) {
		return bountyInput, fmt.Errorf("invalid bounty token address: %s", githubBounty.Token.Address)
	}

	return BountyInput{
		EntityId:     int(userId),
		Url:          issueUrl,
		IssueId:      int(issueId),
		IssueNumber:  prNum,
		RepoId:       int(repoId),
		RepoName:     repoName,
		RepoOwner:    repoOwner,
		OwnerId:      int(ownerId),
		Status:       "open",
		EntityName:   entityName,
		BountyAmount: int(githubBounty.Amount),
		BountyToken:  githubBounty.Token.Address,
	}, nil
}

type Bounty struct {
	BountyInput
	CreatedAt time.Time
	UpdatedAt time.Time
}

// InitBountyOrm initializes the bounty orm
func InitBountyOrm(db *pgx.Conn) BountyORM {
	return BountyORM{db: db}
}

func (b BountyORM) Close() {
	b.db.Close(context.Background())
}

// createBountyCreatorTx
func (b BountyORM) CreateBountyCreator(entityId int, username, entityType string) error {
	_, err := b.db.Exec(context.Background(), `
		INSERT INTO bounty_creator(entity_id,username,entity_type)
		SELECT $1,$2,$3
		WHERE 
			NOT EXISTS (
				SELECT entity_id from bounty_creator where entity_id = $1
			)
	`, entityId, username, entityType)
	return err
}

func (b BountyORM) CreateBounty(ctx context.Context, bountyInput BountyInput) (int, error) {
	tx, err := b.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(ctx, `
	INSERT INTO bounty_creator(entity_id,username,entity_type)
	SELECT $1,$2,$3
	WHERE 
		NOT EXISTS (
			SELECT entity_id from bounty_creator where entity_id = $1
		);
`, bountyInput.EntityId, bountyInput.EntityName, "user")
	if err != nil {
		_ = tx.Rollback(ctx)
		return 0, err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO bounty(
			entity_id,
			url,
			issue_id,
			issue_number,
			repo_id,
			repo_name,
			repo_owner,
			owner_id,
			status,
			amount,
			token_address
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`,
		bountyInput.EntityId,
		bountyInput.Url,
		bountyInput.IssueId,
		bountyInput.IssueNumber,
		bountyInput.RepoId,
		bountyInput.RepoName,
		bountyInput.RepoOwner,
		bountyInput.OwnerId,
		bountyInput.Status,
		bountyInput.BountyAmount,
		bountyInput.BountyToken,
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return 0, err
	}

	return bountyInput.IssueId, tx.Commit(ctx)
}

// getBopunty returns the bounty id for a given issue id
func (b BountyORM) GetBountyOnIssueId(ctx context.Context, issueId int) (Bounty, error) {
	var row Bounty
	err := b.db.QueryRow(ctx,
		`
		SELECT 
			id, 
			entity_id, 
			url, 
			issue_id,
			issue_number,
			repo_id, 
			repo_name, 
			repo_owner,
			owner_id,
			status,
			amount,
			token_address,
			created_at, 
			updated_at 
		FROM bounty WHERE issue_id=$1
		`, issueId).Scan(
		&row.Id,
		&row.EntityId,
		&row.Url,
		&row.IssueId,
		&row.IssueNumber,
		&row.RepoId,
		&row.RepoName,
		&row.RepoOwner,
		&row.OwnerId,
		&row.Status,
		&row.BountyAmount,
		&row.BountyToken,
		&row.CreatedAt,
		&row.UpdatedAt)
	return row, err
}

func (b BountyORM) UpdateBountyStatus(ctx context.Context, issueId int, status string) error {
	_, err := b.db.Exec(ctx,
		`
		UPDATE bounty SET status=$2 WHERE issue_id=$1
		`, issueId, status)
	return err
}

func (b BountyORM) CloseIssue(ctx context.Context, issueId int) error {
	_, err := b.db.Exec(ctx,
		`
		UPDATE bounty SET status='closed' WHERE issue_id=$1
		`, issueId)
	return err
}

func (b BountyORM) IsBountyClosed(ctx context.Context, issueId int) (bool, error) {
	var status string
	err := b.db.QueryRow(ctx,
		`
		SELECT status FROM bounty WHERE issue_id=$1
		`, issueId).Scan(&status)
	if err != nil {
		return false, err
	}
	return status == "closed", nil
}
