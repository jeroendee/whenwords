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
	ref := timestamp
	if len(reference) > 0 {
		ref = reference[0]
	}

	diff := ref - timestamp
	future := diff < 0
	if future {
		diff = -diff
	}

	// Thresholds in seconds
	const (
		secondsPerMinute = 60
		secondsPerHour   = 3600
		secondsPerDay    = 86400
		daysPerMonth     = 30
		daysPerYear      = 365
	)

	var n int
	var unit string

	switch {
	case diff < 45:
		return "just now"

	case diff < 90:
		n = 1
		unit = "minute"

	case diff < 45*secondsPerMinute:
		n = roundHalfUp(float64(diff) / float64(secondsPerMinute))
		unit = "minute"

	case diff < 90*secondsPerMinute:
		n = 1
		unit = "hour"

	case diff < 22*secondsPerHour:
		n = roundHalfUp(float64(diff) / float64(secondsPerHour))
		unit = "hour"

	case diff < 36*secondsPerHour:
		n = 1
		unit = "day"

	case diff < 26*secondsPerDay:
		n = roundHalfUp(float64(diff) / float64(secondsPerDay))
		unit = "day"

	case diff < 46*secondsPerDay:
		n = 1
		unit = "month"

	case diff < 320*secondsPerDay:
		// Use ~30.44 days per month (365/12) for calculation
		n = roundHalfUp(float64(diff) / (365.0 / 12.0 * float64(secondsPerDay)))
		unit = "month"

	case diff < 548*secondsPerDay:
		n = 1
		unit = "year"

	default:
		n = roundHalfUp(float64(diff) / float64(daysPerYear*secondsPerDay))
		unit = "year"
	}

	// Pluralize
	if n != 1 {
		unit += "s"
	}

	if future {
		return "in " + itoa(n) + " " + unit
	}
	return itoa(n) + " " + unit + " ago"
}

// roundHalfUp rounds to nearest integer with half-up rounding (2.5 -> 3).
func roundHalfUp(f float64) int {
	return int(f + 0.5)
}

// itoa converts int to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
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
