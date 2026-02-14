package borgmatic

import (
	"encoding/json"
	"fmt"
	"strings"
)

func getArchives(config string) ([]ListResult, error) {
	var borgmaticCmd = execCommand("borgmatic")
	if config != "" {
		borgmaticCmd.Args = append(borgmaticCmd.Args, "-c")
		configs := strings.Split(config, " ")
		borgmaticCmd.Args = append(borgmaticCmd.Args, configs...)

	}
	borgmaticCmd.Args = append(borgmaticCmd.Args, "list", "--json")

	var archives []ListResult

	cmdResult, err := borgmaticCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("unable to list archives: %w", err)
	}

	if err := json.Unmarshal([]byte(cmdResult), &archives); err != nil {
		return nil, fmt.Errorf("unable to parse archives list: %w", err)
	}

	return archives, nil
}
