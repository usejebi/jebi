package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newExportCommand(handler *handler.Export) *cli.Command {
	return &cli.Command{
		Name:  "export",
		Usage: fmt.Sprintf("Export secrets in various formats: %s export [--env dev] [--format env]", core.AppName),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format (env, json, yaml, k8s)",
				Value:   "env",
			},
		},
		Action: handler.Handle,
	}

}
