package whenwords

import (
	"errors"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
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

// TestTimeago tests TimeAgo function using YAML test data.
func TestTimeago(t *testing.T) {
	suite := loadTestCases()

	for _, tc := range suite.TimeAgo {
		t.Run(tc.Name, func(t *testing.T) {
			got := TimeAgo(tc.Input.Timestamp, tc.Input.Reference)
			if got != tc.Output {
				t.Errorf("TimeAgo(%d, %d) = %q, want %q",
					tc.Input.Timestamp, tc.Input.Reference, got, tc.Output)
			}
		})
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

// --- YAML Test Infrastructure ---

// TestSuite represents the full test data from tests.yaml.
type TestSuite struct {
	Version       string              `yaml:"version"`
	TimeAgo       []TimeAgoTest       `yaml:"timeago"`
	Duration      []DurationTest      `yaml:"duration"`
	ParseDuration []ParseDurationTest `yaml:"parse_duration"`
	HumanDate     []HumanDateTest     `yaml:"human_date"`
	DateRange     []DateRangeTest     `yaml:"date_range"`
}

// TimeAgoInput represents input for timeago tests.
type TimeAgoInput struct {
	Timestamp int64 `yaml:"timestamp"`
	Reference int64 `yaml:"reference"`
}

// TimeAgoTest represents a single timeago test case.
type TimeAgoTest struct {
	Name   string       `yaml:"name"`
	Input  TimeAgoInput `yaml:"input"`
	Output string       `yaml:"output"`
}

// DurationOptions represents optional configuration for duration tests.
type DurationOptions struct {
	Compact  bool `yaml:"compact"`
	MaxUnits int  `yaml:"max_units"`
}

// DurationInput represents input for duration tests.
type DurationInput struct {
	Seconds int64           `yaml:"seconds"`
	Options DurationOptions `yaml:"options"`
}

// DurationTest represents a single duration test case.
type DurationTest struct {
	Name   string        `yaml:"name"`
	Input  DurationInput `yaml:"input"`
	Output string        `yaml:"output"`
	Error  bool          `yaml:"error"`
}

// ParseDurationTest represents a single parse_duration test case.
type ParseDurationTest struct {
	Name   string `yaml:"name"`
	Input  string `yaml:"input"`
	Output int64  `yaml:"output"`
	Error  bool   `yaml:"error"`
}

// HumanDateInput represents input for human_date tests.
type HumanDateInput struct {
	Timestamp int64 `yaml:"timestamp"`
	Reference int64 `yaml:"reference"`
}

// HumanDateTest represents a single human_date test case.
type HumanDateTest struct {
	Name   string         `yaml:"name"`
	Input  HumanDateInput `yaml:"input"`
	Output string         `yaml:"output"`
}

// DateRangeInput represents input for date_range tests.
type DateRangeInput struct {
	Start int64 `yaml:"start"`
	End   int64 `yaml:"end"`
}

// DateRangeTest represents a single date_range test case.
type DateRangeTest struct {
	Name   string         `yaml:"name"`
	Input  DateRangeInput `yaml:"input"`
	Output string         `yaml:"output"`
}

// loadTestCases loads test data from testdata/tests.yaml.
func loadTestCases() TestSuite {
	data, err := os.ReadFile("testdata/tests.yaml")
	if err != nil {
		panic("failed to read testdata/tests.yaml: " + err.Error())
	}

	var suite TestSuite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		panic("failed to parse testdata/tests.yaml: " + err.Error())
	}

	return suite
}

// TestYAMLInfrastructure verifies YAML test loading infrastructure works.
func TestYAMLInfrastructure(t *testing.T) {
	suite := loadTestCases()

	// Count total test cases across all categories
	totalTests := len(suite.TimeAgo) + len(suite.Duration) + len(suite.ParseDuration) + len(suite.HumanDate) + len(suite.DateRange)

	t.Logf("Loaded %d test cases from tests.yaml (version %s)", totalTests, suite.Version)

	if totalTests == 0 {
		t.Error("Expected at least one test case to be loaded")
	}
}
