package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newStatusCommand(handler *handler.Status) *cli.Command {
	return &cli.Command{
		Name:   "status",
		Usage:  fmt.Sprintf("Show the status of the current environment: %s status", core.AppName),
		Action: handler.Handle,
	}
}
