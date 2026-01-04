package github

import (
	"context"
	"net/http"
)

// Client defines the GitHub API client interface.
type Client interface {
	// FetchIssue fetches a single Issue with all comments.
	FetchIssue(ctx context.Context, owner, repo string, number int) (*Issue, error)

	// FetchPullRequest fetches a single PR with all comments.
	FetchPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)

	// SetAuthToken sets the authentication token.
	SetAuthToken(token string)

	// SetTimeout sets the request timeout in seconds.
	SetTimeout(timeout int)

	// SetHTTPClient sets a custom HTTP client (for testing).
	SetHTTPClient(client *http.Client)
}

// NewClient creates a new GitHub client instance.
func NewClient() Client {
	return &clientImpl{
		httpClient: http.DefaultClient,
		baseURL:    "https://api.github.com",
	}
}
