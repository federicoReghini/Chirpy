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

func (c *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// Get user from JWT token (reuse auth logic)
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		marshalError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, c.apiKey)
	if err != nil {
		marshalError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Reuse the same request struct as createUser
	params := createUserBodyRequest{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		marshalError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Reuse password hashing logic from createUser
	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update user with same parameter structure
	user, err := c.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: params.Password,
	})

	if err != nil {
		marshalError(w, 500, "update user error: "+err.Error())
		return
	}

	// Reuse the same marshaling logic from createUser
	dat, err := json.Marshal(user)
	if err != nil {
		marshalError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

type loginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserWithToken struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

// handlerLogin handles user login requests.
// It checks the user's credentials and returns a JWT token if successful.
func (c *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	params := loginRequest{}

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params); err != nil {
		marshalError(w, http.StatusInternalServerError, "Something went wrong while decoding: "+err.Error())
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

	// Create Jwt
	token, err := auth.MakeJWT(user.ID, c.apiKey, time.Duration(60*60)*time.Second)

	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create refresh token
	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	usr := UserWithToken{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Token:        token,
		RefreshToken: refreshToken,
	}

	marshalOkJson(w, http.StatusOK, usr)

}

type token struct {
	Token string `json:"token"`
}

// handlerRefreshToken handles requests to refresh the JWT token using a refresh token.
// It checks the validity of the refresh token and returns a new JWT token if valid.
func (c *apiConfig) handlerRefreshToken(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	refreshToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if the refresh token is valid
	refreshTokenRecord, err := c.db.GetRefreshTokenByToken(req.Context(), refreshToken)

	if err != nil {
		marshalError(w, http.StatusNotFound, err.Error())
		return
	}

	if refreshTokenRecord.RevokedAt.Valid || refreshTokenRecord.ExpiresAt.Before(time.Now()) {
		marshalError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	tkn, err := auth.MakeJWT(refreshTokenRecord.UserID.UUID, c.apiKey, time.Duration(60*60)*time.Second)
	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	marshalOkJson(w, http.StatusOK, token{
		Token: tkn,
	})

}

func (c *apiConfig) handlerRefreshTokenRevoke(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	refreshToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		marshalError(w, http.StatusInternalServerError, err.Error())
		return
	}

	c.db.RevokeRefreshToken(req.Context(), refreshToken)

	w.WriteHeader(http.StatusNoContent)

}
