package main

import (
	"hw3/internal/handler"
	"hw3/internal/repository/memory"
	"hw3/internal/usecase/auth"
	"log"
	"net/http"
	"os"
)

var JWT_SECRET string = os.Getenv("JWT_SECRET")

func main() {
	if JWT_SECRET == "" {
		log.Fatal("secret is none")
	}
	mux := http.NewServeMux()
	repo := memory.NewUserRepo()
	svc := auth.NewService(repo, JWT_SECRET)
	h := handler.NewAuthHandler(svc)
	mux.HandleFunc("POST /auth/register", h.Register)

	mux.HandleFunc("POST /auth/login", h.Login)
	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server didn't start, err: %v", err)
	}
}
