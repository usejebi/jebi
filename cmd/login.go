package cmd

import (
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/urfave/cli/v3"
)

func newLoginCommand(handler *handler.Login) *cli.Command {
	return &cli.Command{
		Name:   "login",
		Usage:  fmt.Sprintf("Login to jebi server via browser: %s login", core.AppName),
		Action: handler.Handle,
	}
}
