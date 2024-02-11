package handler

import (
	"net/http"

	"github.com/jbdoumenjou/go-rssaggregator/internal/api/respond"
)

func Readiness(w http.ResponseWriter, _ *http.Request) {
	respond.WithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
