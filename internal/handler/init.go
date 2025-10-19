package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/urfave/cli/v3"
)

type Init struct {
	appService     appService
	projectService projectService
	envService     envService
	cryptService   cryptService
	slate          slate
}

func NewInitHandler(
	appService appService,
	projectService projectService,
	envService envService,
	cryptService cryptService,
	slate slate,
) *Init {
	return &Init{
		appService:     appService,
		projectService: projectService,
		envService:     envService,
		cryptService:   cryptService,
		slate:          slate,
	}
}

func (h *Init) Handle(ctx context.Context, cmd *cli.Command) error {

	h.slate.ShowHeader(`
🚀 jebi — Git for Secrets
Manage, version, and collaborate on secrets

This will initialize a new jebi project:
 • Create .jebi vault directory
 • Setup your first environment (dev/prod)
 • Generate and encrypt a symmetric key
	`)

	defaultName := "my-project"
	projectName := h.slate.PromptWithDefault("project name:", defaultName)
	projectDescription := h.slate.PromptWithDefault("project description:", "A new gfs project")
	environment := h.slate.PromptWithDefault("environment (dev/prod):", core.DefaultEnvironment)

	// Create hidden directory
	if err := h.appService.CreateAppDir(); err != nil {
		return err
	}

	// Create environment config
	if err := h.envService.CreateEnv(environment); err != nil {
		return err
	}

	// Set current environment
	if err := h.envService.SetCurrentEnv(environment); err != nil {
		return err
	}

	// Save project configuration
	if err := h.projectService.SaveProjectConfig(projectName, projectDescription); err != nil {
		return err
	}

	// Generate symmetric key
	encodedKey, err := h.cryptService.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate symmetric key: %w", err)
	}

	// Save the generated key
	if err := h.cryptService.SaveKey(encodedKey); err != nil {
		return fmt.Errorf("failed to save symmetric key: %w", err)
	}

	fmt.Printf("\n✅ Project '%s' initialized successfully!\n", projectName)
	return nil
}
