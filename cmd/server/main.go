package main

import (
	"encoding/json"
	"hw3/internal/domain"
	"hw3/pkg/token"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var Users []domain.User

var JWT_SECRET string = os.Getenv("JWT_SECRET")

func main() {
	if JWT_SECRET == "" {
		log.Fatal("secret is none")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", PostRegisterHandler)

	mux.HandleFunc("POST /auth/login", PostLoginHandler)
	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server didn't start, err: %v", err)
	}
}

func PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	newUserjson := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&newUserjson); err != nil {
		writeError(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	if newUserjson.Username == "" || newUserjson.Password == "" {
		writeError(w, "Некорректное заполнение имени/пароля", http.StatusBadRequest)
		return
	}
	for _, user := range Users {
		if newUserjson.Username == user.Username {
			writeError(w, "Некорректное заполнение имени", http.StatusConflict)
			return
		}
	}
	// req ^^

	//logic
	hash, err := bcrypt.GenerateFromPassword([]byte(newUserjson.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, "ошибка хеширования пароля", http.StatusInternalServerError)
		return

	}
	newUser := domain.User{
		Username:     newUserjson.Username,
		PasswordHash: hash,
	}
	Users = append(Users, newUser)
	//resp

	t, err := token.GenerateToken(newUser.Username, JWT_SECRET)
	if err != nil {
		writeError(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": t})
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userDatajson := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&userDatajson); err != nil {
		writeError(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	// req ^^

	//logic
	var f bool
	for _, user := range Users {
		if user.Username == userDatajson.Username {

			err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(userDatajson.Password))
			if err == bcrypt.ErrMismatchedHashAndPassword {
				continue
			} else if err == nil {
				f = true

			}
		}
	}
	if f == false {

		writeError(w, "Некорректный password", http.StatusUnauthorized)
		return
	}

	//resp

	t, err := token.GenerateToken(userDatajson.Username, JWT_SECRET)
	if err != nil {
		writeError(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": t})
}

func writeError(w http.ResponseWriter, text string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{"error": text})
}
