package borgmatic

import (
	"context"
	"errors"
	"strings"
)

func getVersion(ctx context.Context) (string, error) {
	borgmaticCmd := execCommand(ctx, "borgmatic", "--version")

	borgmaticVersion, err := borgmaticCmd.Output()
	if err != nil {
		return "", errors.New("borgmatic executable not found")
	}

	return strings.TrimSpace(string(borgmaticVersion)), nil
}
