package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newCloneCommand(handler *handler.Clone) *cli.Command {
	return &cli.Command{
		Name:   "clone",
		Usage:  fmt.Sprintf("Clone an existing project: %s clone PROJECT_SLUG", core.AppName),
		Action: handler.Handle,
	}
}
