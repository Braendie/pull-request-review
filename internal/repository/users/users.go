package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Braendie/pull-request-review/internal/models"
	repoerrors "github.com/Braendie/pull-request-review/internal/repository/repo_errors"
	"github.com/lib/pq"
)

type UsersRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (ur *UsersRepository) Create(ctx context.Context, userID string, username string, teamID int, isActive bool) error {
	const op = "repository.users.Create"

	stmt, err := ur.db.Prepare("INSERT INTO users(user_id, username, team_id, is_active) VALUES($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = stmt.ExecContext(ctx, userID, username, teamID, isActive)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code.Name() == repoerrors.ErrAlreadyExistsCode {
			return fmt.Errorf("%s:%w", op, repoerrors.ErrAlreadyExists)
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

// func (ur *UsersRepository) Get(ctx context.Context, teamID string) (models.User, error) {
// 	const op = "repository.users.Get"

// 	stmt, err := ur.db.Prepare("SELECT user_id, username, is_active FROM users WHERE team_id = $1")
// 	if err != nil {
// 		return models.User{}, fmt.Errorf("%s:%w", op, err)
// 	}

// 	user := models.User{}
// 	err = stmt.QueryRowContext(ctx).Scan(&user.UserID, &user.Username, &user.IsActive)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return models.User{}, fmt.Errorf("%s:%w", op, repoerrors.ErrNotFound)
// 		}

// 		return models.User{}, fmt.Errorf("%s:%w", op, err)
// 	}

// 	return user, nil
// }

func (ur *UsersRepository) List(ctx context.Context, teamID string) ([]models.User, error) {
	const op = "repository.users.List"

	stmt, err := ur.db.Prepare("SELECT user_id, username, is_active FROM users WHERE team_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	rows, err := stmt.QueryContext(ctx, teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.User{}, fmt.Errorf("%s:%w", op, repoerrors.ErrNotFound)
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var users []models.User
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.IsActive); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UsersRepository) UpdateActivity(ctx context.Context, userID string, isActive bool) error {
	const op = "repository.users.UpdateActivity"

	stmt, err := ur.db.Prepare("UPDATE users SET is_active = $1 WHERE user_id = $2")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	res, err := stmt.ExecContext(ctx, isActive, userID)
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
