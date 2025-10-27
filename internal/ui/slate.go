package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/jawahars16/jebi/internal/core"
)

type slate struct {
	accentColor    lipgloss.Color
	activeSpinner  *spinner.Model
	spinnerMessage string
}

func NewSlate(accentColor lipgloss.Color) *slate {
	return &slate{accentColor: accentColor}
}

func (s *slate) PromptWithDefault(message, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [%s]: ", message, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

func (s *slate) ShowHeader(title string) {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(s.accentColor).
		Padding(0, 2).
		Margin(1, 0, 1, 0).
		Align(lipgloss.Left)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("213"))

	box := borderStyle.Render(title)
	fmt.Println(titleStyle.Render(box))
}

func (s *slate) ShowList(title string, items []string, highlight string) {
	l := list.New(items)
	l = l.Enumerator(func(ls list.Items, i int) string {
		if ls.At(i).Value() == highlight {
			return lipgloss.NewStyle().Foreground(s.accentColor).Render("➤")
		}
		return ""
	})
	l = l.ItemStyleFunc(func(ls list.Items, i int) lipgloss.Style {
		if ls.At(i).Value() == highlight {
			return lipgloss.NewStyle().Foreground(s.accentColor).Bold(true)
		}
		return lipgloss.NewStyle().Bold(false)
	})
	fmt.Println(l)
}

func (s *slate) WriteStatus(changes []core.Change) {
	var (
		addStyle = lipgloss.NewStyle().Padding(0, 0, 0, 1).Foreground(lipgloss.Color("34"))  // Green
		modStyle = lipgloss.NewStyle().Padding(0, 0, 0, 1).Foreground(lipgloss.Color("214")) // Orange
		delStyle = lipgloss.NewStyle().Padding(0, 0, 0, 1).Foreground(lipgloss.Color("131")) // Red
		padding  = 2
	)

	// Build action label map
	labelMap := map[string]string{
		string(core.ChangeTypeAdd):    "added:",
		string(core.ChangeTypeModify): "modified:",
		string(core.ChangeTypeRemove): "removed:",
	}

	// Find longest label for alignment
	maxLen := 0
	for _, c := range changes {
		if l := len(labelMap[string(c.Type)]); l > maxLen {
			maxLen = l
		}
	}

	// Print all aligned lines
	for _, c := range changes {
		label := labelMap[string(c.Type)]
		var style lipgloss.Style
		var symbol string

		switch c.Type {
		case core.ChangeTypeAdd:
			style = addStyle
			symbol = "+"
		case core.ChangeTypeModify:
			style = modStyle
			symbol = "~"
		case core.ChangeTypeRemove:
			style = delStyle
			symbol = "-"
		default:
			style = lipgloss.NewStyle()
			symbol = " "
		}

		// Align columns with dynamic padding
		formatted := fmt.Sprintf("%-*s %s", maxLen+padding, label, c.Key)
		fmt.Println(style.Render(fmt.Sprintf("%s %s", symbol, formatted)))
	}
}

func (s *slate) RenderMarkdown(md string) {
	r, err := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		// wrap output at specific width (default is 80)
		glamour.WithWordWrap(0),
	)
	if err != nil {
		fmt.Printf("failed to create markdown renderer: %v", err)
	}

	out, err := r.Render(md)
	if err != nil {
		fmt.Printf("failed to render markdown: %v", err)
	}

	fmt.Println(out)
}

func (s *slate) ShowWarning(msg string) {
	border := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("178")). // warm amber border
		Padding(0, 1).
		Margin(1, 0, 1, 0)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("178"))

	box := border.Render(msg)
	fmt.Println(title.Render(box))
}

func (s *slate) ShowError(msg string) {
	border := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("131")). // bright red border
		Padding(0, 1).
		Margin(1, 0, 1, 0)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")) // red title color

	box := border.Render(msg)
	fmt.Println(title.Render(box))
}

// StyleOptions defines styling options for WriteStyledText
type StyleOptions struct {
	Color       lipgloss.Color // Text color (e.g., "196" for red, "34" for green)
	Background  lipgloss.Color // Background color (optional)
	Bold        bool           // Make text bold
	Italic      bool           // Make text italic
	Underline   bool           // Make text underlined
	Padding     []int          // Padding [top, right, bottom, left] or [vertical, horizontal] or [all]
	Margin      []int          // Margin [top, right, bottom, left] or [vertical, horizontal] or [all]
	Border      bool           // Add border around text
	BorderColor lipgloss.Color // Border color (uses text color if not specified)
}

// WriteStyledText writes text to terminal with custom color and styling
func (s *slate) WriteStyledText(text string, options StyleOptions) {
	style := lipgloss.NewStyle()

	// Apply color
	if options.Color != "" {
		style = style.Foreground(options.Color)
	}

	// Apply background color
	if options.Background != "" {
		style = style.Background(options.Background)
	}

	// Apply text styling
	if options.Bold {
		style = style.Bold(true)
	}
	if options.Italic {
		style = style.Italic(true)
	}
	if options.Underline {
		style = style.Underline(true)
	}

	// Apply padding
	if len(options.Padding) > 0 {
		switch len(options.Padding) {
		case 1:
			style = style.Padding(options.Padding[0])
		case 2:
			style = style.Padding(options.Padding[0], options.Padding[1])
		case 4:
			style = style.Padding(options.Padding[0], options.Padding[1], options.Padding[2], options.Padding[3])
		}
	}

	// Apply margin
	if len(options.Margin) > 0 {
		switch len(options.Margin) {
		case 1:
			style = style.Margin(options.Margin[0])
		case 2:
			style = style.Margin(options.Margin[0], options.Margin[1])
		case 4:
			style = style.Margin(options.Margin[0], options.Margin[1], options.Margin[2], options.Margin[3])
		}
	}

	// Apply border
	if options.Border {
		borderColor := options.BorderColor
		if borderColor == "" {
			borderColor = options.Color // Use text color if border color not specified
		}
		style = style.Border(lipgloss.RoundedBorder()).BorderForeground(borderColor)
	}

	fmt.Println(style.Render(text))
}

// WriteColoredText is a convenience function for simple colored text output
func (s *slate) WriteColoredText(text string, color lipgloss.Color) {
	s.WriteStyledText(text, StyleOptions{Color: color})
}

// WriteIndentedText writes text with default left padding and optional styling
// Provides consistent indentation across the application
func (s *slate) WriteIndentedText(text string, options StyleOptions) {
	// Set default left padding if no padding is specified
	if len(options.Padding) == 0 {
		options.Padding = []int{0, 0, 0, 2} // [top, right, bottom, left] - 2 spaces left padding
	} else {
		// If padding is specified, ensure left padding is at least 2
		switch len(options.Padding) {
		case 1:
			// Convert single value to [top, right, bottom, left] with minimum left padding
			if options.Padding[0] < 2 {
				options.Padding = []int{options.Padding[0], options.Padding[0], options.Padding[0], 2}
			} else {
				options.Padding = []int{options.Padding[0], options.Padding[0], options.Padding[0], options.Padding[0]}
			}
		case 2:
			// Convert [vertical, horizontal] to [top, right, bottom, left] with minimum left padding
			leftPadding := options.Padding[1]
			if leftPadding < 2 {
				leftPadding = 2
			}
			options.Padding = []int{options.Padding[0], options.Padding[1], options.Padding[0], leftPadding}
		case 4:
			// Ensure left padding (index 3) is at least 2
			if options.Padding[3] < 2 {
				options.Padding[3] = 2
			}
		}
	}

	s.WriteStyledText(text, options)
}

// ShowSuccess displays a success message with consistent styling
func (s *slate) ShowSuccess(message string) {
	s.WriteStyledText(message, StyleOptions{
		Color:   "34", // Green
		Bold:    true,
		Padding: []int{0, 1, 0, 2}, // Left padding with some right padding
	})
}

// ShowEnvironmentContext displays current environment information
func (s *slate) ShowEnvironmentContext(env string) {
	s.WriteStyledText(fmt.Sprintf("Environment: %s", env), StyleOptions{
		Color:   "82", // Light green
		Italic:  true,
		Padding: []int{0, 0, 0, 2}, // Left padding only
	})
}

// ShowSecretOperation displays information about a secret operation (add/set/remove)
func (s *slate) ShowSecretOperation(operation core.ChangeType, key, env string, isPlaintext bool) {
	// Display environment context
	s.ShowEnvironmentContext(env)
	var (
		color  lipgloss.Color
		status string
	)
	switch operation {
	case core.ChangeTypeAdd:
		color = "34" // Green
		status = "added"
	case core.ChangeTypeModify:
		color = "214" // Orange
		status = "modified"
	case core.ChangeTypeRemove:
		color = "131" // Red
		status = "removed"
	default:
		color = "15" // White
	}
	// Display operation summary
	operationText := fmt.Sprintf("Secret '%s' %s successfully", key, status)
	if isPlaintext {
		operationText += " (plaintext)"
	}
	s.WriteIndentedText(operationText, StyleOptions{
		Color: color,
		Bold:  true,
	})
}

func (s *slate) RenderInitHeader() {
	// Define color palette
	subtle := lipgloss.Color("#94A3B8")  // Slate-400
	success := lipgloss.Color("#22C55E") // Green-500

	// Title styles
	titleStyle := lipgloss.NewStyle().
		Foreground(s.accentColor).
		Bold(true)

	taglineStyle := lipgloss.NewStyle().
		Foreground(subtle).
		Faint(true)

	// Body styles
	sectionStyle := lipgloss.NewStyle().
		MarginLeft(0)

	bulletStyle := lipgloss.NewStyle().
		Foreground(s.accentColor)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5E7EB")) // Neutral light gray

	footerStyle := lipgloss.NewStyle().
		Foreground(success).
		MarginTop(1)

	// Construct header text
	header := fmt.Sprintf(
		"%s  %s\n%s\n\n%s\n%s\n%s\n%s",
		titleStyle.Render("jebi"),
		taglineStyle.Render("\ngit for secrets"),
		lipgloss.NewStyle().Faint(true).Render("───────────────────────────────────────────────"),
		sectionStyle.Render("Initializing a new jebi project:"),
		sectionStyle.Render(bulletStyle.Render("• ")+textStyle.Render("Set up your first environment (dev/prod)")),
		sectionStyle.Render(bulletStyle.Render("• ")+textStyle.Render("Generate and encrypt a symmetric key")),
		sectionStyle.Render(bulletStyle.Render("• ")+textStyle.Render("Prepare project for secret versioning")),
	)

	footer := footerStyle.Render("Let's get your secret workflow started!")

	// Combine everything neatly in a styled box
	boxStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.ThickBorder()).
		BorderForeground(s.accentColor)

	fmt.Println(boxStyle.Render(header + "\n\n" + footer))
}

// StartSpinner starts a spinner with a message for long-running operations
func (s *slate) StartSpinner(message string) {
	s.activeSpinner = &spinner.Model{}
	*s.activeSpinner = spinner.New()
	s.activeSpinner.Spinner = spinner.Dot
	s.activeSpinner.Style = lipgloss.NewStyle().Foreground(s.accentColor)
	s.spinnerMessage = message

	// Print initial spinner state
	fmt.Printf("\r%s %s", s.activeSpinner.View(), message)
}

// UpdateSpinner updates the spinner animation and optionally changes the message
func (s *slate) UpdateSpinner(newMessage ...string) {
	if s.activeSpinner == nil {
		return
	}

	// Update message if provided
	if len(newMessage) > 0 && newMessage[0] != "" {
		s.spinnerMessage = newMessage[0]
	}

	time.Sleep(100 * time.Millisecond)
	// Tick the spinner animation
	cmd := s.activeSpinner.Tick()
	if cmd != nil {
		// Update the spinner model
		newModel, _ := s.activeSpinner.Update(cmd)
		*s.activeSpinner = newModel
	}

	// Clear the line and print updated spinner
	fmt.Printf("\r\033[K%s %s", s.activeSpinner.View(), s.spinnerMessage)
}

// StopSpinner stops the spinner and shows a completion message
func (s *slate) StopSpinner() {
	if s.activeSpinner == nil {
		return
	}

	// Clear the spinner line
	fmt.Printf("\r\033[K")
	// Reset spinner state
	s.activeSpinner = nil
	s.spinnerMessage = ""
}

// ShowSpinnerOperation runs a spinner for the duration of an operation
// This is a convenience method for simple operations
func (s *slate) ShowSpinnerOperation(message string, operation func() error) error {
	s.StartSpinner(message)

	// Simple animation loop in a goroutine
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				s.UpdateSpinner()
			}
		}
	}()

	// Run the operation
	err := operation()

	// Stop the animation
	done <- true

	// Stop spinner with appropriate message
	if err != nil {
		s.StopSpinnerWithError(err.Error())
	} else {
		s.StopSpinner()
	}

	return err
}

// StopSpinnerWithError stops the spinner and shows an error message
func (s *slate) StopSpinnerWithError(errorMessage string) {
	if s.activeSpinner == nil {
		return
	}

	// Clear the spinner line
	fmt.Printf("\r\033[K")

	// Show error message
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	fmt.Printf("%s %s\n", errorStyle.Render("❌"), errorMessage)

	// Reset spinner state
	s.activeSpinner = nil
	s.spinnerMessage = ""
}

/*
Usage examples:

// Simple colored text
slate.WriteColoredText("Success!", "34") // Green text
slate.WriteColoredText("Error!", "196")  // Red text

// Indented text with default left padding (2 spaces)
slate.WriteIndentedText("This is indented", StyleOptions{
    Color: "34",  // Green text with automatic left padding
    Bold:  true,
})

// Indented text with custom styling
slate.WriteIndentedText("Status update", StyleOptions{
    Color:   "214", // Orange
    Italic:  true,
    // Left padding will be at least 2, even if you specify less
})

// Styled text with multiple options
slate.WriteStyledText("Important Notice", StyleOptions{
    Color:     "214",    // Orange text
    Bold:      true,     // Make it bold
    Border:    true,     // Add border
    Padding:   []int{1, 2}, // Vertical padding: 1, Horizontal: 2
    Margin:    []int{1},    // All margins: 1
})

// Text with background
slate.WriteStyledText("Highlighted", StyleOptions{
    Color:      "15",  // White text
    Background: "34",  // Green background
    Bold:       true,
    Padding:    []int{0, 1}, // Add some horizontal padding
})

// Underlined and italic text
slate.WriteStyledText("Emphasized text", StyleOptions{
    Color:     "213", // Pink
    Italic:    true,
    Underline: true,
})

// Success and context helpers
slate.ShowSuccess("Operation completed successfully!")
slate.ShowEnvironmentContext("production")
slate.ShowSecretOperation("added", "DATABASE_URL", "production", false)

// Progress indicators for remote operations
slate.ShowProgress("Connecting to server...")
slate.ShowProgressWithSpinner("Uploading secrets...", 0) // Call with incrementing index for animation
// Spinner examples
// Simple spinner usage
slate.StartSpinner("Processing...")
// ... do some work ...
slate.UpdateSpinner("Almost done...")
// ... finish work ...
slate.StopSpinner("Process completed!")

// Spinner with error
slate.StartSpinner("Connecting to server...")
slate.StopSpinnerWithError("Failed to connect to server")

// Convenient spinner wrapper
err := slate.ShowSpinnerOperation("Uploading secrets...", func() error {
    // Your operation here
    time.Sleep(2 * time.Second) // Simulate work
    return nil // or return an error
})

Common colors:
- "34"  - Green
- "196" - Red
- "214" - Orange
- "82"  - Light green
- "213" - Pink/Purple
- "178" - Yellow/Amber
- "15"  - White
- "0"   - Black
*/
