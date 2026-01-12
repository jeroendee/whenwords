// Package whenwords provides human-friendly time formatting and parsing.
package whenwords

import (
	"errors"
	"strconv"
)

// Sentinel errors for the whenwords package.
var (
	// ErrNegativeDuration is returned when a negative duration is provided.
	ErrNegativeDuration = errors.New("negative duration not allowed")

	// ErrEmptyInput is returned when an empty string is provided for parsing.
	ErrEmptyInput = errors.New("empty input")

	// ErrUnparseable is returned when input cannot be parsed as a duration.
	ErrUnparseable = errors.New("unable to parse duration")

	// ErrNegativeValue is returned when a parsed value would be negative.
	ErrNegativeValue = errors.New("negative value not allowed")
)

// durationOptions holds configuration for duration formatting.
type durationOptions struct {
	compact  bool
	maxUnits int
}

// DurationOption is a functional option for configuring duration formatting.
type DurationOption func(*durationOptions)

// WithCompact returns an option that enables compact output format.
func WithCompact() DurationOption {
	return func(o *durationOptions) {
		o.compact = true
	}
}

// WithMaxUnits returns an option that limits the number of time units displayed.
func WithMaxUnits(n int) DurationOption {
	return func(o *durationOptions) {
		o.maxUnits = n
	}
}

// TimeAgo returns a human-readable relative time string.
// Positive diff (timestamp < reference) means past: "{n} {units} ago"
// Negative diff (timestamp > reference) means future: "in {n} {units}"
func TimeAgo(timestamp, reference int64) string {
	diff := reference - timestamp
	future := diff < 0
	if future {
		diff = -diff
	}

	seconds := float64(diff)
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	months := days / 30    // approximation
	years := days / 365.25 // approximation

	var value int
	var unit string

	switch {
	case seconds < 45:
		return "just now"
	case seconds < 90:
		value = 1
		unit = "minute"
	case minutes < 45:
		value = roundHalfUp(minutes)
		unit = "minute"
	case minutes < 90:
		value = 1
		unit = "hour"
	case hours < 22:
		value = roundHalfUp(hours)
		unit = "hour"
	case hours < 36:
		value = 1
		unit = "day"
	case days < 26:
		value = roundHalfUp(days)
		unit = "day"
	case days < 46:
		value = 1
		unit = "month"
	case days < 320:
		value = int(months)
		if value < 2 {
			value = 2 // minimum is 2 months in this range (46+ days)
		}
		unit = "month"
	case days < 548:
		value = 1
		unit = "year"
	default:
		value = roundHalfUp(years)
		unit = "year"
	}

	if value != 1 {
		unit += "s"
	}

	if future {
		return "in " + formatInt(value) + " " + unit
	}
	return formatInt(value) + " " + unit + " ago"
}

// roundHalfUp rounds a float to the nearest integer using half-up rounding.
func roundHalfUp(x float64) int {
	return int(x + 0.5)
}

// formatInt converts an integer to its string representation.
func formatInt(n int) string {
	return strconv.Itoa(n)
}
