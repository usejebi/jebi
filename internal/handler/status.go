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

	in := fmt.Sprintf("On environment %s", currentEnv.Env)
	if len(currentEnv.Changes) == 0 {
		in += "\n(no pending changes)"
		fmt.Println(in)
		return nil
	}
	// in = fmt.Sprintf("Pending changes on environment %s", currentEnv.Env)

	fmt.Println(in)
	h.slate.WriteStatus(currentEnv.Changes)
	return nil
}
