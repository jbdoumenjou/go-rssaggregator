package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

func WithJSONError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with %d error: %v", code, msg)
	}
	WithJSON(w, code, map[string]string{"error": msg})
}

func WithJSON(w http.ResponseWriter, code int, payload any) {
	content, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(content)
}
