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

	userID, err := getUserIDFromValidateJWT(c, w, req)

	if err != nil {
		marshalError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check if length of the chirp is valid
	chirpValidated, isValid := handlerValidateChirp(w, req)
	if !isValid {
		return
	}

	chirpValidated.UserID = uuid.NullUUID{UUID: userID, Valid: true}

	chirp, err := c.db.CreateChirp(req.Context(), chirpValidated)

	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	marshalOkJson(w, http.StatusCreated, chirp)
}

func (c *apiConfig) handlerGetChips(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	chirps, err := c.db.GetChirps(req.Context())

	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
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
	decoder := json.NewDecoder(req.Body)
	params := database.CreateChirpParams{}

	err := decoder.Decode(&params)
	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
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
				if _, exist := badWords[strings.ToLower(word)]; exist {
					msg.WriteString("****" + " ")
				} else {
					msg.WriteString(word + " ")
				}
			}
			// Update params.Body with censored text and trim trailing space
			params.Body = strings.TrimSpace(msg.String())
		}

		return params, true
	}
}

// handlerDeleteChirp deletes a chirp from the database.
// It checks if the user is authorized to delete the chirp by comparing the user ID from
// the JWT token with the user ID of the chirp.
// If the user is authorized, it deletes the chirp and returns a 204 No Content response.
// If the user is not authorized, it returns a 403 Forbidden response.
// If the chirp does not exist, it returns a 404 Not Found response.
// It uses the database.DeleteChirpParams struct to validate the chirp parameters.
// The function is part of the apiConfig struct which contains the database connection.
// It is registered as a handler for the "/chirps/{chirpID}" endpoint with the DELETE method.
func (c *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	userID, err := getUserIDFromValidateJWT(c, w, req)
	if err != nil {
		marshalError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		marshalError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := c.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		marshalError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if userID == chirp.UserID.UUID {
		c.db.DeleteChirp(req.Context(), database.DeleteChirpParams{
			UserID: uuid.NullUUID{UUID: userID, Valid: true},
			ID:     chirp.ID,
		})
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}

}
