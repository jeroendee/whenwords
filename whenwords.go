// Package whenwords provides human-friendly time formatting and parsing.
package whenwords

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
// It returns ErrNegativeDuration for negative input.
func Duration(seconds int64, opts ...DurationOption) (string, error) {
	if seconds < 0 {
		return "", ErrNegativeDuration
	}

	cfg := &durationConfig{maxUnits: 2}
	for _, opt := range opts {
		opt(cfg)
	}

	// Time unit constants
	const (
		secondsPerMinute = 60
		secondsPerHour   = 3600
		secondsPerDay    = 86400
		secondsPerMonth  = 30 * secondsPerDay
		secondsPerYear   = 365 * secondsPerDay
	)

	// Unit names for compact and verbose output
	type unit struct {
		seconds int64
		compact string
		verbose string
	}
	units := []unit{
		{secondsPerYear, "y", "year"},
		{secondsPerMonth, "mo", "month"},
		{secondsPerDay, "d", "day"},
		{secondsPerHour, "h", "hour"},
		{secondsPerMinute, "m", "minute"},
		{1, "s", "second"},
	}

	// Extract non-zero units
	type part struct {
		value   int64
		compact string
		verbose string
	}
	var parts []part
	remaining := seconds

	for _, u := range units {
		if remaining >= u.seconds {
			count := remaining / u.seconds
			remaining %= u.seconds
			parts = append(parts, part{count, u.compact, u.verbose})
		}
	}

	// Handle zero case
	if len(parts) == 0 {
		if cfg.compact {
			return "0s", nil
		}
		return "0 seconds", nil
	}

	// Apply max_units limit
	if cfg.maxUnits > 0 && len(parts) > cfg.maxUnits {
		parts = parts[:cfg.maxUnits]
	}

	// Build output string
	var result string
	for i, p := range parts {
		if i > 0 {
			if cfg.compact {
				result += " "
			} else {
				result += ", "
			}
		}

		if cfg.compact {
			result += itoa(int(p.value)) + p.compact
		} else {
			result += itoa(int(p.value)) + " " + p.verbose
			if p.value != 1 {
				result += "s"
			}
		}
	}

	return result, nil
}

// ParseDuration parses a human-written duration string into seconds.
func ParseDuration(input string) (int64, error) {
	// Handle empty input
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, ErrEmptyInput
	}

	// Check for negative sign
	if strings.HasPrefix(input, "-") {
		return 0, ErrNegativeValue
	}

	// Try colon notation first (h:mm or h:mm:ss)
	if seconds, ok := parseColonNotation(input); ok {
		return seconds, nil
	}

	// Tokenize with regex for unit-value pairs
	seconds, found := parseUnitValuePairs(input)
	if !found {
		return 0, ErrUnparseable
	}

	return seconds, nil
}

// parseColonNotation parses h:mm or h:mm:ss format.
func parseColonNotation(input string) (int64, bool) {
	// Match patterns like 2:30 or 1:30:00 or 0:05:30
	colonPattern := regexp.MustCompile(`^(\d+):(\d{2})(?::(\d{2}))?$`)
	matches := colonPattern.FindStringSubmatch(input)
	if matches == nil {
		return 0, false
	}

	hours, _ := strconv.ParseInt(matches[1], 10, 64)
	minutes, _ := strconv.ParseInt(matches[2], 10, 64)
	var seconds int64
	if matches[3] != "" {
		seconds, _ = strconv.ParseInt(matches[3], 10, 64)
	}

	return hours*3600 + minutes*60 + seconds, true
}

// parseUnitValuePairs extracts value-unit pairs and sums them.
func parseUnitValuePairs(input string) (int64, bool) {
	// Unit multipliers in seconds
	unitMultipliers := map[string]int64{
		"w":       604800, // week
		"week":    604800,
		"weeks":   604800,
		"d":       86400, // day
		"day":     86400,
		"days":    86400,
		"h":       3600, // hour
		"hr":      3600,
		"hrs":     3600,
		"hour":    3600,
		"hours":   3600,
		"m":       60, // minute
		"min":     60,
		"mins":    60,
		"minute":  60,
		"minutes": 60,
		"s":       1, // second
		"sec":     1,
		"secs":    1,
		"second":  1,
		"seconds": 1,
	}

	// Pattern: number (possibly decimal) followed by unit
	// Handles: 2h, 2.5h, 2 hours, 2.5 hours
	pattern := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(w|weeks?|d|days?|h|hrs?|hours?|m|mins?|minutes?|s|secs?|seconds?)`)

	matches := pattern.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return 0, false
	}

	var total int64
	for _, match := range matches {
		valueStr := match[1]
		unit := strings.ToLower(match[2])

		multiplier, ok := unitMultipliers[unit]
		if !ok {
			continue
		}

		// Handle decimal values
		if strings.Contains(valueStr, ".") {
			value, _ := strconv.ParseFloat(valueStr, 64)
			total += int64(value * float64(multiplier))
		} else {
			value, _ := strconv.ParseInt(valueStr, 10, 64)
			total += value * multiplier
		}
	}

	return total, true
}

// HumanDate returns a contextual date string.
// The optional reference parameter is used for comparison to determine relative output.
func HumanDate(timestamp int64, reference ...int64) string {
	ref := timestamp
	if len(reference) > 0 {
		ref = reference[0]
	}

	// Convert to UTC time objects
	tsTime := time.Unix(timestamp, 0).UTC()
	refTime := time.Unix(ref, 0).UTC()

	// Truncate to UTC midnight for calendar day comparison
	tsDay := time.Date(tsTime.Year(), tsTime.Month(), tsTime.Day(), 0, 0, 0, 0, time.UTC)
	refDay := time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate day difference
	dayDiff := int(tsDay.Sub(refDay).Hours() / 24)

	switch {
	case dayDiff == 0:
		return "Today"
	case dayDiff == -1:
		return "Yesterday"
	case dayDiff == 1:
		return "Tomorrow"
	case dayDiff >= -6 && dayDiff <= -2:
		return "Last " + tsTime.Weekday().String()
	case dayDiff >= 2 && dayDiff <= 6:
		return "This " + tsTime.Weekday().String()
	default:
		// Format as date
		if tsTime.Year() == refTime.Year() {
			return tsTime.Format("January 2")
		}
		return tsTime.Format("January 2, 2006")
	}
}

// DateRange formats a date range with smart abbreviation.
func DateRange(start, end int64) string {
	return ""
}
