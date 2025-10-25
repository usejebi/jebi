package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newPushCommand(handler *handler.Push) *cli.Command {
	return &cli.Command{
		Name:   "push",
		Usage:  fmt.Sprintf("Push project to remote server: %s push", core.AppName),
		Action: handler.Handle,
	}
}
