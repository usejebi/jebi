package handler

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v3"
)

type Run struct {
	envService     envService
	cryptService   cryptService
	projectService projectService
	slate          slate
}

func NewRunHandler(envService envService, cryptService cryptService, projectService projectService, slate slate) *Run {
	return &Run{
		envService:     envService,
		cryptService:   cryptService,
		projectService: projectService,
		slate:          slate,
	}
}

func (h *Run) Handle(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("no command provided; usage: jebi run -- <command> [args...]")
	}

	project, err := h.projectService.LoadProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	currentEnv, err := h.envService.CurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}

	secrets, err := h.cryptService.LoadSecrets(project.ID, currentEnv)
	if err != nil {
		return fmt.Errorf("failed to load secrets: %w", err)
	}

	// Build the full shell command string
	commandLine := strings.Join(args, " ")

	// Spawn using a shell to support npm/yarn/next, etc.
	child := exec.Command("/bin/sh", "-c", commandLine)
	child.Env = append(os.Environ(), flattenEnv(secrets)...)
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	child.Stdin = os.Stdin

	h.slate.ShowHeader("Jebi is spawning your command with secrets injected at runtime.\nThese secrets are available only to this process and its children.")

	err = child.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}

	return err
}

func flattenEnv(m map[string]string) []string {
	env := make([]string, 0, len(m))
	for k, v := range m {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}
