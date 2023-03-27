package borg

import (
	"errors"
	"strings"
)

func getVersion() (string, error) {
	borgmaticCmd := execCommand("borg", "--version")

	borgVersion, err := borgmaticCmd.Output()
	if err != nil {
		return "", errors.New("executable not found")
	}

	return strings.ReplaceAll(string(borgVersion), "\n", ""), nil
}
