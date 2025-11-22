package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Braendie/pull-request-review/internal/models"
	repoerrors "github.com/Braendie/pull-request-review/internal/repository/repo_errors"
	"github.com/lib/pq"
)

type TeamRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (tr *TeamRepository) Create(ctx context.Context, teamName string) (int64, error) {
	const op = "repository.team.Create"

	stmt, err := tr.db.Prepare("INSERT INTO team(team_name) VALUES($1) RETURNING team_id")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	var id int64
	err = stmt.QueryRowContext(ctx, teamName).Scan(&id)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code.Name() == repoerrors.ErrAlreadyExistsCode {
			return 0, fmt.Errorf("%s:%w", op, repoerrors.ErrAlreadyExists)
		}

		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (tr *TeamRepository) Get(ctx context.Context, teamID int) (models.Team, error) {
	const op = "repository.team.Get"

	stmt, err := tr.db.Prepare("SELECT team_name FROM team WHERE team_id=$1")
	if err != nil {
		return models.Team{}, fmt.Errorf("%s:%w", op, err)
	}

	var teamName string
	err = stmt.QueryRowContext(ctx, teamID).Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Team{}, fmt.Errorf("%s:%w", op, repoerrors.ErrNotFound)
		}

		return models.Team{}, fmt.Errorf("%s:%w", op, err)
	}

	return models.Team{
		ID:   teamID,
		Name: teamName,
	}, nil
}
