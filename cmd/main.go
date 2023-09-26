package main

import (
	"context"
	"fmt"
	"market/api"

	"market/api/handler"
	"market/config"
	"market/pkg/logger"
	"market/storage/postgres"
)

func main() {
	cfg := config.Load()
	log := logger.NewLogger("market-project", logger.LevelInfo)
	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		return
	}

	h := handler.NewHandler(strg, log)

	r := api.NewServer(h)
	r.Run(fmt.Sprintf(":%s", cfg.Port))
}
