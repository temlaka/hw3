package main

import (
	"context"
	"hw3/internal/client/coingecko"
	"hw3/internal/handler"
	"hw3/internal/repository/memory"
	"hw3/internal/scheduler"
	"hw3/internal/usecase/auth"
	"hw3/internal/usecase/crypto"
	"hw3/internal/usecase/schedule"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env not exist")
	}
	JWT_SECRET := os.Getenv("JWT_SECRET")
	COINGECKO_API_KEY := os.Getenv("COINGECKO_API_KEY")

	if JWT_SECRET == "" {
		log.Fatal("secret is none")
	}
	if COINGECKO_API_KEY == "" {
		log.Fatal("api key is none")
	}

	cgAPI := coingecko.NewClient("https://api.coingecko.com/api/v3", COINGECKO_API_KEY)

	userRepo := memory.NewUserRepo()
	cryptoRepo := memory.NewCryptoRepo(100)
	scheduleRepo := memory.NewScheduleRepo()

	authSvc := auth.NewService(userRepo, JWT_SECRET)
	authHandler := handler.NewAuthHandler(authSvc)

	cryptoSvc := crypto.NewService(cryptoRepo, cgAPI)
	cryptoHandler := handler.NewCryptoHandler(cryptoSvc)

	scheduleSvc := schedule.NewService(scheduleRepo, cryptoSvc)
	scheduleHandler := handler.NewScheduleHandler(scheduleSvc)

	mux := handler.Routes(authHandler, cryptoHandler, scheduleHandler, JWT_SECRET)

	// запуск schedule
	ctx, cancel := context.WithCancel(context.Background())
	go scheduler.UpdateAllPrice(ctx, scheduleSvc)
	defer cancel()

	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server didn't start, err: %v", err)
	}

}
