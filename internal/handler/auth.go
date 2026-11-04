package handler

import (
	"encoding/json"
	"errors"
	"hw3/internal/domain"
	"hw3/internal/usecase/auth"
	"log"
	"net/http"
)

type AuthHandler struct {
	svc *auth.Service
}

type credentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := credentialsRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	t, err := h.svc.Register(req.Username, req.Password)
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, "invalid username or password format", http.StatusBadRequest)
		return
	case errors.Is(err, domain.ErrUserExists):
		writeError(w, "user already exists", http.StatusConflict)
		return
	case err != nil:
		log.Printf("register: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, tokenResponse{Token: t})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req := credentialsRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	t, err := h.svc.Login(req.Username, req.Password)
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, "invalid credentials", http.StatusUnauthorized)
		return
	case err != nil:
		log.Printf("login: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{Token: t})
}
