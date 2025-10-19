package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newCommitCommand(handler *handler.Commit) *cli.Command {
	return &cli.Command{
		Name:  "commit",
		Usage: fmt.Sprintf("Create a commit with a message: %s commit -m 'message'", AppName),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "message",
				Aliases:  []string{"m"},
				Usage:    "Commit message",
				Required: true,
			},
		},
		Action: handler.Handle,
	}
}
