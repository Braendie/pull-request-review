package models

import "time"

type PullRequest struct {
	ID string
	Name string
	AuthorID string
	Status PullRequestStatus
	CreatedAt time.Time
	MergedAt time.Time
}

type PullRequestStatus string

const (
	StatusOpen   = "OPEN"
	StatusMerged = "MERGED"
)
