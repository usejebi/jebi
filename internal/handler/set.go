package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/urfave/cli/v3"
)

type Set struct {
	cryptService        cryptService
	envService          envService
	secretService       secretService
	changeRecordService changeRecordService
	projectService      projectService
}

func NewSetHandler(
	projectService projectService,
	cryptService cryptService,
	envService envService,
	secretService secretService,
	changeRecordService changeRecordService) *Set {
	return &Set{
		cryptService:        cryptService,
		envService:          envService,
		secretService:       secretService,
		changeRecordService: changeRecordService,
		projectService:      projectService,
	}
}

func (s *Set) Handle(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 2 {
		return fmt.Errorf("usage: %s set KEY VALUE", core.AppName)
	}

	key := cmd.Args().Get(0)
	value := cmd.Args().Get(1)

	encryptionKey, err := s.cryptService.LoadKey()
	if err != nil {
		return fmt.Errorf("failed to retrieve encryption key: %w", err)
	}

	ciphertext, nonce, err := s.cryptService.Encrypt(encryptionKey, value)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	env, err := s.envService.CurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}

	noSecret := cmd.Bool("no-secret")
	var secret core.Secret
	if noSecret {
		secret = core.Secret{
			Value: value,
		}
	} else {
		secret = core.Secret{
			Value: ciphertext,
			Nonce: nonce,
		}
	}
	var action string
	if action, err = s.secretService.SetSecret(key, env, secret); err != nil {
		return fmt.Errorf("failed to set secret: %w", err)
	}

	if err := s.changeRecordService.AddChangeRecord(env, action, key); err != nil {
		return fmt.Errorf("failed to record change: %w", err)
	}

	fmt.Printf("✅ Secret '%s' %s successfully!\n", key, action)
	return nil
}
