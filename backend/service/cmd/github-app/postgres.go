// siple package for interacting with the
package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type BountyORM struct {
	db *pgx.Conn
}

type Bounty struct {
	Id        int
	EntityId  int
	Url       string
	IssueId   int
	RepoId    int
	RepoName  string
	OwnerId   int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// createBountyCreatorTx
func (b *BountyORM) createBountyCreator(entityId int, username, entityType string) error {
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

func (b *BountyORM) createBounty(ctx context.Context, entityId int, entity_name, url string, issueId int, repoId int, repoName string, ownerId int, status string) (int, error) {
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
`, entityId, entity_name, "user")
	if err != nil {
		_ = tx.Rollback(ctx)
		return 0, err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO bounty(entity_id,url,issue_id,repo_id,repo_name,owner_id,status)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, entityId, url, issueId, repoId, repoName, ownerId, status)
	if err != nil {
		_ = tx.Rollback(ctx)
		return 0, err
	}

	return issueId, tx.Commit(ctx)
}

// getBopunty returns the bounty id for a given issue id
func (b *BountyORM) getBounty(ctx context.Context, issueId int) (Bounty, error) {
	var row Bounty
	err := b.db.QueryRow(ctx,
		`
		SELECT id, entity_id, url, issue_id,repo_id, repo_name, owner_id,status,created_at, updated_at FROM bounty WHERE issue_id=$1
		`, issueId).Scan(&row.Id, &row.EntityId, &row.Url, &row.IssueId, &row.RepoId, &row.RepoName, &row.OwnerId, &row.Status, &row.CreatedAt, &row.UpdatedAt)
	return row, err
}

func (b *BountyORM) updateBountyStatus(ctx context.Context, issueId int, status string) error {
	_, err := b.db.Exec(ctx,
		`
		UPDATE bounty SET status=$2 WHERE issue_id=$1
		`, issueId, status)
	return err
}
