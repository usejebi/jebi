package io

import (
	"bufio"
	"fmt"
	"strings"
)

// ToEnv converts a map to .env string
func ToEnv(secrets map[string]string, env string) string {
	lines := make([]string, 0, len(secrets))
	lines = append(lines, fmt.Sprintf("# Exported variables for environment: %s", env))
	for k, v := range secrets {
		lines = append(lines, k+"="+v)
	}
	return strings.Join(lines, "\n") + "\n"
}

// FromEnv parses .env text into a map
func FromEnv(data string) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result, scanner.Err()
}
