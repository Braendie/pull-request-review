package app

import (
	"log/slog"

	"github.com/Braendie/pull-request-review/internal/config"
)

func Start() {
	cfg := config.MustLoad()

	logger := slog.Default()

	
}
