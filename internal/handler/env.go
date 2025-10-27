package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/ui"
	"github.com/urfave/cli/v3"
)

type Env struct {
	envService envService
	slate      slate
}

func NewEnvHandler(envService envService, slate slate) *Env {
	return &Env{
		envService: envService,
		slate:      slate,
	}
}

func (h *Env) HandleList(ctx context.Context, cmd *cli.Command) error {
	envs, err := h.envService.ListEnvs()
	if err != nil {
		return err
	}
	current, err := h.envService.CurrentEnv()
	if err != nil {
		return err
	}
	h.slate.ShowList("Environments", envs, current)
	return nil
}

func (h *Env) HandleNew(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 {
		return fmt.Errorf("usage: %s env new <name>", core.AppName)
	}
	env := cmd.Args().Get(0)
	if err := h.envService.CreateEnv(env); err != nil {
		return err
	}
	if err := h.envService.SetCurrentEnv(env); err != nil {
		return err
	}
	h.slate.WriteIndentedText(fmt.Sprintf("Created environment '%s' and switched to it.\n", env), ui.StyleOptions{
		Color: "34", // Green
		Bold:  true,
	})
	return nil
}

func (h *Env) HandleUse(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 {
		return fmt.Errorf("usage: %s env use <name>", core.AppName)
	}
	env := cmd.Args().Get(0)
	if err := h.envService.SetCurrentEnv(env); err != nil {
		return err
	}
	h.slate.WriteIndentedText(fmt.Sprintf("Switched to '%s'\n", env), ui.StyleOptions{
		Color: "34", // Green
		Bold:  true,
	})
	h.HandleList(ctx, cmd)
	return nil
}

func (h *Env) HandleRemove(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 {
		return fmt.Errorf("usage: %s env remove <name>", core.AppName)
	}
	env := cmd.Args().Get(0)
	if err := h.envService.RemoveEnv(env); err != nil {
		return err
	}
	h.slate.RenderMarkdown(fmt.Sprintf("Removed environment `%s`", env))
	return nil
}
