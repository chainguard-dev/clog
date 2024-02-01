package init

import (
	"log/slog"

	"github.com/chainguard-dev/clog/gcp"
)

// Set up structured logging
func init() { slog.SetDefault(slog.New(gcp.NewHandler())) }
