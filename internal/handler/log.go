package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/ui"
	"github.com/urfave/cli/v3"
)

type Log struct {
	envService    envService
	commitService commitService
	slate         slate
}

func NewLogHandler(envService envService, commitService commitService, slate slate) *Log {
	return &Log{
		envService:    envService,
		commitService: commitService,
		slate:         slate,
	}
}

func (h *Log) Handle(ctx context.Context, cmd *cli.Command) error {
	env, err := h.envService.CurrentEnv()
	if err != nil {
		h.slate.WriteStyledText("Current environment is not set", ui.StyleOptions{
			Color: "178", // Amber
			Bold:  true,
		})
		h.slate.WriteIndentedText("Use `jebi env use <env name>` to set the current environment", ui.StyleOptions{
			Color:  "248", // Gray
			Italic: true,
		})
		return nil
	}

	commits, err := h.commitService.ListCommits(env)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		h.slate.WriteStyledText("No commits yet in this environment", ui.StyleOptions{
			Color:  "248", // Gray
			Italic: true,
		})
		h.slate.WriteIndentedText("Use `jebi commit -m \"your message\"` to create your first commit", ui.StyleOptions{
			Color:  "248", // Gray
			Italic: true,
		})
		return nil
	}

	head, err := h.commitService.GetHead(env)
	if err != nil {
		return err
	}

	// Show header
	h.slate.WriteStyledText(fmt.Sprintf("Commit History - Environment: %s", env), ui.StyleOptions{
		Color:  "82", // Light green
		Bold:   true,
		Margin: []int{0, 0, 1, 0}, // Bottom margin
	})

	// Display each commit
	renderer := ui.NewCommitRenderer(h.slate)
	for i, commit := range commits {
		renderer.RenderCommit(commit, head, i > 0) // Add spacing for all but first commit
	}

	// Show summary
	h.slate.WriteStyledText(fmt.Sprintf("Total: %d commits", len(commits)), ui.StyleOptions{
		Color:  "248", // Gray
		Italic: true,
		Margin: []int{1, 0, 0, 0}, // Top margin
	})

	return nil
}
