package core

import (
	"fmt"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/io"
)

type changeRecordService struct {
	workingDir     string
	currentEnvPath string
}

func NewChangeRecordService(workingDir string) *changeRecordService {
	return &changeRecordService{
		workingDir:     workingDir,
		currentEnvPath: filepath.Join(workingDir, fmt.Sprintf(".%s", AppName), EnvDirPath, CurrentFileName),
	}
}

func (s *changeRecordService) AddChangeRecord(env, action, key, value, nonce string, noSecret bool) error {
	curr, err := io.ReadJSONFile[CurrentEnv](s.currentEnvPath)
	if err != nil {
		return fmt.Errorf("failed to read current environment: %w", err)
	}
	if curr.Changes == nil {
		curr.Changes = []Change{}
	}
	curr.Changes = append(curr.Changes, Change{
		Type:     ChangeType(action),
		Key:      key,
		Value:    value,
		Nonce:    nonce,
		NoSecret: noSecret,
	})
	curr.Changes = normalizeChanges(curr.Changes)

	err = io.WriteJSONToFile(s.currentEnvPath, curr)
	if err != nil {
		return fmt.Errorf("failed to write current environment: %w", err)
	}
	return nil
}

func (s *changeRecordService) ClearPendingChanges() error {
	curr, err := io.ReadJSONFile[CurrentEnv](s.currentEnvPath)
	if err != nil {
		return fmt.Errorf("failed to read current environment: %w", err)
	}
	curr.Changes = []Change{}

	err = io.WriteJSONToFile(s.currentEnvPath, curr)
	if err != nil {
		return fmt.Errorf("failed to write current environment: %w", err)
	}
	return nil
}

// normalizeChanges removes duplicate changes and applies conflict resolution
// Similar to the existing change normalization logic but for commitstore.Change
func normalizeChanges(changes []Change) []Change {
	// Track the latest change for each key
	changeMap := make(map[string]Change)

	for _, change := range changes {
		existing, exists := changeMap[change.Key]

		if !exists {
			changeMap[change.Key] = change
			continue
		}

		// Handle conflicts based on change types
		switch {
		case existing.Type == ChangeTypeAdd && change.Type == ChangeTypeRemove:
			// Add then remove = no change
			delete(changeMap, change.Key)
		case existing.Type == ChangeTypeRemove && change.Type == ChangeTypeAdd:
			// Remove then add = modify (if it existed before) or add (if new)
			changeMap[change.Key] = Change{
				Type:     ChangeTypeAdd, // Treat as add for simplicity
				Key:      change.Key,
				Value:    change.Value,
				Nonce:    change.Nonce,
				NoSecret: change.NoSecret,
			}
		default:
			// Later change wins
			changeMap[change.Key] = change
		}
	}

	// Convert map back to slice
	var normalized []Change
	for _, change := range changeMap {
		normalized = append(normalized, change)
	}

	return normalized
}

// func normalizeChanges(changes []Change) []Change {
// 	latest := make(map[string]Change)

// 	for _, c := range changes {
// 		prev, exists := latest[c.Key]
// 		if !exists {
// 			latest[c.Key] = c
// 			continue
// 		}

// 		switch prev.Action {
// 		case ActionAdd:
// 			if c.Action == ActionRemove {
// 				delete(latest, c.Key) // add+remove → no-op
// 			} else {
// 				latest[c.Key] = Change{Action: ActionAdd, Key: c.Key}
// 			}
// 		case ActionRemove:
// 			if c.Action == ActionAdd {
// 				latest[c.Key] = Change{Action: ActionUpdate, Key: c.Key} // remove+add → modify
// 			} else {
// 				latest[c.Key] = c
// 			}
// 		case ActionUpdate:
// 			if c.Action == ActionRemove {
// 				latest[c.Key] = c
// 			} else {
// 				latest[c.Key] = c // override with latest modify
// 			}
// 		default:
// 			latest[c.Key] = c
// 		}
// 	}

// 	result := make([]Change, 0, len(latest))
// 	for _, c := range latest {
// 		result = append(result, c)
// 	}
// 	return result
// }
