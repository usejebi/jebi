package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/jawahars16/jebi/internal/core"
)

type slate struct {
	accentColor lipgloss.Color
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
		core.ActionAdd:    "added:",
		core.ActionUpdate: "modified:",
		core.ActionRemove: "removed:",
	}

	// Find longest label for alignment
	maxLen := 0
	for _, c := range changes {
		if l := len(labelMap[c.Action]); l > maxLen {
			maxLen = l
		}
	}

	// Print all aligned lines
	for _, c := range changes {
		label := labelMap[c.Action]
		var style lipgloss.Style
		var symbol string

		switch c.Action {
		case core.ActionAdd:
			style = addStyle
			symbol = "+"
		case core.ActionUpdate:
			style = modStyle
			symbol = "~"
		case core.ActionRemove:
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
		fmt.Printf("failed to create markdown renderer: %w", err)
	}

	out, err := r.Render(md)
	if err != nil {
		fmt.Printf("failed to render markdown: %w", err)
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
		BorderForeground(lipgloss.Color("196")). // bright red border
		Padding(0, 1).
		Margin(1, 0, 1, 0)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")) // red title color

	box := border.Render(msg)
	fmt.Println(title.Render(box))
}
