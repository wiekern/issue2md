package github

import "time"

// Issue represents a GitHub Issue.
type Issue struct {
	Title         string
	ID            int
	State         string
	URL           string
	Body          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ClosedAt      *time.Time
	Author        string
	AuthorURL     string
	Labels        []string
	Milestone     *string
	Assignees     []string
	CommentsCount int
	Comments      []Comment
}

// PullRequest represents a GitHub Pull Request.
type PullRequest struct {
	Issue
	MergedAt *time.Time
	MergedBy *string
}

// Comment represents a GitHub Issue/PR comment.
type Comment struct {
	Author    string
	AuthorURL string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
