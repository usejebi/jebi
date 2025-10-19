package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/urfave/cli/v3"
)

type Remove struct {
	cryptService        cryptService
	envService          envService
	secretService       secretService
	changeRecordService changeRecordService
}

func NewRemoveHandler(
	cryptService cryptService,
	envService envService,
	secretService secretService,
	changeRecordService changeRecordService) *Remove {
	return &Remove{
		cryptService:        cryptService,
		envService:          envService,
		secretService:       secretService,
		changeRecordService: changeRecordService,
	}
}

func (s *Remove) Handle(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 {
		return fmt.Errorf("usage: %s remove KEY", core.AppName)
	}

	key := cmd.Args().Get(0)
	env, err := s.envService.CurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}

	if err := s.secretService.RemoveSecret(key, env); err != nil {
		return fmt.Errorf("failed to remove secret: %w", err)
	}

	if err := s.changeRecordService.AddChangeRecord(env, core.ActionRemove, key); err != nil {
		return fmt.Errorf("failed to record change: %w", err)
	}

	fmt.Printf("✅ Secret '%s' removed successfully!\n", key)
	return nil
}
