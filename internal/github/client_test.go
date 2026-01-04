package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchIssue(t *testing.T) {
	// Mock GitHub Issue API response
	issueJSON := `{
		"id": 1,
		"number": 12345,
		"title": "Test Issue Title",
		"state": "open",
		"html_url": "https://github.com/testowner/testrepo/issues/12345",
		"body": "This is a test issue body.",
		"created_at": "2024-01-15T10:00:00Z",
		"updated_at": "2024-01-16T11:00:00Z",
		"closed_at": null,
		"user": {
			"login": "testuser",
			"html_url": "https://github.com/testuser"
		},
		"labels": [
			{"name": "bug"},
			{"name": "enhancement"}
		],
		"milestone": {
			"title": "v1.0"
		},
		"assignees": [
			{"login": "assignee1"},
			{"login": "assignee2"}
		],
		"comments": 2
	}`

	commentsJSON := `[
		{
			"user": {
				"login": "commenter1",
				"html_url": "https://github.com/commenter1"
			},
			"body": "First comment",
			"created_at": "2024-01-15T12:00:00Z",
			"updated_at": "2024-01-15T12:00:00Z"
		},
		{
			"user": {
				"login": "commenter2",
				"html_url": "https://github.com/commenter2"
			},
			"body": "Second comment",
			"created_at": "2024-01-15T13:00:00Z",
			"updated_at": "2024-01-15T13:00:00Z"
		}
	]`

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check endpoint path
		if strings.Contains(r.URL.Path, "/comments") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(commentsJSON))
		} else if strings.Contains(r.URL.Path, "/issues/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(issueJSON))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Note: This test will fail because FetchIssue is not implemented yet.
	// We need to modify the client to use a custom HTTP client that points to the mock server.
	// For now, this test documents the expected behavior.

	tests := []struct {
		name    string
		owner   string
		repo    string
		number  int
		want    *Issue
		wantErr bool
	}{
		{
			name:   "fetch valid issue",
			owner:  "testowner",
			repo:   "testrepo",
			number: 12345,
			want: &Issue{
				Title:         "Test Issue Title",
				ID:            12345,
				State:         "open",
				URL:           "https://github.com/testowner/testrepo/issues/12345",
				Body:          "This is a test issue body.",
				CreatedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2024, 1, 16, 11, 0, 0, 0, time.UTC),
				ClosedAt:      nil,
				Author:        "testuser",
				AuthorURL:     "https://github.com/testuser",
				Labels:        []string{"bug", "enhancement"},
				Milestone:     stringPtr("v1.0"),
				Assignees:     []string{"assignee1", "assignee2"},
				CommentsCount: 2,
				Comments: []Comment{
					{
						Author:    "commenter1",
						AuthorURL: "https://github.com/commenter1",
						Body:      "First comment",
						CreatedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
					},
					{
						Author:    "commenter2",
						AuthorURL: "https://github.com/commenter2",
						Body:      "Second comment",
						CreatedAt: time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
	}

	ctx := context.Background()
	client := NewClient()
	client.SetHTTPClient(server.Client())

	// Set baseURL to point to mock server
	c, ok := client.(*clientImpl)
	if !ok {
		t.Fatal("failed to type assert client to *clientImpl")
	}
	c.baseURL = server.URL

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.FetchIssue(ctx, tt.owner, tt.repo, tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				// Compare title
				if got.Title != tt.want.Title {
					t.Errorf("FetchIssue() Title = %v, want %v", got.Title, tt.want.Title)
				}
				// Compare ID
				if got.ID != tt.want.ID {
					t.Errorf("FetchIssue() ID = %v, want %v", got.ID, tt.want.ID)
				}
				// Compare State
				if got.State != tt.want.State {
					t.Errorf("FetchIssue() State = %v, want %v", got.State, tt.want.State)
				}
				// Compare Author
				if got.Author != tt.want.Author {
					t.Errorf("FetchIssue() Author = %v, want %v", got.Author, tt.want.Author)
				}
				// Compare CommentsCount
				if got.CommentsCount != tt.want.CommentsCount {
					t.Errorf("FetchIssue() CommentsCount = %v, want %v", got.CommentsCount, tt.want.CommentsCount)
				}
			}
		})
	}
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
