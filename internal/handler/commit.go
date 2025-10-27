package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/ui"
	"github.com/urfave/cli/v3"
)

type Commit struct {
	envService          envService
	commitService       commitService
	changeRecordService changeRecordService
	userService         userService
	secretService       secretService
	projectService      projectService
	slate               slate
}

func NewCommitHandler(
	envService envService,
	commitService commitService,
	changeRecordService changeRecordService,
	userService userService,
	secretService secretService,
	projectService projectService,
	slate slate,
) *Commit {
	return &Commit{
		envService:          envService,
		commitService:       commitService,
		changeRecordService: changeRecordService,
		userService:         userService,
		secretService:       secretService,
		projectService:      projectService,
		slate:               slate,
	}
}

func (h *Commit) Handle(ctx context.Context, cmd *cli.Command) error {
	msg := cmd.String("message")

	env, err := h.envService.CurrentEnv()
	if err != nil {
		return err
	}

	// Get pending changes from the existing change tracking system
	currentEnv, err := h.envService.GetCurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}

	if len(currentEnv.Changes) == 0 {
		h.slate.WriteStyledText("No changes to commit", ui.StyleOptions{
			Color:   "178", // Yellow/amber
			Padding: []int{0, 0, 0, 2},
		})
		return nil
	}

	// Get the best available author information (authenticated user or system username)
	author := h.userService.GetCommitAuthor()

	// Show commit author info (helpful for users to know what identity is being used)
	if user, err := h.userService.LoadCurrentUser(); err != nil || user == nil {
		// User is not logged in, show that we're using system username
		h.slate.WriteStyledText(fmt.Sprintf("Committing as: %s (system user)", author), ui.StyleOptions{
			Color:   "82", // Light green
			Italic:  true,
			Padding: []int{0, 0, 1, 2},
		})
		h.slate.WriteStyledText("💡 Use 'jebi login' to commit with your authenticated identity", ui.StyleOptions{
			Color:   "180", // Light yellow
			Italic:  true,
			Padding: []int{0, 0, 0, 2},
		})
	}

	// Create commit using commitstore
	commit, err := h.commitService.AddCommit(env, msg, author, currentEnv.Changes)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Clear uncommitted changes after successful commit using the existing service
	err = h.changeRecordService.ClearPendingChanges()
	if err != nil {
		return fmt.Errorf("failed to clear uncommitted changes: %w", err)
	}

	// Get HEAD information for consistent rendering
	head, err := h.commitService.GetHead(env)
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Use the shared commit renderer for consistent formatting
	renderer := ui.NewCommitRenderer(h.slate)
	renderer.RenderSingleCommit(*commit, head)
	return nil
}

// displayCommit renders the newly created commit in a beautiful format
