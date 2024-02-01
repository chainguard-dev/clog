package init

import (
	"log/slog"

	"github.com/imjasonh/gcpslog"
)

// Set up structured logging
func init() { slog.SetDefault(slog.New(gcpslog.NewHandler())) }
