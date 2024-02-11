package api

import "net/http"

func readinessHandler(w http.ResponseWriter, _ *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
