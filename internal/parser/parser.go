package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// ResourceType represents the type of GitHub resource.
type ResourceType string

const (
	ResourceTypeIssue       ResourceType = "issue"
	ResourceTypePullRequest ResourceType = "pull_request"
	ResourceTypeDiscussion ResourceType = "discussion"
)

// URLInfo represents parsed GitHub URL information.
type URLInfo struct {
	Type     ResourceType
	Owner    string
	Repo     string
	Number   int
	Original string
}

// Parser defines the URL parser interface.
type Parser interface {
	Parse(url string) (*URLInfo, error)
	Validate(url string) bool
}

type parser struct{}

// New creates a new Parser instance.
func New() Parser {
	return &parser{}
}

// Parse parses a GitHub URL and returns structured information.
func (p *parser) Parse(url string) (*URLInfo, error) {
	// Check URL prefix
	if !strings.HasPrefix(url, "https://github.com/") {
		return nil, fmt.Errorf("invalid GitHub URL format: %s", url)
	}

	// Remove prefix and split by /
	parts := strings.TrimPrefix(url, "https://github.com/")
	segments := strings.Split(parts, "/")

	// Expected format: owner/repo/type/number
	if len(segments) != 4 {
		return nil, fmt.Errorf("invalid GitHub URL format: %s", url)
	}

	owner := segments[0]
	repo := segments[1]
	resourceType := segments[2]
	numberStr := segments[3]

	// Parse number
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, fmt.Errorf("invalid issue number: %w", err)
	}

	// Map resource type
	var rt ResourceType
	switch resourceType {
	case "issues":
		rt = ResourceTypeIssue
	case "pull":
		rt = ResourceTypePullRequest
	case "discussions":
		rt = ResourceTypeDiscussion
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return &URLInfo{
		Type:     rt,
		Owner:    owner,
		Repo:     repo,
		Number:   number,
		Original: url,
	}, nil
}

// Validate checks if the URL format is valid.
func (p *parser) Validate(url string) bool {
	_, err := p.Parse(url)
	return err == nil
}
