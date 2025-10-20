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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Usage:   "Name of the project",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:    "description",
				Usage:   "Description of the project",
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:    "environment",
				Usage:   "Environment (dev/prod)",
				Aliases: []string{"e"},
			},
		},
	}
}
