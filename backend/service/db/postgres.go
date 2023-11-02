// siple package for interacting with the
package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type BountyOrm interface {
	CreateBountyCreator(entityId int, username, entityType string) error
	CreateBounty(ctx context.Context, bountyInput BountyInput, entityName string) (int, error)
	GetBountyOnIssueId(ctx context.Context, issueId int) (Bounty, error)
	UpdateBountyStatus(ctx context.Context, issueId int, status string) error
	Close()
}

type BountyORM struct {
	db *pgx.Conn
}

type BountyInput struct {
	Id          int
	EntityId    int
	Url         string
	IssueId     int
	IssueNumber int
	RepoId      int
	RepoName    string
	RepoOwner   string
	OwnerId     int
	Status      string
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

func (b BountyORM) CreateBounty(ctx context.Context, bountyInput BountyInput, entityName string) (int, error) {
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
`, bountyInput.EntityId, entityName, "user")
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
			status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		bountyInput.EntityId,
		bountyInput.Url,
		bountyInput.IssueId,
		bountyInput.IssueNumber,
		bountyInput.RepoId,
		bountyInput.RepoName,
		bountyInput.RepoOwner,
		bountyInput.OwnerId,
		bountyInput.Status)
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
