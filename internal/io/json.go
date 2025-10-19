package io

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJSONFile[T any](path string) (T, error) {
	var result T

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return zero value of T, not an error
			return result, nil
		}
		return result, fmt.Errorf("failed to read %q: %w", path, err)
	}

	if len(b) == 0 {
		// Return zero value of T, not an error
		return result, nil
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return result, fmt.Errorf("failed to parse %q: %w", path, err)
	}

	return result, nil
}

func WriteJSONToFile[T any](path string, data T) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	if err := os.WriteFile(path, out, 0600); err != nil {
		return fmt.Errorf("failed to write to %q: %w", path, err)
	}
	return nil
}
