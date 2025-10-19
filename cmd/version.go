package cmd

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/urfave/cli/v3"
)

func newVersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show the current CLI version",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Printf("%s version %s", core.AppName, core.AppVersion)
			return nil
		},
	}
}
