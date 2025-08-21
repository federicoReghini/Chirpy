package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("Expected non-empty hash")
	}

	if hash == password {
		t.Fatal("Hash should not equal original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Expected no error for correct password, got %v", err)
	}

	// Test wrong password
	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Fatal("Expected error for wrong password")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("Expected non-empty token")
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	// Create a valid token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Validate the token
	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Expected no error for valid token, got %v", err)
	}

	if validatedUserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, validatedUserID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := -time.Hour // Expired token

	// Create an expired token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Validate the expired token
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	wrongSecret := "wrong-secret"
	expiresIn := time.Hour

	// Create a token with one secret
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error for wrong secret")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "test-secret"
	invalidToken := "invalid.token.string"

	// Try to validate an invalid token
	_, err := ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Fatal("Expected error for invalid token")
	}
}

func TestValidateJWT_EmptyToken(t *testing.T) {
	tokenSecret := "test-secret"
	emptyToken := ""

	// Try to validate an empty token
	_, err := ValidateJWT(emptyToken, tokenSecret)
	if err == nil {
		t.Fatal("Expected error for empty token")
	}
}

func TestMakeJWT_ZeroExpiration(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Duration(0) // Immediate expiration

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Token should be created but might be expired immediately
	if token == "" {
		t.Fatal("Expected non-empty token")
	}
}

func TestGetBearerToken_ValidToken(t *testing.T) {

	headers := http.Header{}

	headers.Set("Authorization", "Bearer fsdjfksk3ek4wk4")

	_, err := GetBearerToken(headers)

	if err != nil {
		t.Fatal("Expected token: ", err)
	}

}
func TestGetBearerToken_NoBearerPrefix(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "fsdjfksk3ek4wk4")

	_, err := GetBearerToken(headers)

	if err == nil {
		t.Fatal("Expected error for authorization header without Bearer prefix")
	}
}

func TestGetBearerToken_NoAuthorizationHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)

	if err == nil {
		t.Fatal("Expected error for missing authorization header")
	}
}
func TestGetBearerToken_EmptyBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	_, err := GetBearerToken(headers)

	if err == nil {
		t.Fatal("Expected error for empty bearer token")
	}
}

func TestGetBearerToken_OnlyBearer(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer")

	_, err := GetBearerToken(headers)

	if err == nil {
		t.Fatal("Expected error for authorization header with only 'Bearer'")
	}
}
