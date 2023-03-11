package borgmatic

import (
	"errors"
	"os/exec"
	"strings"
)

func getVersion() (string, error) {
	borgmaticCmd := exec.Command("borgmatic", "--version")

	borgmaticVersion, err := borgmaticCmd.Output()
	if err != nil {
		return "", errors.New("borgmatic executable not found")
	}

	return strings.ReplaceAll(string(borgmaticVersion), "\n", ""), nil
}
