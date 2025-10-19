package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newRunCommand(handler *handler.Run) *cli.Command {
	return &cli.Command{
		Name:   "run",
		Usage:  fmt.Sprintf("Run a command in the current environment: %s run -- <command> [args...]", core.AppName),
		Action: handler.Handle,
	}
}
