package repository

import (
	"database/sql"
	"log"

	pullrequest "github.com/Braendie/pull-request-review/internal/repository/pull_request"
	"github.com/Braendie/pull-request-review/internal/repository/team"
	"github.com/Braendie/pull-request-review/internal/repository/users"
)

type Storage struct {
	Team        *team.TeamRepository
	Users       *users.UsersRepository
	PullRequest *pullrequest.PullRequestRepository
	db          *sql.DB
}

func New(storageCon string) *Storage {
	db, err := sql.Open("postgres", storageCon)
	if err != nil {
		log.Fatalf("Failed to connect to database:%v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database:%v", err)
	}

	return &Storage{
		Team:        team.New(db),
		Users:       users.New(db),
		PullRequest: pullrequest.New(db),
		db:          db,
	}
}
