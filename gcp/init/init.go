package init

import (
	"log/slog"

	"github.com/chainguard-dev/clog/gcp"
)

// Set up structured logging at Info+ level.
// TODO: Make the level configurable by env var or flag; or just remove the init package.
func init() { slog.SetDefault(slog.New(gcp.NewHandler(slog.LevelInfo))) }
