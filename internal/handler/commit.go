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
}

func NewCommitHandler(envService envService, commitService commitService, changeRecordService changeRecordService) *Commit {
	return &Commit{
		envService:          envService,
		commitService:       commitService,
		changeRecordService: changeRecordService,
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

	fmt.Printf("✅ [%s] Commit created for environment '%s'\n", core.AppName, env)
	return nil
}
