package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	writeJSON(w, http.StatusOK, payload)
}

func (app *Config) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := readJSON(w, r, &requestPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		authenticate(w, requestPayload.Auth)
	default:
		errorJSON(w, errors.New("unknown action"))
	}
}

func authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.Marshal(a)

	request, err := http.NewRequest(http.MethodPost, "http://authenticationservice/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		errorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	writeJSON(w, http.StatusAccepted, payload)
}
