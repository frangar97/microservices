package main

import (
	"net/http"

	"github.com/frangar97/loggerservice/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) writeLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = readJSON(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	writeJSON(w, http.StatusAccepted, resp)
}
