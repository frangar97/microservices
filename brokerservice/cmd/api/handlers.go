package main

import (
	"net/http"
)

func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	writeJSON(w, http.StatusOK, payload)
}
