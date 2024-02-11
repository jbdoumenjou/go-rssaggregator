package handler

import (
	"net/http"

	"github.com/jbdoumenjou/go-rssaggregator/internal/api/respond"
)

func Error(w http.ResponseWriter, _ *http.Request) {
	respond.WithJSONError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
