package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/federicoReghini/Chirpy/internal/database"
	"github.com/google/uuid"
)

// handlerCreateChirp creates a new chirp in the database.
// It reads the request body, validates the chirp, and if valid, inserts it into the database.
// If the chirp is invalid, it returns an error response.
// If the chirp is valid, it returns the created chirp with a 201 status code.
// It also censors bad words in the chirp body.
// The chirp body must not exceed 140 characters.
// The function uses the database.CreateChirpParams struct to validate the chirp parameters.
// It returns a JSON response with the created chirp or an error message.
// The function is part of the apiConfig struct which contains the database connection.
// It is registered as a handler for the "/chirps" endpoint with the POST method.
func (c *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// Check if length of the chirp is valid
	chirpValidated, isValid := handlerValidateChirp(w, req)
	if !isValid {
		return
	}

	chirp, err := c.db.CreateChirp(req.Context(), chirpValidated)

	if err != nil {
		marshalError(w, 500, err.Error())
		return
	}

	marshalOkJson(w, http.StatusCreated, chirp)
}

func (c *apiConfig) handlerGetChips(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	chirps, err := c.db.GetChirps(req.Context())

	if err != nil {
		marshalError(w, 500, err.Error())
		return
	}

	marshalOkJson(w, http.StatusOK, chirps)
}
func (c *apiConfig) handlerGetChipByID(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	uuid, err := uuid.Parse(req.PathValue("chirpID"))

	if err != nil {
		marshalError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := c.db.GetChirp(req.Context(), uuid)

	if err != nil {
		marshalError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	marshalOkJson(w, http.StatusOK, chirp)
}

// handlerValidateChirp validates the chirp request body and returns the chirp parameters if valid.
// If the chirp is invalid, it returns an error response and false.
// It checks for the length of the chirp body and censors bad words.
func handlerValidateChirp(w http.ResponseWriter, req *http.Request) (database.CreateChirpParams, bool) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := database.CreateChirpParams{}

	err := decoder.Decode(&params)
	if err != nil {
		marshalError(w, http.StatusInternalServerError, "Something went wrong")
		return database.CreateChirpParams{}, false
	}

	if len(params.Body) > 140 {
		marshalError(w, http.StatusBadRequest, "Chirp is too long")
		return database.CreateChirpParams{}, false

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
