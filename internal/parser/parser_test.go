package parser

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *URLInfo
		wantErr bool
	}{
		{
			name: "valid issue URL",
			url:  "https://github.com/golang/go/issues/12345",
			want: &URLInfo{
				Type:     ResourceTypeIssue,
				Owner:    "golang",
				Repo:     "go",
				Number:   12345,
				Original: "https://github.com/golang/go/issues/12345",
			},
			wantErr: false,
		},
		{
			name: "valid pull request URL",
			url:  "https://github.com/golang/go/pull/67890",
			want: &URLInfo{
				Type:     ResourceTypePullRequest,
				Owner:    "golang",
				Repo:     "go",
				Number:   67890,
				Original: "https://github.com/golang/go/pull/67890",
			},
			wantErr: false,
		},
		{
			name: "valid discussion URL",
			url:  "https://github.com/org/community/discussions/123",
			want: &URLInfo{
				Type:     ResourceTypeDiscussion,
				Owner:    "org",
				Repo:     "community",
				Number:   123,
				Original: "https://github.com/org/community/discussions/123",
			},
			wantErr: false,
		},
		{
			name:    "invalid URL - missing protocol",
			url:     "github.com/golang/go/issues/123",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong domain",
			url:     "https://example.com/golang/go/issues/123",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - missing number",
			url:     "https://github.com/golang/go/issues/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - repository homepage",
			url:     "https://github.com/golang/go",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - unsupported resource type",
			url:     "https://github.com/golang/go/releases/123",
			want:    nil,
			wantErr: true,
		},
	}

	p := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("Parse() returned nil, want %+v", tt.want)
					return
				}
				if got.Type != tt.want.Type ||
					got.Owner != tt.want.Owner ||
					got.Repo != tt.want.Repo ||
					got.Number != tt.want.Number ||
					got.Original != tt.want.Original {
					t.Errorf("Parse() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}
