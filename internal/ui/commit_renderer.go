package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jawahars16/jebi/internal/core"
)

// CommitRenderer provides shared functionality for rendering commits consistently
type CommitRenderer struct {
	slate Slate
}

// Slate interface for commit rendering (using concrete methods we need)
type Slate interface {
	WriteStyledText(text string, options StyleOptions)
	WriteIndentedText(text string, options StyleOptions)
}

// NewCommitRenderer creates a new commit renderer
func NewCommitRenderer(s Slate) *CommitRenderer {
	return &CommitRenderer{slate: s}
} // RenderCommit renders a single commit with beautiful styling
// showSpacing controls whether to add empty line before commit (for multiple commits)
func (r *CommitRenderer) RenderCommit(commit core.Commit, head *core.Head, showSpacing bool) {
	// Add spacing between commits if requested
	if showSpacing {
		fmt.Println() // Empty line between commits
	}

	// Determine HEAD status and styling
	var headLabel string
	var headColor string

	if head != nil {
		if commit.ID == head.RemoteHead && commit.ID == head.LocalHead {
			headLabel = "LOCAL & REMOTE HEAD"
			headColor = "82" // Green
		} else if commit.ID == head.LocalHead {
			headLabel = "LOCAL HEAD"
			headColor = "214" // Orange
		} else if commit.ID == head.RemoteHead {
			headLabel = "REMOTE HEAD"
			headColor = "196" // Red
		}
	}

	// Commit ID with HEAD label
	commitText := fmt.Sprintf("commit %s", commit.ID)
	if headLabel != "" {
		r.slate.WriteStyledText(commitText, StyleOptions{
			Color: "15", // White
			Bold:  true,
		})
		r.slate.WriteStyledText(fmt.Sprintf("(%s)", headLabel), StyleOptions{
			Color:      "0", // Black text
			Background: lipgloss.Color(headColor),
			Bold:       true,
			Padding:    []int{0, 1}, // Horizontal padding
		})
	} else {
		r.slate.WriteStyledText(commitText, StyleOptions{
			Color: "15", // White
			Bold:  true,
		})
	}

	// Author
	r.slate.WriteIndentedText(fmt.Sprintf("Author: %s", commit.Author), StyleOptions{
		Color: "248", // Gray
	})

	// Date - format relative time
	timeAgo := r.formatTimeAgo(commit.Timestamp)
	r.slate.WriteIndentedText(fmt.Sprintf("Date: %s (%s)",
		commit.Timestamp.Format("Mon Jan 2 15:04:05 2006"),
		timeAgo), StyleOptions{
		Color: "248", // Gray
	})

	// Commit message
	r.slate.WriteIndentedText(commit.Message, StyleOptions{
		Color:  "15", // White
		Italic: true,
		Margin: []int{0, 0, 0, 0},
	})

	// Show changes summary if available
	if len(commit.Changes) > 0 {
		r.displayChangesWithColors(commit.Changes)
	}
}

// RenderSingleCommit renders a single commit (for use in commit command)
func (r *CommitRenderer) RenderSingleCommit(commit core.Commit, head *core.Head) {
	r.slate.WriteStyledText("Created commit:", StyleOptions{
		Color:  "82", // Green
		Bold:   true,
		Margin: []int{0, 0, 1, 0}, // Bottom margin
	})
	r.RenderCommit(commit, head, false)
}

// displayChangesWithColors shows changes with appropriate colors for each type
func (r *CommitRenderer) displayChangesWithColors(changes []core.Change) {
	var adds, mods, dels int

	for _, change := range changes {
		switch change.Type {
		case core.ChangeTypeAdd:
			adds++
		case core.ChangeTypeModify:
			mods++
		case core.ChangeTypeRemove:
			dels++
		}
	}

	// Create colored parts
	var parts []string

	if adds > 0 {
		parts = append(parts, fmt.Sprintf("+%d", adds))
	}
	if mods > 0 {
		parts = append(parts, fmt.Sprintf("~%d", mods))
	}
	if dels > 0 {
		parts = append(parts, fmt.Sprintf("-%d", dels))
	}

	// Display with single line but indicate it has colors
	changesText := fmt.Sprintf("Changes: %s", strings.Join(parts, " "))
	r.slate.WriteIndentedText(changesText, StyleOptions{
		Color: "4", // Blue for the summary
	})
}

// formatTimeAgo returns a human-readable relative time
func (r *CommitRenderer) formatTimeAgo(timestamp time.Time) string {
	now := time.Now()
	diff := now.Sub(timestamp)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
}
