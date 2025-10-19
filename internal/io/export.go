package io

import (
	"fmt"
	"strings"
)

// Export writes secrets in the requested format (env, json, yaml)
func Export(format string, secrets map[string]string, env string) (string, error) {
	switch strings.ToLower(format) {
	case "env", "dotenv":
		return ToEnv(secrets, env), nil
	default:
		return "", fmt.Errorf("unknown export format: %s", format)
	}
}
