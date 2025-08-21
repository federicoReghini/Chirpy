package main

import (
	"encoding/json"
	"net/http"

	"github.com/federicoReghini/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type upgradeRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (c *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	token, apiKeyError := auth.GetAPIKey(req.Header)

	if apiKeyError != nil {
		marshalError(w, http.StatusUnauthorized, apiKeyError.Error())
		return
	}

	if token != c.polkaKey {
		marshalError(w, http.StatusUnauthorized, "Invalid apiKey")
		return

	}

	params := upgradeRequest{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&params)

	if err != nil {
		marshalError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Event != "user.upgraded" {
		marshalError(w, http.StatusNoContent, "Unsupported event type")
		return
	}
	uuid, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		marshalError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	error := c.db.UpdateUserToPremium(req.Context(), uuid)

	if error != nil {
		marshalError(w, http.StatusNotFound, "User not found or already upgraded")
		return
	}
	marshalOkJson(w, http.StatusNoContent, nil)
}
