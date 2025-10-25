package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/urfave/cli/v3"
)

type Add struct {
	cryptService        cryptService
	envService          envService
	secretService       secretService
	changeRecordService changeRecordService
	projectService      projectService
	slate               slate
}

func NewAddHandler(
	projectService projectService,
	cryptService cryptService,
	envService envService,
	secretService secretService,
	changeRecordService changeRecordService,
	slate slate) *Add {
	return &Add{
		cryptService:        cryptService,
		envService:          envService,
		secretService:       secretService,
		changeRecordService: changeRecordService,
		slate:               slate,
		projectService:      projectService,
	}
}

func (s *Add) Handle(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 2 {
		return fmt.Errorf("usage: %s add KEY VALUE", core.AppName)
	}

	key := cmd.Args().Get(0)
	value := cmd.Args().Get(1)

	encryptionKey, err := s.cryptService.LoadKey()
	if err != nil {
		return fmt.Errorf("failed to retrieve encryption key: %w", err)
	}

	var secret core.Secret

	noSecret := cmd.Bool("no-secret")
	if noSecret {
		secret = core.Secret{
			Value: value,
		}
	} else {
		ciphertext, nonce, err := s.cryptService.Encrypt(encryptionKey, value)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}
		secret = core.Secret{
			Value: ciphertext,
			Nonce: nonce,
		}
	}

	env, err := s.envService.CurrentEnv()
	if err != nil {
		return fmt.Errorf("failed to get current environment: %w", err)
	}
	fmt.Printf("Current environment: %s\n", env)
	if err := s.secretService.AddSecret(key, env, secret); err != nil {
		if errors.Is(err, core.ErrSecretAlreadyExists) {
			s.slate.ShowError(fmt.Sprintf("secret with key '%s' already exists", key))
			return nil
		}
		return fmt.Errorf("failed to add secret: %w", err)
	}

	if err := s.changeRecordService.AddChangeRecord(env, core.ActionAdd, key, secret.Value); err != nil {
		return fmt.Errorf("failed to record change: %w", err)
	}

	fmt.Printf("✅ Secret '%s' added successfully!\n", key)
	return nil
}
