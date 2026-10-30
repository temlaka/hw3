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

	mux := http.NewServeMux()
	repo := memory.NewUserRepo()
	authSvc := auth.NewService(repo, JWT_SECRET)
	auth := handler.NewAuthHandler(authSvc)

	cryptoRepo := memory.NewCryptoRepo(100)
	cgAPI := coingecko.NewClient("https://api.coingecko.com/api/v3", COINGECKO_API_KEY)
	cryptoSvc := crypto.NewService(cryptoRepo, cgAPI)

	cryptoHandler := handler.NewCryptoHandler(cryptoSvc)

	mux.HandleFunc("POST /auth/register", auth.Register)
	mux.HandleFunc("POST /auth/login", auth.Login)
	mux.HandleFunc("GET /crypto", cryptoHandler.GetCryptos)
	mux.HandleFunc("POST /crypto", cryptoHandler.AddCrypto)
	mux.HandleFunc("GET /crypto/{symbol}", cryptoHandler.GetCryptoInfo)
	mux.HandleFunc("PUT /crypto/{symbol}/refresh", cryptoHandler.UpdateCryptoInfo)
	mux.HandleFunc("GET /crypto/{symbol}/history", cryptoHandler.GetCryptoHistory)
	mux.HandleFunc("GET /crypto/{symbol}/stats", cryptoHandler.GetStatsCrypto)
	mux.HandleFunc("DELETE /crypto/{symbol}", cryptoHandler.DeleteCrypto)

	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server didn't start, err: %v", err)
	}
}
