package borgmatic

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func getInfo(config string) ([]InfoResult, error) {
	var borgmaticCmd = exec.Command("borgmatic")
	if config != "" {
		borgmaticCmd.Args = append(borgmaticCmd.Args, "-c")
		configs := strings.Split(config, " ")
		borgmaticCmd.Args = append(borgmaticCmd.Args, configs...)

	}
	borgmaticCmd.Args = append(borgmaticCmd.Args, "info", "--json")

	var info []InfoResult

	cmdResult, err := borgmaticCmd.Output()
	if err != nil {
		return nil, errors.New("unable to get borgmatic info")
	}

	json.Unmarshal([]byte(cmdResult), &info)

	return info, nil
}
