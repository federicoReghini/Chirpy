package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(psw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return err
	}

	return nil
}

// MakeJWT generates a JWT token for the given user ID with the specified expiration time.
// It returns the token as a string or an error if the token could not be created.
// The tokenSecret is used to sign the JWT.
// The expiresIn parameter specifies the duration after which the token will expire.
// The userID is the unique identifier for the user for whom the token is being created.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {

		return "", err
	}

	return signedToken, nil
}

// ValidateJWT checks the validity of a JWT token.
// It returns the user ID if the token is valid or an error if the token is invalid
// The tokenString is the JWT token to validate.
// The tokenSecret is used to verify the signature of the JWT.
// If the token is valid, it returns the user ID as a uuid.UUID.
// If the token is invalid, it returns an error.
// The function does not check the expiration of the token; it only verifies the signature.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil

	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok || !token.Valid {
		return uuid.Nil, jwt.ErrInvalidKey
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// GetBearerToken extracts the Bearer token from the Authorization header.
// It returns the token as a string or an error if the header is not formatted correctly.
func GetBearerToken(headers http.Header) (string, error) {

	bearer := headers.Get("Authorization")

	if bearer == "" {
		return "", errors.New("authorization header is missing")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(bearer, bearerPrefix) {

		return "", errors.New("Authorization header must start with 'Bearer '")
	}

	token := strings.TrimSpace(bearer[len(bearerPrefix):])
	if token == "" {
		return "", errors.New("Bearer token is empty")
	}

	return token, nil
}

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}
