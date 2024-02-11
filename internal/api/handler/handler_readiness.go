package handler

import (
	"net/http"

	"github.com/jbdoumenjou/go-rssaggregator/internal/api/json"
)

func ReadinessHandler(w http.ResponseWriter, _ *http.Request) {
	json.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
