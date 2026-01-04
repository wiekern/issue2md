package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type clientImpl struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
}

// SetAuthToken sets the authentication token.
func (c *clientImpl) SetAuthToken(token string) {
	c.authToken = token
}

// SetTimeout sets the request timeout in seconds.
func (c *clientImpl) SetTimeout(timeout int) {
	c.httpClient.Timeout = time.Duration(timeout) * time.Second
}

// SetHTTPClient sets a custom HTTP client (for testing).
func (c *clientImpl) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// FetchIssue fetches a single Issue with all comments.
func (c *clientImpl) FetchIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	// Build issue URL
	issueURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d", c.baseURL, owner, repo, number)

	// Fetch issue
	issue, err := c.fetchIssue(ctx, issueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue: %w", err)
	}

	// Fetch comments
	commentsURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", c.baseURL, owner, repo, number)
	comments, err := c.fetchComments(ctx, commentsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	issue.Comments = comments
	issue.CommentsCount = len(comments)

	return issue, nil
}

// FetchPullRequest fetches a single PR with all comments.
func (c *clientImpl) FetchPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// fetchIssue fetches a single issue from GitHub API.
func (c *clientImpl) fetchIssue(ctx context.Context, url string) (*Issue, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Message string `json:"message"`
			Docs    string `json:"documentation_url"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, errResp.Message)
	}

	var apiIssue struct {
		ID        int    `json:"id"`
		Number    int    `json:"number"`
		Title     string `json:"title"`
		State     string `json:"state"`
		HTMLURL   string `json:"html_url"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		ClosedAt  *string `json:"closed_at"`
		User      struct {
			Login    string `json:"login"`
			HTMLURL  string `json:"html_url"`
		} `json:"user"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
		Milestone *struct {
			Title string `json:"title"`
		} `json:"milestone"`
		Assignees []struct {
			Login string `json:"login"`
		} `json:"assignees"`
		Comments int `json:"comments"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiIssue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to Issue
	issue := &Issue{
		ID:        apiIssue.Number,
		Title:     apiIssue.Title,
		State:     apiIssue.State,
		URL:       apiIssue.HTMLURL,
		Body:      apiIssue.Body,
		Author:    apiIssue.User.Login,
		AuthorURL: apiIssue.User.HTMLURL,
	}

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, apiIssue.CreatedAt)
	if err == nil {
		issue.CreatedAt = createdAt
	}
	updatedAt, err := time.Parse(time.RFC3339, apiIssue.UpdatedAt)
	if err == nil {
		issue.UpdatedAt = updatedAt
	}
	if apiIssue.ClosedAt != nil {
		closedAt, err := time.Parse(time.RFC3339, *apiIssue.ClosedAt)
		if err == nil {
			issue.ClosedAt = &closedAt
		}
	}

	// Convert labels
	for _, label := range apiIssue.Labels {
		issue.Labels = append(issue.Labels, label.Name)
	}

	// Convert milestone
	if apiIssue.Milestone != nil {
		issue.Milestone = &apiIssue.Milestone.Title
	}

	// Convert assignees
	for _, assignee := range apiIssue.Assignees {
		issue.Assignees = append(issue.Assignees, assignee.Login)
	}

	issue.CommentsCount = apiIssue.Comments

	return issue, nil
}

// fetchComments fetches comments for an issue.
func (c *clientImpl) fetchComments(ctx context.Context, url string) ([]Comment, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Message string `json:"message"`
			Docs    string `json:"documentation_url"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, errResp.Message)
	}

	var apiComments []struct {
		User struct {
			Login   string `json:"login"`
			HTMLURL string `json:"html_url"`
		} `json:"user"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiComments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var comments []Comment
	for _, apiComment := range apiComments {
		comment := Comment{
			Author:    apiComment.User.Login,
			AuthorURL: apiComment.User.HTMLURL,
			Body:      apiComment.Body,
		}

		createdAt, err := time.Parse(time.RFC3339, apiComment.CreatedAt)
		if err == nil {
			comment.CreatedAt = createdAt
		}
		updatedAt, err := time.Parse(time.RFC3339, apiComment.UpdatedAt)
		if err == nil {
			comment.UpdatedAt = updatedAt
		}

		comments = append(comments, comment)
	}

	return comments, nil
}
