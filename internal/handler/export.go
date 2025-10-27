package handler

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jawahars16/jebi/internal/io"
	"github.com/urfave/cli/v3"
)

type Export struct {
	envService     envService
	cryptService   cryptService
	projectService projectService
	slate          slate
}

func NewExportHandler(envService envService, cryptService cryptService, projectService projectService, slate slate) *Export {
	return &Export{
		envService:     envService,
		cryptService:   cryptService,
		projectService: projectService,
		slate:          slate,
	}
}

func (h *Export) Handle(ctx context.Context, cmd *cli.Command) error {
	format := cmd.String("output")
	env, err := h.envService.CurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}

	project, err := h.projectService.LoadProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	secrets, err := h.cryptService.LoadSecrets(project.ID, env)
	if err != nil {
		return fmt.Errorf("failed to load secrets: %w", err)
	}
	projectName := sanitizeK8sName(project.Name)
	output, err := io.Export(format, secrets, env, projectName)
	if err != nil {
		return fmt.Errorf("failed to export secrets: %w", err)
	}

	fmt.Println(output)
	return nil
}

func sanitizeK8sName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")

	// Remove any remaining invalid characters
	re := regexp.MustCompile(`[^a-z0-9-]`)
	name = re.ReplaceAllString(name, "")

	// Trim leading/trailing hyphens
	name = strings.Trim(name, "-")

	if len(name) == 0 {
		name = "default"
	}

	return name
}
