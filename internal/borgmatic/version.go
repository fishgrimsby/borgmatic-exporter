package borgmatic

import (
	"errors"
	"strings"
)

func getVersion() (string, error) {
	borgmaticCmd := execCommand("borgmatic", "--version")

	borgmaticVersion, err := borgmaticCmd.Output()
	if err != nil {
		return "", errors.New("borgmatic executable not found")
	}

	return strings.ReplaceAll(string(borgmaticVersion), "\n", ""), nil
}
