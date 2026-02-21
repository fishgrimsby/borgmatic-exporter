package borgmatic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/fishgrimsby/borgmatic-exporter/internal/logs"
)

func getArchives(ctx context.Context, config string) ([]ListResult, error) {
	if config == "" {
		// Single execution with no -c flag (borgmatic default configs)
		return getArchivesForConfig(ctx, "")
	}

	// Split by space and run each config in parallel
	configs := strings.Split(config, " ")

	var wg sync.WaitGroup
	resultsChan := make(chan []ListResult, len(configs))

	for _, cfg := range configs {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			results, err := getArchivesForConfig(ctx, strings.TrimSpace(c))
			if err != nil {
				logs.Logger.Error("failed to get archives for config", "config", c, "error", err.Error())
				return
			}
			resultsChan <- results
		}(cfg)
	}

	wg.Wait()
	close(resultsChan)

	// Aggregate results
	var allResults []ListResult
	for results := range resultsChan {
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

func getArchivesForConfig(ctx context.Context, config string) ([]ListResult, error) {
	var borgmaticCmd = execCommand(ctx, "borgmatic")
	if config != "" {
		borgmaticCmd.Args = append(borgmaticCmd.Args, "-c", config)
	}
	borgmaticCmd.Args = append(borgmaticCmd.Args, "list", "--json", "-v", "-1")

	var archives []ListResult

	cmdResult, err := borgmaticCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("unable to list archives for config %s: %w", config, err)
	}

	if err := json.Unmarshal([]byte(cmdResult), &archives); err != nil {
		return nil, fmt.Errorf("unable to parse archives list for config %s: %w", config, err)
	}

	return archives, nil
}
