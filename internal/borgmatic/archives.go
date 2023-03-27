package borgmatic

import (
	"encoding/json"
	"errors"
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
		return nil, errors.New("unable to list archives")
	}

	json.Unmarshal([]byte(cmdResult), &archives)

	return archives, nil
}
