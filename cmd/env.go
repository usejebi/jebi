package cmd

import (
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

const AppName = "AppName" // change this once, affects all paths

func newEnvCommand(handler *handler.Env) *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "Manage environments (list, new, use, remove)",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Usage:   "List available environments",
				Action:  handler.HandleList,
				Aliases: []string{"ls"},
			},
			{
				Name:   "new",
				Usage:  "Create a new environment",
				Action: handler.HandleNew,
			},
			{
				Name:   "use",
				Usage:  "Switch current environment",
				Action: handler.HandleUse,
			},
			{
				Name:    "remove",
				Usage:   "Remove an environment",
				Action:  handler.HandleRemove,
				Aliases: []string{"rem", "rm", "delete", "del"},
			},
		},
	}
}
