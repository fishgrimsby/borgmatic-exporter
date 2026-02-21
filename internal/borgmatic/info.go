package borgmatic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/fishgrimsby/borgmatic-exporter/internal/logs"
)

func getInfo(ctx context.Context, config string) ([]InfoResult, error) {
	if config == "" {
		// Single execution with no -c flag (borgmatic default configs)
		return getInfoForConfig(ctx, "")
	}

	// Split by space and run each config in parallel
	configs := strings.Split(config, " ")

	var wg sync.WaitGroup
	resultsChan := make(chan []InfoResult, len(configs))

	for _, cfg := range configs {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			results, err := getInfoForConfig(ctx, strings.TrimSpace(c))
			if err != nil {
				logs.Logger.Error("failed to get info for config", "config", c, "error", err.Error())
				return
			}
			resultsChan <- results
		}(cfg)
	}

	wg.Wait()
	close(resultsChan)

	// Aggregate results
	var allResults []InfoResult
	for results := range resultsChan {
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

func getInfoForConfig(ctx context.Context, config string) ([]InfoResult, error) {
	var borgmaticCmd = execCommand(ctx, "borgmatic")
	if config != "" {
		borgmaticCmd.Args = append(borgmaticCmd.Args, "-c", config)
	}
	borgmaticCmd.Args = append(borgmaticCmd.Args, "info", "--json", "-v", "-1")

	var info []InfoResult

	cmdResult, err := borgmaticCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("unable to get borgmatic info for config %s: %w", config, err)
	}

	if err := json.Unmarshal([]byte(cmdResult), &info); err != nil {
		return nil, fmt.Errorf("unable to parse borgmatic info for config %s: %w", config, err)
	}

	return info, nil
}
