package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/remote"
	"github.com/urfave/cli/v3"
)

type Push struct {
	projectService projectService
	envService     envService
	secretService  secretService
	commitService  commitService
	apiClient      apiClient
	slate          slate
}

func NewPushHandler(projectService projectService, envService envService, secretService secretService, commitService commitService, apiClient apiClient, slate slate) *Push {
	return &Push{
		projectService: projectService,
		envService:     envService,
		secretService:  secretService,
		commitService:  commitService,
		apiClient:      apiClient,
		slate:          slate,
	}
}

func (h *Push) Handle(ctx context.Context, cmd *cli.Command) error {
	// Get current environment
	currentEnv, err := h.envService.CurrentEnv()
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to get current environment: %v", err))
		return nil
	}

	// Get commits to push since remote HEAD
	commitsToPush, err := h.commitService.GetCommitsSinceRemoteHead(currentEnv)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to get commits to push: %v", err))
		return nil
	}

	// Check if there are any commits to push
	if len(commitsToPush) == 0 {
		fmt.Println("No new commits to push. Everything up-to-date.")
		return nil
	}

	h.slate.StartSpinner("Preparing to push commits...")

	// Load project configuration
	project, err := h.projectService.LoadProjectConfig()
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to load project: %v", err))
		return nil
	}
	// add metadata to commits
	for i := range commitsToPush {
		commitsToPush[i].ProjectID = project.ID
		commitsToPush[i].Environment = currentEnv
	}

	h.slate.UpdateSpinner("Loading project configuration...")

	// Get current HEAD to compute final state
	head, err := h.commitService.GetHead(currentEnv)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to get HEAD: %v", err))
		return nil
	}

	h.slate.UpdateSpinner("Computing final state...")
	// Get final state map from commits
	stateMap, err := h.commitService.ComputeState(currentEnv, head.LocalHead)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to compute final state: %v", err))
		return nil
	}

	h.slate.UpdateSpinner("Enriching final state with metadata...")
	// Enrich with metadata by retrieving all secrets from disk
	allSecrets, err := h.secretService.ListSecrets(project.ID, currentEnv)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to list secrets for metadata: %v", err))
		return nil
	}

	// Create a map of secrets by key for quick lookup
	secretsMap := make(map[string]core.Secret)
	for _, secret := range allSecrets {
		secretsMap[secret.Key] = secret
	}

	// Build final state by combining committed values with metadata
	var finalState []core.Secret
	for _, value := range stateMap {
		value.EnvironmentName = currentEnv
		value.ProjectId = project.ID
		finalState = append(finalState, value)
	}

	// Create environment object for API
	environment := core.Environment{
		Name:      currentEnv,
		ProjectID: project.ID,
	}

	// Create push request
	pushReq := remote.PushRequest{
		Project:        *project,
		Environment:    environment,
		Commits:        commitsToPush,
		FinalState:     finalState,
		RemoteHeadHash: head.RemoteHead,
	}

	h.slate.UpdateSpinner("Making push request to remote...")
	// Make the push request using injected API client
	response, err := h.apiClient.Push(pushReq)
	if err != nil {
		if errors.Is(err, remote.ErrProjectNameAlreadyExists) {
			h.slate.StopSpinner()
			h.slate.ShowError(fmt.Sprintf("Project with name '%s' already exists on remote.", project.Name))
			return nil
		}
		if errors.Is(err, remote.ErrUnauthorized) {
			h.slate.StopSpinner()
			h.slate.ShowError("Unauthorized access to remote server. Please try logging in.")
			return nil
		}
		h.slate.StopSpinner()
		h.slate.ShowError(fmt.Sprintf("Failed to push project: %v", err))
		return nil
	}

	// Update remote HEAD after successful push if we have commits
	if len(commitsToPush) > 0 {
		latestCommit := commitsToPush[len(commitsToPush)-1]
		err = h.commitService.UpdateRemoteHead(currentEnv, latestCommit.ID)
		if err != nil {
			h.slate.ShowWarning(fmt.Sprintf("Push succeeded but failed to update remote HEAD locally: %v", err))
		}
	}

	h.slate.StopSpinner()
	h.slate.WriteColoredText(response.Message, "")

	return nil
}
