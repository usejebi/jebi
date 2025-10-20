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
	exists, err := h.appService.Exists()
	if err != nil {
		return err
	}
	if exists {
		h.slate.ShowWarning("A Jebi project is already initialized in this directory.\nRemove the .jebi folder and run 'jebi init' again to start fresh.")
		return nil
	}

	// Create hidden directory
	if err := h.appService.CreateAppDir(); err != nil {
		return err
	}

	h.slate.ShowHeader(`
🚀 jebi — Git for Secrets
Manage, version, and collaborate on secrets

This will initialize a new jebi project:
 • Create .jebi vault directory
 • Setup your first environment (dev/prod)
 • Generate and encrypt a symmetric key
	`)

	projectName := cmd.String("name")
	projectDescription := cmd.String("description")
	environment := cmd.String("environment")

	// Fallback to prompts only if not provided
	if projectName == "" {
		projectName = h.slate.PromptWithDefault("project name:", core.DefaultProjectName)
	}
	if projectDescription == "" {
		projectDescription = h.slate.PromptWithDefault("project description:", "A new jebi project to manage my secrets")
	}
	if environment == "" {
		environment = h.slate.PromptWithDefault("environment (dev/prod):", core.DefaultEnvironment)
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
	_, err = h.projectService.SaveProjectConfig(projectName, projectDescription)
	if err != nil {
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
