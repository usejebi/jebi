package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/remote"
	"github.com/urfave/cli/v3"
)

type Push struct {
	projectService projectService
	envService     envService
	secretService  secretService
	apiClient      apiClient
	slate          slate
}

func NewPushHandler(projectService projectService, envService envService, secretService secretService, apiClient apiClient, slate slate) *Push {
	return &Push{
		projectService: projectService,
		envService:     envService,
		secretService:  secretService,
		apiClient:      apiClient,
		slate:          slate,
	}
}

func (h *Push) Handle(ctx context.Context, cmd *cli.Command) error {
	h.slate.ShowHeader("Pushing project to remote server...")

	// Load project configuration
	project, err := h.projectService.LoadProjectConfig()
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to load project: %v", err))
		return nil
	}

	currentEnv, err := h.envService.CurrentEnv()
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to get current environment: %v", err))
		return nil
	}
	currentEnvironmentToPPush := core.Environment{
		Name:      currentEnv,
		ProjectID: project.ID,
	}

	secrets, err := h.secretService.ListSecrets(project.ID, currentEnv)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to load secrets: %v", err))
		return nil
	}

	// Create push request with project details
	pushReq := remote.PushRequest{
		Project:     *project,
		Environment: currentEnvironmentToPPush,
		Secrets:     secrets,
	}

	// Make the push request using injected API client
	_, err = h.apiClient.Push(pushReq)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to push project: %v", err))
		return nil
	}

	// Show success message
	fmt.Printf("Successfully pushed changes\n")
	return nil
}
