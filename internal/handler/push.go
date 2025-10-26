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

	// Load project configuration
	project, err := h.projectService.LoadProjectConfig()
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to load project: %v", err))
		return nil
	}

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
	// add metadata to commits
	for i := range commitsToPush {
		commitsToPush[i].ProjectID = project.ID
		commitsToPush[i].Environment = currentEnv
	}

	// Get current HEAD to compute final state
	head, err := h.commitService.GetHead(currentEnv)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to get HEAD: %v", err))
		return nil
	}

	// Get final state map from commits
	stateMap, err := h.commitService.ComputeState(currentEnv, head.LocalHead)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Failed to compute final state: %v", err))
		return nil
	}

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
	h.slate.ShowHeader(fmt.Sprintf("Pushing %d commits with %d secrets...\n", len(commitsToPush), len(finalState)))

	// Make the push request using injected API client
	response, err := h.apiClient.Push(pushReq)
	if err != nil {
		if errors.Is(err, remote.ErrProjectNameAlreadyExists) {
			h.slate.ShowError(fmt.Sprintf("Project with name '%s' already exists on remote. Please rename your local project and try again.", project.Name))
			return nil
		}
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

	// Show success message based on response
	if response.IsFirstPush {
		fmt.Printf("✅ Created and pushed project '%s' with %d commits\n", response.Name, response.CommitsPushed)
	} else {
		fmt.Printf("✅ Pushed %d new commits to project '%s'\n", response.CommitsPushed, response.Name)
	}

	return nil
}
