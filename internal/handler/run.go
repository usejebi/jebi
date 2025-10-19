package handler

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v3"
)

type Run struct {
	envService   envService
	cryptService cryptService
	slate        slate
}

func NewRunHandler(envService envService, cryptService cryptService, slate slate) *Run {
	return &Run{
		envService:   envService,
		cryptService: cryptService,
		slate:        slate,
	}
}

func (h *Run) Handle(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("no command provided; usage: jebi run -- <command> [args...]")
	}

	currentEnvenv, err := h.envService.CurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}

	secrets, err := h.cryptService.LoadSecrets(currentEnvenv)
	if err != nil {
		return fmt.Errorf("failed to load secrets: %w", err)
	}

	// Prepare the child process
	child := exec.Command(args[0], args[1:]...)
	child.Env = append(os.Environ(), flattenEnv(secrets)...)
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	child.Stdin = os.Stdin

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
