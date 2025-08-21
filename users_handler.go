package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/federicoReghini/Chirpy/internal/auth"
	"github.com/federicoReghini/Chirpy/internal/database"
	"github.com/google/uuid"
)

type createUserBodyRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	params := createUserBodyRequest{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&params)

	if err != nil {
		marshalError(w, 500, "Something went wrong")
		return
	}

	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		marshalError(w, 500, err.Error())
		return
	}

	user, err := c.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: params.Password,
	})

	if err != nil {
		marshalError(w, 501, "create user error: "+err.Error())
		return
	}

	dat, err := json.Marshal(user)
	if err != nil {
		marshalError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

type loginRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds:omitempty"`
}

type UserWithToken struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

// handlerLogin handles user login requests.
// It checks the user's credentials and returns a JWT token if successful.
func (c *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	params := loginRequest{}

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params); err != nil {
		marshalError(w, http.StatusInternalServerError, "Something went wrong while decoding")
		return
	}

	user, err := c.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		marshalError(w, http.StatusNotFound, "User not found")
		return
	}

	// Check if psw is correct
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		marshalError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 60*60 {
		params.ExpiresInSeconds = 60 * 60
	}

	// Create Jwt
	token, err := auth.MakeJWT(user.ID, c.apiKey, time.Duration(params.ExpiresInSeconds)*time.Second)

	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	usr := UserWithToken{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     token,
	}

	marshalOkJson(w, http.StatusOK, usr)

}
