package pullrequest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Braendie/pull-request-review/internal/models"
	repoerrors "github.com/Braendie/pull-request-review/internal/repository/repo_errors"
	"github.com/lib/pq"
)

type PullRequestRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PullRequestRepository {
	return &PullRequestRepository{
		db: db,
	}
}

func (prr *PullRequestRepository) Create(ctx context.Context, pullRequestID, pullRequestName, authorID string) error {
	const op = "repository.pull_request.Create"

	stmt, err := prr.db.Prepare("INSERT INTO pull_requests(pull_request_id, pull_request_name, author_id, status, created_at) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = stmt.ExecContext(ctx, pullRequestID, pullRequestName, authorID, models.StatusOpen, time.Now())
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code.Name() == repoerrors.ErrAlreadyExistsCode {
			return fmt.Errorf("%s:%w", op, repoerrors.ErrAlreadyExists)
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (prr *PullRequestRepository) List(ctx context.Context, userID string) ([]models.PullRequest, error){
	const op = "repository.pull_request.List"

	stmt, err := prr.db.Prepare(`
	SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at FROM pull_requests AS pr
	JOIN pull_requests_users AS pru ON pru.pull_request_id = pr.pull_request_id
	WHERE pru.user_id = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s:%w", op, repoerrors.ErrNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var pullRequests []models.PullRequest
	for rows.Next() {
		var pullRequest models.PullRequest
		if err := rows.Scan(&pullRequest.ID, &pullRequest.Name, &pullRequest.AuthorID, &pullRequest.Status, &pullRequest.CreatedAt, &pullRequest.MergedAt); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		pullRequests = append(pullRequests, pullRequest)
	}

	return pullRequests, nil

}

func (prr *PullRequestRepository) UpdateStatus(ctx context.Context, pullRequestID string, status models.PullRequestStatus) error {
	const op = "repository.pull_request.UpdateStatus"

	stmt, err := prr.db.Prepare("UPDATE pull_requests SET status = $1 WHERE pull_request_id = $2")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.ExecContext(ctx, status, pullRequestID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s:%w", op, repoerrors.ErrNotFound)
	}

	return nil
}

func (prr *PullRequestRepository) AssignUser(ctx context.Context, userID, pullRequestID string) error {
	const op = "repository.pull_request.AssignUser"

	stmt, err := prr.db.Prepare("INSERT INTO pull_requests_users(user_id, pull_request_id) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = stmt.ExecContext(ctx, userID, pullRequestID)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code.Name() == repoerrors.ErrAlreadyExistsCode {
			return fmt.Errorf("%s:%w", op, repoerrors.ErrAlreadyExists)
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
