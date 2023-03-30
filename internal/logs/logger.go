package logs

import (
	"os"

	"golang.org/x/exp/slog"
)

var Logger *slog.Logger

func Configure(debug bool, format string) {
	var handler slog.Handler
	var opts slog.HandlerOptions

	if debug {
		opts = slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}

	switch format {
	case "json":
		handler = opts.NewJSONHandler(os.Stdout)
	default:
		handler = opts.NewTextHandler(os.Stdout)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}
