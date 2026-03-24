package main

import (
	"context"

	"github.com/slavatrudu/auth/internal/app"
	"github.com/slavatrudu/auth/internal/config"
	"github.com/slavatrudu/auth/internal/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	l := logger.New(cfg)

	application := app.New(&l, cfg)

	if err := application.Run(ctx); err != nil {
		l.Fatal().Err(err).Msg("error")
	}
}
