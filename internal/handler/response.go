package handler

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, text string, status int) {

	writeJSON(w, status, map[string]string{"error": text})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(v)
}
