package main

import "net/http"

func errorHandler(w http.ResponseWriter, _ *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}
