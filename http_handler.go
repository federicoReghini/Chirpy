package main

import (
	"encoding/json"
	"github.com/federicoReghini/Chirpy/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type myError struct {
	Error string `json:"error"`
}

func handlerHealth(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) (chirpBodyRequest, bool) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := chirpBodyRequest{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		marshalError(w, 500, "Something went wrong")
		return chirpBodyRequest{}, false
	}

	if len(params.Body) > 140 {
		marshalError(w, 400, "Chirp is too long")
		return chirpBodyRequest{}, false

	} else {
		// not allowed words
		badWords := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
		uncensuredText := strings.Split(params.Body, " ")

		msg := strings.Builder{}

		if params.Body != "" {
			for _, word := range uncensuredText {
				if _, exist := badWords[word]; exist {
					msg.WriteString("****" + " ")
				} else {
					msg.WriteString(word + " ")
				}
			}
		} else {
			msg.WriteString(params.Body)
		}

		return params, true
	}
}

func marshalError(w http.ResponseWriter, statusCode int, msg string) {
	error := myError{
		Error: msg,
	}

	dat, err := json.Marshal(error)
	if err != nil {
		log.Printf("Error marshal JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}

// Marshal and write on http.ResponseWriter the status code and struct marshalled
// w is the http.ResponseWriter
// statusCode should be any status code for http status response (2xx,5xx,4xx etc..)

func marshalOkJson(w http.ResponseWriter, statusCode int, data any) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	dat, err := json.Marshal(data)
	ifErrCondition(w, err, 500, "")

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(dat)
	}
}

type createUserBodyRequest struct {
	Email string `json:"email"`
}

func (c *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	params := createUserBodyRequest{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		marshalError(w, 500, "Something went wrong")
		return
	}

	user, err := c.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Fatal(err)
	}

	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshal JSON: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

type chirpBodyRequest struct {
	Body    string    `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

func (c *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// Check if length of the chirp is valid
	chirpValidated, isValid := handlerValidateChirp(w, req)
	if !isValid {
		return
	}

	chirp, err := c.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   chirpValidated.Body,
		UserID: uuid.NullUUID{UUID: chirpValidated.User_id, Valid: true},
	})

	if err != nil {
		marshalError(w, 500, err.Error())
		return
	}

	marshalOkJson(w, http.StatusCreated, chirp)
}

// A wrapper for if err nil condition that also write on the http.ResponseWriter.
// w is the http.ResponseWriter
// statusCode should be any status code for http status response (2xx,5xx,4xx etc..)
// msg if valuated as "" will be set to Something went wrong
func ifErrCondition(w http.ResponseWriter, err error, statusCode int, msg string) {

	if err != nil {

		if msg == "" {
			msg = err.Error()
		}

		marshalError(w, statusCode, msg)
		return
	}
}
