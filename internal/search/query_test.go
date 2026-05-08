package search

import (
	"testing"
	"time"
)

func TestParseRelativeDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    string
		wantDays int
		wantErr  bool
	}{
		{"empty string", "", 0, false},
		{"1d", "1d", 1, false},
		{"7d", "7d", 7, false},
		{"30d", "30d", 30, false},
		{"100d", "100d", 100, false},
		{"absolute date", "2024-01-15", 0, false},
		{"invalid relative", "1w", 0, true},
		{"invalid format", "abc", 0, true},
		{"empty number", "0d", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRelativeDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRelativeDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.input == "" {
					if !got.IsZero() {
						t.Errorf("ParseRelativeDate(%q) = %v, want zero time", tt.input, got)
					}
				} else if tt.wantDays > 0 {
					expected := now.AddDate(0, 0, -tt.wantDays)
					diff := got.Sub(expected)
					if diff < 0 {
						diff = -diff
					}
					if diff > time.Second {
						t.Errorf("ParseRelativeDate(%q) = %v, want approximately %v", tt.input, got, expected)
					}
				}
			}
		})
	}
}

func TestCommitSearchQuery_Validate(t *testing.T) {
	futureDate := time.Now().AddDate(0, 0, 10)
	pastDate := time.Now().AddDate(0, 0, -10)

	tests := []struct {
		name    string
		query   CommitSearchQuery
		wantErr error
	}{
		{"valid query", CommitSearchQuery{Q: "fix bug"}, nil},
		{"valid query with all fields", CommitSearchQuery{Q: "fix", Author: "Alice", Repo: "/repo", Limit: 50}, nil},
		{"valid query with since/until", CommitSearchQuery{Q: "fix", Since: &pastDate, Until: &futureDate}, nil},
		{"empty query", CommitSearchQuery{}, ErrEmptyQuery},
		{"since after until", CommitSearchQuery{Q: "fix", Since: &futureDate, Until: &pastDate}, ErrInvalidDateRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if err != tt.wantErr {
				t.Errorf("CommitSearchQuery.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommitSearchQuery_SetDefaults(t *testing.T) {
	q := CommitSearchQuery{}
	q.SetDefaults()
	if q.Limit != 50 {
		t.Errorf("expected default limit 50, got %d", q.Limit)
	}

	q = CommitSearchQuery{Limit: 0}
	q.SetDefaults()
	if q.Limit != 50 {
		t.Errorf("expected default limit 50 for zero, got %d", q.Limit)
	}

	q = CommitSearchQuery{Limit: 500}
	q.SetDefaults()
	if q.Limit != 200 {
		t.Errorf("expected capped limit 200, got %d", q.Limit)
	}

	q = CommitSearchQuery{Offset: -5}
	q.SetDefaults()
	if q.Offset != 0 {
		t.Errorf("expected zero offset for negative, got %d", q.Offset)
	}
}