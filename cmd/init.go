package cmd

import (
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newInitCommand(handler *handler.Init) *cli.Command {
	return &cli.Command{
		Name:   "init",
		Usage:  "Initialize a new project",
		Action: handler.Handle,
	}
}
