package handler

import (
	"net/http"

	"github.com/jbdoumenjou/go-rssaggregator/internal/api/json"
)

func ErrorHandler(w http.ResponseWriter, _ *http.Request) {
	json.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
