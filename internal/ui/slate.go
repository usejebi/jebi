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

func (s *slate) WriteStatus(key, acton string) {
	var (
		style  lipgloss.Style
		symbol string
	)

	switch acton {
	case core.ActionAdd:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("34")) // Green
		symbol = "+"
	case core.ActionUpdate:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange
		symbol = "~"
	case core.ActionRemove:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("160")) // Red
		symbol = "-"
	default:
		style = lipgloss.NewStyle()
	}
	fmt.Println(style.Render(fmt.Sprintf("%s %s", symbol, key)))
}

func (s *slate) RenderMarkdown(md string) (string, error) {
	// Use "dark" for syntax highlighting. (AutoStyle causes padding in this version)
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"), // or "light" depending on preference
		glamour.WithWordWrap(0),       // disable hard wrap
	)
	if err != nil {
		return "", fmt.Errorf("failed to create markdown renderer: %w", err)
	}

	out, err := r.Render(md)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	// Remove leading/trailing newlines and collapse multiple blank lines
	out = strings.TrimSpace(out)
	out = trimExtraBlankLines(out)
	return out, nil
}

func trimExtraBlankLines(s string) string {
	lines := strings.Split(s, "\n")
	var compact []string
	lastBlank := false
	for _, line := range lines {
		blank := strings.TrimSpace(line) == ""
		if blank && lastBlank {
			continue
		}
		compact = append(compact, line)
		lastBlank = blank
	}
	return strings.Join(compact, "\n")
}
