package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newSetCommand(handler *handler.Set) *cli.Command {
	return &cli.Command{
		Name:   "set",
		Usage:  fmt.Sprintf("Set or update a secret: %s set KEY VALUE", core.AppName),
		Action: handler.Handle,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-secret",
				Usage: "Do not encrypt the secret value before storing",
			},
		},
	}
}
