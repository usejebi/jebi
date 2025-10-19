package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/io"
	"github.com/urfave/cli/v3"
)

type Export struct {
	envService    envService
	exportService cryptService
	slate         slate
}

func NewExportHandler(envService envService, exportService cryptService, slate slate) *Export {
	return &Export{
		envService:    envService,
		exportService: exportService,
		slate:         slate,
	}
}

func (h *Export) Handle(ctx context.Context, cmd *cli.Command) error {
	format := cmd.String("format")
	env, err := h.envService.CurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}

	secrets, err := h.exportService.LoadSecrets(env)
	if err != nil {
		return fmt.Errorf("failed to load secrets: %w", err)
	}
	output, err := io.Export(format, secrets, env)
	if err != nil {
		return fmt.Errorf("failed to export secrets: %w", err)
	}

	fmt.Println(output)
	return nil
}
