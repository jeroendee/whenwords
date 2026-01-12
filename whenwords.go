// Package whenwords provides human-friendly time formatting and parsing.
package whenwords

import "errors"

// Error variables for common error conditions.
var (
	ErrNegativeDuration = errors.New("duration cannot be negative")
	ErrEmptyInput       = errors.New("input cannot be empty")
	ErrUnparseable      = errors.New("unable to parse duration")
	ErrNegativeValue    = errors.New("negative values are not allowed")
)

// DurationOption configures Duration output.
type DurationOption func(*durationConfig)

// durationConfig holds options for Duration formatting.
type durationConfig struct {
	compact  bool
	maxUnits int
}

// WithCompact enables compact output format (e.g., "2h 30m" instead of "2 hours, 30 minutes").
func WithCompact() DurationOption {
	return func(c *durationConfig) {
		c.compact = true
	}
}

// WithMaxUnits limits the number of time units in the output.
func WithMaxUnits(n int) DurationOption {
	return func(c *durationConfig) {
		c.maxUnits = n
	}
}

// TimeAgo returns a human-readable relative time string.
// The optional reference parameter defaults to the timestamp itself (returns "just now").
func TimeAgo(timestamp int64, reference ...int64) string {
	return ""
}

// Duration formats a duration in seconds as a human-readable string.
func Duration(seconds float64, opts ...DurationOption) string {
	return ""
}

// ParseDuration parses a human-written duration string into seconds.
func ParseDuration(input string) (int64, error) {
	return 0, nil
}

// HumanDate returns a contextual date string.
// The optional reference parameter is used for comparison to determine relative output.
func HumanDate(timestamp int64, reference ...int64) string {
	return ""
}

// DateRange formats a date range with smart abbreviation.
func DateRange(start, end int64) string {
	return ""
}
