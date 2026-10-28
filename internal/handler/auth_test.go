package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hw3/internal/handler"
	"hw3/internal/repository/memory"
	"hw3/internal/usecase/auth"
	"hw3/pkg/token"
)

const secret = "test-secret"

func newServer() *httptest.Server {
	h := handler.NewAuthHandler(auth.NewService(memory.NewUserRepo(), secret))

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	return httptest.NewServer(mux)
}

func TestAuth(t *testing.T) {
	srv := newServer()
	defer srv.Close()

	// Порядок важен: регистрация alice идёт первой, остальные кейсы на неё опираются.
	cases := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantError  string // ожидаемый текст в {"error": ...}, пусто — проверяем токен
	}{
		{"register ok", "POST", "/auth/register", `{"username":"alice","password":"secret"}`, 201, ""},
		{"register duplicate", "POST", "/auth/register", `{"username":"alice","password":"other"}`, 409, "user already exists"},
		{"register empty password", "POST", "/auth/register", `{"username":"bob","password":""}`, 400, "username and password are required"},
		{"register empty username", "POST", "/auth/register", `{"username":"","password":"secret"}`, 400, "username and password are required"},
		{"register bad json", "POST", "/auth/register", `{"username":`, 400, "invalid JSON"},
		{"login ok", "POST", "/auth/login", `{"username":"alice","password":"secret"}`, 200, ""},
		{"login wrong password", "POST", "/auth/login", `{"username":"alice","password":"wrong"}`, 401, "invalid credentials"},
		{"login unknown user", "POST", "/auth/login", `{"username":"nobody","password":"secret"}`, 401, "invalid credentials"},
		{"login bad json", "POST", "/auth/login", `not json`, 400, "invalid JSON"},
		{"login wrong method", "GET", "/auth/login", ``, 405, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, srv.URL+tc.path, strings.NewReader(tc.body))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tc.wantStatus)
			}
			if tc.wantStatus == 405 {
				return // ответ 405 отдаёт сам mux, обычным текстом
			}
			if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			var body map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}

			if tc.wantError != "" {
				if body["error"] != tc.wantError {
					t.Errorf("error = %q, want %q", body["error"], tc.wantError)
				}
				return
			}

			sub, err := token.ParseToken(body["token"], secret)
			if err != nil {
				t.Fatalf("token is invalid: %v", err)
			}
			if sub != "alice" {
				t.Errorf("token sub = %q, want alice", sub)
			}
		})
	}
}
