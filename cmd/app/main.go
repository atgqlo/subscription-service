package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"subscriptons-service/internal/config"
	"subscriptons-service/internal/handlers"
	"subscriptons-service/internal/repository/postgres"
	"subscriptons-service/internal/server"
	"subscriptons-service/internal/service"
)

func main() {
	cfg := config.Load()
	logger := log.New(os.Stdout, "app: ", log.LstdFlags)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.PostgresURL())
	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Close()
	repo := postgres.NewSubscriptionRepository(pool)
	svc := service.NewSubscriptionService(repo, logger)
	h := handlers.NewSubscriptionHandler(svc, logger)
	srv := server.NewServer(h, logger)
	if err := srv.Run(":8080"); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}

}
