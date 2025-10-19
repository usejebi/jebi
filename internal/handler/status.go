package handler

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

type Status struct {
	envService envService
	slate      slate
}

func NewStatusHandler(envService envService, slate slate) *Status {
	return &Status{
		envService: envService,
		slate:      slate,
	}
}

func (h *Status) Handle(ctx context.Context, cmd *cli.Command) error {
	currentEnv, err := h.envService.GetCurrentEnv()
	if err != nil {
		return err
	}

	if len(currentEnv.Changes) == 0 {
		println("No pending changes")
		return nil
	}
	in := fmt.Sprintf("# Pending changes on environment `%s`", currentEnv.Env)

	header, err := h.slate.RenderMarkdown(in)
	if err != nil {
		return err
	}
	fmt.Print(header)

	h.slate.WriteStatus(currentEnv.Changes)
	return nil
}
