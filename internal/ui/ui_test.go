package ui

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		text     string
		noColor  bool
		expected string
	}{
		{
			name:     "with color",
			color:    Red,
			text:     "error",
			noColor:  false,
			expected: Red + "error" + Reset,
		},
		{
			name:     "no color mode",
			color:    Red,
			text:     "error",
			noColor:  true,
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldNoColor := NoColor
			NoColor = tt.noColor
			defer func() { NoColor = oldNoColor }()

			result := Colorize(tt.color, tt.text)
			if result != tt.expected {
				t.Errorf("Colorize() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	buf := &bytes.Buffer{}
	oldOutput := Output
	Output = buf
	defer func() { Output = oldOutput }()

	Success("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Success() output does not contain message")
	}
}

func TestError(t *testing.T) {
	buf := &bytes.Buffer{}
	oldOutput := Output
	Output = buf
	defer func() { Output = oldOutput }()

	Error("test error")

	output := buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("Error() output does not contain message")
	}
}

func TestProgressBar(t *testing.T) {
	buf := &bytes.Buffer{}
	oldOutput := Output
	Output = buf
	defer func() { Output = oldOutput }()

	pb := NewProgressBar(10, "Testing")

	for i := 0; i < 10; i++ {
		pb.Increment()
	}

	output := buf.String()
	if !strings.Contains(output, "Testing") {
		t.Errorf("ProgressBar output does not contain prefix")
	}
}

func TestSpinner(t *testing.T) {
	buf := &bytes.Buffer{}
	oldOutput := Output
	Output = buf
	defer func() { Output = oldOutput }()

	spinner := NewSpinner("Loading")
	spinner.Start()
	time.Sleep(200 * time.Millisecond)
	spinner.Stop()

	// Just verify it doesn't crash
}

func TestTable(t *testing.T) {
	buf := &bytes.Buffer{}
	oldOutput := Output
	Output = buf
	defer func() { Output = oldOutput }()

	table := NewTable([]string{"Name", "Age", "City"})
	table.AddRow([]string{"Alice", "30", "NYC"})
	table.AddRow([]string{"Bob", "25", "LA"})
	table.Render()

	output := buf.String()
	if !strings.Contains(output, "Name") || !strings.Contains(output, "Alice") {
		t.Errorf("Table output is incorrect")
	}
}
