package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/urfave/cli/v3"
)

type Commit struct {
	envService          envService
	commitService       commitService
	changeRecordService changeRecordService
	slate               slate
}

func NewCommitHandler(
	envService envService,
	commitService commitService,
	changeRecordService changeRecordService,
	slate slate,
) *Commit {
	return &Commit{
		envService:          envService,
		commitService:       commitService,
		changeRecordService: changeRecordService,
		slate:               slate,
	}
}

func (h *Commit) Handle(ctx context.Context, cmd *cli.Command) error {
	msg := cmd.String("message")

	env, err := h.envService.CurrentEnv()
	if err != nil {
		return err
	}

	currentEnv, err := h.envService.GetCurrentEnv()
	if err != nil {
		return err
	}
	if len(currentEnv.Changes) == 0 {
		fmt.Printf("ℹ️  No changes to commit for environment '%s'\n", env)
		return nil
	}

	if err := h.commitService.AddCommit(env, msg, currentEnv.Changes); err != nil {
		return err
	}

	if err := h.changeRecordService.ClearPendingChanges(); err != nil {
		return err
	}

	var (
		output    string
		additions int
		deletions int
		changes   int
	)

	for _, change := range currentEnv.Changes {
		switch change.Action {
		case core.ActionAdd:
			additions++
		case core.ActionRemove:
			deletions++
		case core.ActionUpdate:
			changes++
		}
	}

	output += fmt.Sprintf("[environment `%s`] %s\n", env, msg)
	output += fmt.Sprintf(" %d additions(+), %d deletions(-), %d changes(~)\n", additions, deletions, changes)
	h.slate.RenderMarkdown(output)
	return nil
}
