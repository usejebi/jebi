package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/remote"
	"github.com/urfave/cli/v3"
)

type Clone struct {
	projectService projectService
	envService     envService
	secretService  secretService
	commitService  commitService
	cryptService   cryptService
	apiClient      apiClient
	slate          slate
	appService     appService
}

func NewCloneHandler(projectService projectService, envService envService, secretService secretService, commitService commitService, cryptService cryptService, apiClient apiClient, slate slate, appService appService) *Clone {
	return &Clone{
		projectService: projectService,
		envService:     envService,
		secretService:  secretService,
		commitService:  commitService,
		cryptService:   cryptService,
		apiClient:      apiClient,
		slate:          slate,
		appService:     appService,
	}
}

func (h *Clone) Handle(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 {
		return fmt.Errorf("usage: %s clone PROJECT_SLUG", core.AppName)
	}

	slug := cmd.Args().Get(0)
	h.slate.StartSpinner("Cloning project...")
	resp, err := h.apiClient.Clone(remote.CloneRequest{ProjectSlug: slug})
	if err != nil {
		h.slate.StopSpinner()
		return fmt.Errorf("failed to clone project: %w", err)
	}
	h.slate.UpdateSpinner("Setting up project...")
	data := resp.Data
	// Create hidden directory
	if err := h.appService.CreateAppDir(); err != nil {
		return err
	}

	h.slate.UpdateSpinner("Setting up environment...")
	// Create environment config
	if err := h.envService.CreateEnv(data.Environment.Name); err != nil {
		return err
	}

	// Set current environment
	if err := h.envService.SetCurrentEnv(data.Environment.Name); err != nil {
		return err
	}

	h.slate.UpdateSpinner("Setting up project...")
	projectId, err := h.projectService.SaveProjectConfig(data.Project.ID, data.Project.Name, data.Project.Description, data.Project.DefaultEnvironment)
	if err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	// Save the generated key
	if err := h.cryptService.SaveKey(data.Project.Key, projectId); err != nil {
		return fmt.Errorf("failed to save symmetric key: %w", err)
	}

	h.slate.UpdateSpinner("Importing commits...")
	var latestCommit *core.Commit
	for _, commit := range data.Commits {
		addedCommit, err := h.commitService.AddCommit(commit.ID, data.Environment.Name, commit.Message, commit.Author, commit.Changes, commit.Timestamp)
		if err != nil {
			h.slate.UpdateSpinner(fmt.Sprintf("Failed to import commit '%s': %v", commit.ID, err))
		}
		latestCommit = addedCommit
		h.slate.UpdateSpinner(fmt.Sprintf("Imported commit '%s'", addedCommit.ID))
	}
	h.commitService.UpdateRemoteHead(data.Environment.Name, latestCommit.ID)

	h.slate.UpdateSpinner("Importing secrets...")
	for _, secret := range data.Secrets {
		err = h.secretService.AddSecret(secret.Key, data.Environment.Name, secret)
		if err != nil {
			h.slate.UpdateSpinner(fmt.Sprintf("Failed to import secret '%s': %v", secret.Key, err))
		}
	}

	h.slate.StopSpinner()

	return nil
}
