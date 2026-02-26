package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m"
)

var (
	// NoColor disables color output
	NoColor = false
	// Output is the writer for UI output
	Output io.Writer = os.Stdout
)

// Colorize wraps text with color codes
func Colorize(color, text string) string {
	if NoColor {
		return text
	}
	return color + text + Reset
}

// Success prints a success message
func Success(message string) {
	fmt.Fprintf(Output, "%s %s\n", Colorize(Green, "✓"), message)
}

// Error prints an error message
func Error(message string) {
	fmt.Fprintf(Output, "%s %s\n", Colorize(Red, "✗"), message)
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Fprintf(Output, "%s %s\n", Colorize(Yellow, "⚠"), message)
}

// Info prints an info message
func Info(message string) {
	fmt.Fprintf(Output, "%s %s\n", Colorize(Blue, "ℹ"), message)
}

// Step prints a step message
func Step(message string) {
	fmt.Fprintf(Output, "%s %s\n", Colorize(Cyan, "→"), message)
}

// Header prints a header message
func Header(message string) {
	fmt.Fprintf(Output, "\n%s\n", Colorize(Bold+Cyan, message))
}

// ProgressBar represents a simple progress bar
type ProgressBar struct {
	total   int
	current int
	width   int
	prefix  string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, prefix string) *ProgressBar {
	return &ProgressBar{
		total:  total,
		width:  40,
		prefix: prefix,
	}
}

// Increment increments the progress bar
func (pb *ProgressBar) Increment() {
	pb.current++
	pb.Render()
}

// Set sets the current progress
func (pb *ProgressBar) Set(current int) {
	pb.current = current
	pb.Render()
}

// Render renders the progress bar
func (pb *ProgressBar) Render() {
	if NoColor {
		fmt.Fprintf(Output, "\r%s: %d/%d", pb.prefix, pb.current, pb.total)
		return
	}

	percent := float64(pb.current) / float64(pb.total)
	filled := int(percent * float64(pb.width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)
	fmt.Fprintf(Output, "\r%s [%s] %d/%d (%.0f%%)",
		pb.prefix,
		Colorize(Green, bar),
		pb.current,
		pb.total,
		percent*100)

	if pb.current >= pb.total {
		fmt.Fprintln(Output)
	}
}

// Spinner represents a loading spinner
type Spinner struct {
	message string
	frames  []string
	stop    chan bool
	done    chan bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:    make(chan bool),
		done:    make(chan bool),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	if NoColor {
		fmt.Fprintf(Output, "%s...\n", s.message)
		return
	}

	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Fprintf(Output, "\r%s\r", strings.Repeat(" ", len(s.message)+10))
				s.done <- true
				return
			default:
				frame := s.frames[i%len(s.frames)]
				fmt.Fprintf(Output, "\r%s %s", Colorize(Cyan, frame), s.message)
				time.Sleep(80 * time.Millisecond)
				i++
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	if NoColor {
		return
	}
	s.stop <- true
	<-s.done
}

// Table represents a simple table
type Table struct {
	headers []string
	rows    [][]string
}

// NewTable creates a new table
func NewTable(headers []string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

// Render renders the table
func (t *Table) Render() {
	if len(t.headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(t.headers))
	for i, header := range t.headers {
		widths[i] = len(header)
	}

	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	fmt.Fprint(Output, Colorize(Bold, ""))
	for i, header := range t.headers {
		fmt.Fprintf(Output, "%-*s", widths[i]+2, header)
	}
	fmt.Fprintln(Output, Colorize(Reset, ""))

	// Print separator
	for _, width := range widths {
		fmt.Fprint(Output, strings.Repeat("─", width+2))
	}
	fmt.Fprintln(Output)

	// Print rows
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Fprintf(Output, "%-*s", widths[i]+2, cell)
			}
		}
		fmt.Fprintln(Output)
	}
}

// Prompt prompts the user for input
func Prompt(message string) string {
	fmt.Fprintf(Output, "%s ", Colorize(Cyan, message+":"))
	var input string
	fmt.Scanln(&input)
	return input
}

// Confirm prompts the user for confirmation
func Confirm(message string) bool {
	fmt.Fprintf(Output, "%s %s ", Colorize(Yellow, "?"), message)
	fmt.Fprint(Output, Colorize(Bold, "(y/N):"))
	var input string
	fmt.Scanln(&input)
	return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}
