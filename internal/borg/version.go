package borg

import (
	"context"
	"errors"
	"strings"
)

func getVersion(ctx context.Context) (string, error) {
	borgmaticCmd := execCommand(ctx, "borg", "--version")

	borgVersion, err := borgmaticCmd.Output()
	if err != nil {
		return "", errors.New("executable not found")
	}

	return strings.TrimSpace(string(borgVersion)), nil
}
