package main

import (
	"hw3/internal/client/coingecko"
	"hw3/internal/handler"
	"hw3/internal/repository/memory"
	"hw3/internal/usecase/auth"
	"hw3/internal/usecase/crypto"
	"log"
	"net/http"
	"os"
)

func main() {

	JWT_SECRET := os.Getenv("JWT_SECRET")
	COINGECKO_API_KEY := os.Getenv("COINGECKO_API_KEY")

	if JWT_SECRET == "" {
		log.Fatal("secret is none")
	}
	if COINGECKO_API_KEY == "" {
		log.Fatal("api key is none")
	}

	cgAPI := coingecko.NewClient("https://api.coingecko.com/api/v3", COINGECKO_API_KEY)

	repo := memory.NewUserRepo()
	cryptoRepo := memory.NewCryptoRepo(100)

	authSvc := auth.NewService(repo, JWT_SECRET)
	authHandler := handler.NewAuthHandler(authSvc)

	cryptoSvc := crypto.NewService(cryptoRepo, cgAPI)
	cryptoHandler := handler.NewCryptoHandler(cryptoSvc)

	mux := handler.Routes(authHandler, cryptoHandler, JWT_SECRET)

	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server didn't start, err: %v", err)
	}
}
