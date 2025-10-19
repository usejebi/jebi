package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newLogCommand(handler *handler.Log) *cli.Command {
	return &cli.Command{
		Name:   "log",
		Usage:  fmt.Sprintf("Show commit history: %s log", AppName),
		Action: handler.Handle,
	}
}
