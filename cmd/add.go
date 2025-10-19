package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newAddCommand(handler *handler.Add) *cli.Command {
	return &cli.Command{
		Name:   "add",
		Usage:  fmt.Sprintf("Add a new secret: %s add KEY VALUE", core.AppName),
		Action: handler.Handle,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-secret",
				Usage: "Do not encrypt the secret value before storing",
			},
		},
	}
}
