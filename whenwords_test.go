package whenwords

import (
	"errors"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestSpec represents the full YAML test specification.
type TestSpec struct {
	Version       string              `yaml:"version"`
	Timeago       []TimeagoTestCase   `yaml:"timeago"`
	Duration      []DurationTestCase  `yaml:"duration"`
	ParseDuration []ParseTestCase     `yaml:"parse_duration"`
	HumanDate     []HumanDateTestCase `yaml:"human_date"`
	DateRange     []DateRangeTestCase `yaml:"date_range"`
}

// TimeagoTestCase represents a single timeago test.
type TimeagoTestCase struct {
	Name   string `yaml:"name"`
	Input  struct {
		Timestamp int64 `yaml:"timestamp"`
		Reference int64 `yaml:"reference"`
	} `yaml:"input"`
	Output string `yaml:"output"`
}

// DurationTestCase represents a single duration formatting test.
type DurationTestCase struct {
	Name  string `yaml:"name"`
	Input struct {
		Seconds int64 `yaml:"seconds"`
		Options struct {
			Compact  bool `yaml:"compact"`
			MaxUnits int  `yaml:"max_units"`
		} `yaml:"options"`
	} `yaml:"input"`
	Output string `yaml:"output"`
	Error  bool   `yaml:"error"`
}

// ParseTestCase represents a single parse_duration test.
type ParseTestCase struct {
	Name   string `yaml:"name"`
	Input  string `yaml:"input"`
	Output int64  `yaml:"output"`
	Error  bool   `yaml:"error"`
}

// HumanDateTestCase represents a single human_date test.
type HumanDateTestCase struct {
	Name  string `yaml:"name"`
	Input struct {
		Timestamp int64 `yaml:"timestamp"`
		Reference int64 `yaml:"reference"`
	} `yaml:"input"`
	Output string `yaml:"output"`
}

// DateRangeTestCase represents a single date_range test.
type DateRangeTestCase struct {
	Name  string `yaml:"name"`
	Input struct {
		Start int64 `yaml:"start"`
		End   int64 `yaml:"end"`
	} `yaml:"input"`
	Output string `yaml:"output"`
}

func loadTestSpec(t *testing.T) TestSpec {
	t.Helper()
	data, err := os.ReadFile("testdata/tests.yaml")
	if err != nil {
		t.Fatalf("failed to read testdata/tests.yaml: %v", err)
	}
	var spec TestSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		t.Fatalf("failed to parse testdata/tests.yaml: %v", err)
	}
	return spec
}

func TestLoadTestSpec(t *testing.T) {
	spec := loadTestSpec(t)

	if spec.Version == "" {
		t.Error("expected version to be set")
	}

	// Verify expected test case counts (actual counts from tests.yaml)
	if got := len(spec.Timeago); got != 36 {
		t.Errorf("timeago: expected 36 test cases, got %d", got)
	}
	if got := len(spec.Duration); got != 26 {
		t.Errorf("duration: expected 26 test cases, got %d", got)
	}
	if got := len(spec.ParseDuration); got != 32 {
		t.Errorf("parse_duration: expected 32 test cases, got %d", got)
	}
	if got := len(spec.HumanDate); got != 20 {
		t.Errorf("human_date: expected 20 test cases, got %d", got)
	}
	if got := len(spec.DateRange); got != 9 {
		t.Errorf("date_range: expected 9 test cases, got %d", got)
	}
}

func TestSentinelErrorsDefined(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNegativeDuration", ErrNegativeDuration},
		{"ErrEmptyInput", ErrEmptyInput},
		{"ErrUnparseable", ErrUnparseable},
		{"ErrNegativeValue", ErrNegativeValue},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Errorf("%s should not be nil", tc.name)
			}
			// Verify each error is unique
			for _, other := range tests {
				if other.name != tc.name && errors.Is(tc.err, other.err) {
					t.Errorf("%s should not be equal to %s", tc.name, other.name)
				}
			}
		})
	}
}

func TestDurationOptionType(t *testing.T) {
	// Verify the functional options pattern works
	opts := &durationOptions{}

	// Apply options
	WithCompact()(opts)
	if !opts.compact {
		t.Error("WithCompact should set compact to true")
	}

	opts2 := &durationOptions{}
	WithMaxUnits(2)(opts2)
	if opts2.maxUnits != 2 {
		t.Errorf("WithMaxUnits(2) should set maxUnits to 2, got %d", opts2.maxUnits)
	}
}

func TestTimeago(t *testing.T) {
	spec := loadTestSpec(t)

	for _, tc := range spec.Timeago {
		t.Run(tc.Name, func(t *testing.T) {
			got := TimeAgo(tc.Input.Timestamp, tc.Input.Reference)
			if got != tc.Output {
				t.Errorf("TimeAgo(%d, %d) = %q, want %q",
					tc.Input.Timestamp, tc.Input.Reference, got, tc.Output)
			}
		})
	}
}
