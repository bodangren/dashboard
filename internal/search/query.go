package search

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrEmptyQuery          = errors.New("query text cannot be empty")
	ErrInvalidDateFormat   = errors.New("invalid date format, use YYYY-MM-DD or relative like 1d, 7d, 30d")
	ErrInvalidRelativeDate = errors.New("invalid relative date format, use <number>d (e.g., 1d, 7d, 30d)")
	ErrInvalidDateRange    = errors.New("since date must be before until date")
)

var relativeDateRegex = regexp.MustCompile(`^(\d+)d$`)

type CommitSearchQuery struct {
	Q      string
	Author string
	Repo   string
	Since  *time.Time
	Until  *time.Time
	Limit  int
	Offset int
}

func ParseRelativeDate(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, nil
	}

	if matches := relativeDateRegex.FindStringSubmatch(input); matches != nil {
		days, err := parseDays(matches[1])
		if err != nil {
			return time.Time{}, ErrInvalidRelativeDate
		}
		return time.Now().AddDate(0, 0, -days), nil
	}

	return time.Parse("2006-01-02", input)
}

func parseDays(s string) (int, error) {
	var days int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("invalid number")
		}
		days = days*10 + int(c-'0')
	}
	if days < 0 {
		return 0, errors.New("invalid number")
	}
	return days, nil
}

func (q *CommitSearchQuery) Validate() error {
	if q.Q == "" {
		return ErrEmptyQuery
	}
	if q.Since != nil && q.Until != nil && q.Since.After(*q.Until) {
		return ErrInvalidDateRange
	}
	return nil
}

func (q *CommitSearchQuery) SetDefaults() {
	if q.Limit <= 0 {
		q.Limit = 50
	}
	if q.Limit > 200 {
		q.Limit = 200
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
}
