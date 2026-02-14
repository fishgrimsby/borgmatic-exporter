package borgmatic

import (
	"encoding/json"
	"fmt"
	"strings"
)

func getInfo(config string) ([]InfoResult, error) {
	var borgmaticCmd = execCommand("borgmatic")
	if config != "" {
		borgmaticCmd.Args = append(borgmaticCmd.Args, "-c")
		configs := strings.Split(config, " ")
		borgmaticCmd.Args = append(borgmaticCmd.Args, configs...)

	}
	borgmaticCmd.Args = append(borgmaticCmd.Args, "info", "--json")

	var info []InfoResult

	cmdResult, err := borgmaticCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("unable to get borgmatic info: %w", err)
	}

	if err := json.Unmarshal([]byte(cmdResult), &info); err != nil {
		return nil, fmt.Errorf("unable to parse borgmatic info: %w", err)
	}

	return info, nil
}
