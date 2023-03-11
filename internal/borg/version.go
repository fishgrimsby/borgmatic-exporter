package borg

import (
	"errors"
	"os/exec"
	"strings"
)

func getVersion() (string, error) {
	borgmaticCmd := exec.Command("borg", "--version")

	borgVersion, err := borgmaticCmd.Output()
	if err != nil {
		return "", errors.New("executable not found")
	}

	return strings.ReplaceAll(string(borgVersion), "\n", ""), nil
}
