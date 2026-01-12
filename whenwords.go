// Package whenwords provides human-friendly time formatting and parsing.
package whenwords

import "errors"

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
