package core

import (
	"fmt"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/io"
)

type changeRecordService struct {
	workingDir string
}

func NewChangeRecordService(workingDir string) *changeRecordService {
	return &changeRecordService{
		workingDir: workingDir,
	}
}

func (s *changeRecordService) AddChangeRecord(env, action, key string) error {
	currentEnv := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), CurrentFileName)
	curr, err := io.ReadJSONFile[CurrentEnv](currentEnv)
	if err != nil {
		return fmt.Errorf("failed to read current environment: %w", err)
	}
	if curr.Changes == nil {
		curr.Changes = []Change{}
	}
	curr.Changes = append(curr.Changes, Change{
		Action: action,
		Key:    key,
	})
	curr.Changes = normalizeChanges(curr.Changes)

	err = io.WriteJSONToFile(currentEnv, curr)
	if err != nil {
		return fmt.Errorf("failed to write current environment: %w", err)
	}
	return nil
}

func (s *changeRecordService) GetPendingChanges() ([]Change, error) {
	currentEnv := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), CurrentFileName)
	curr, err := io.ReadJSONFile[CurrentEnv](currentEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to read current environment: %w", err)
	}
	return curr.Changes, nil
}

func (s *changeRecordService) ClearPendingChanges() error {
	currentEnv := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), CurrentFileName)
	curr, err := io.ReadJSONFile[CurrentEnv](currentEnv)
	if err != nil {
		return fmt.Errorf("failed to read current environment: %w", err)
	}
	curr.Changes = []Change{}

	err = io.WriteJSONToFile(currentEnv, curr)
	if err != nil {
		return fmt.Errorf("failed to write current environment: %w", err)
	}
	return nil
}

func normalizeChanges(changes []Change) []Change {
	latest := make(map[string]Change)

	for _, c := range changes {
		prev, exists := latest[c.Key]
		if !exists {
			latest[c.Key] = c
			continue
		}

		switch prev.Action {
		case ActionAdd:
			if c.Action == ActionRemove {
				delete(latest, c.Key) // add+remove → no-op
			} else {
				latest[c.Key] = Change{Action: ActionAdd, Key: c.Key}
			}
		case ActionRemove:
			if c.Action == ActionAdd {
				latest[c.Key] = Change{Action: ActionUpdate, Key: c.Key} // remove+add → modify
			} else {
				latest[c.Key] = c
			}
		case ActionUpdate:
			if c.Action == ActionRemove {
				latest[c.Key] = c
			} else {
				latest[c.Key] = c // override with latest modify
			}
		default:
			latest[c.Key] = c
		}
	}

	result := make([]Change, 0, len(latest))
	for _, c := range latest {
		result = append(result, c)
	}
	return result
}
