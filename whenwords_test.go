package whenwords

import (
	"errors"
	"testing"
)

// TestErrorVariables verifies exported error variables exist and are non-nil.
func TestErrorVariables(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNegativeDuration", ErrNegativeDuration},
		{"ErrEmptyInput", ErrEmptyInput},
		{"ErrUnparseable", ErrUnparseable},
		{"ErrNegativeValue", ErrNegativeValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
		})
	}
}

// TestWithCompact verifies the WithCompact option exists and returns a DurationOption.
func TestWithCompact(t *testing.T) {
	opt := WithCompact()
	if opt == nil {
		t.Error("WithCompact() should return a non-nil DurationOption")
	}
}

// TestWithMaxUnits verifies the WithMaxUnits option exists and returns a DurationOption.
func TestWithMaxUnits(t *testing.T) {
	opt := WithMaxUnits(2)
	if opt == nil {
		t.Error("WithMaxUnits(2) should return a non-nil DurationOption")
	}
}

// TestTimeAgoStub verifies TimeAgo function exists and returns empty string (stub).
func TestTimeAgoStub(t *testing.T) {
	result := TimeAgo(0)
	if result != "" {
		t.Errorf("TimeAgo stub should return empty string, got %q", result)
	}
}

// TestTimeAgoWithReferenceStub verifies TimeAgo accepts optional reference.
func TestTimeAgoWithReferenceStub(t *testing.T) {
	result := TimeAgo(0, 0)
	if result != "" {
		t.Errorf("TimeAgo stub with reference should return empty string, got %q", result)
	}
}

// TestDurationStub verifies Duration function exists and returns empty string (stub).
func TestDurationStub(t *testing.T) {
	result := Duration(0)
	if result != "" {
		t.Errorf("Duration stub should return empty string, got %q", result)
	}
}

// TestDurationWithOptionsStub verifies Duration accepts options.
func TestDurationWithOptionsStub(t *testing.T) {
	result := Duration(0, WithCompact(), WithMaxUnits(1))
	if result != "" {
		t.Errorf("Duration stub with options should return empty string, got %q", result)
	}
}

// TestParseDurationStub verifies ParseDuration exists and returns (0, nil) stub.
func TestParseDurationStub(t *testing.T) {
	result, err := ParseDuration("1h")
	if result != 0 {
		t.Errorf("ParseDuration stub should return 0, got %d", result)
	}
	if err != nil {
		t.Errorf("ParseDuration stub should return nil error, got %v", err)
	}
}

// TestHumanDateStub verifies HumanDate exists and returns empty string (stub).
func TestHumanDateStub(t *testing.T) {
	result := HumanDate(0)
	if result != "" {
		t.Errorf("HumanDate stub should return empty string, got %q", result)
	}
}

// TestHumanDateWithReferenceStub verifies HumanDate accepts optional reference.
func TestHumanDateWithReferenceStub(t *testing.T) {
	result := HumanDate(0, 0)
	if result != "" {
		t.Errorf("HumanDate stub with reference should return empty string, got %q", result)
	}
}

// TestDateRangeStub verifies DateRange exists and returns empty string (stub).
func TestDateRangeStub(t *testing.T) {
	result := DateRange(0, 0)
	if result != "" {
		t.Errorf("DateRange stub should return empty string, got %q", result)
	}
}

// TestErrorsAreDistinct verifies each error is unique (for errors.Is checks).
func TestErrorsAreDistinct(t *testing.T) {
	errs := []error{ErrNegativeDuration, ErrEmptyInput, ErrUnparseable, ErrNegativeValue}
	for i := 0; i < len(errs); i++ {
		for j := i + 1; j < len(errs); j++ {
			if errors.Is(errs[i], errs[j]) {
				t.Errorf("Error %d and %d should be distinct", i, j)
			}
		}
	}
}
