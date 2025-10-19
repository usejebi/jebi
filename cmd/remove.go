package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newRemoveCommand(handler *handler.Remove) *cli.Command {
	return &cli.Command{
		Name:   "remove",
		Usage:  fmt.Sprintf("Remove an existing secret: %s remove KEY", core.AppName),
		Action: handler.Handle,
	}
}
