package security

import (
	"fmt"
	"time"

	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	UserID    string
	TokenType string
	ExpiresAt int64
	IssuedAt  int64
	TokenID   string
}

// GenerateTokenPair generates both access and refresh JWTs
func GenerateTokenPair(userID string, secret string) (string, string, error) {
	// Generate access token
	accessToken, err := generateAccessToken(userID, secret)
	if err != nil {
		return "", "", fmt.Errorf("error generating access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(userID, secret)
	if err != nil {
		return "", "", fmt.Errorf("error generating refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateAccessToken creates a short-lived JWT token for API access
func generateAccessToken(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 168).Unix(),
		"type":    "access",
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateRefreshToken creates a long-lived JWT token for obtaining new access tokens
func generateRefreshToken(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 720).Unix(), //30 days
		"type":    "refresh",
		"iat":     time.Now().Unix(),
		// "jti":     uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAccessToken is a convenience function for parsing access tokens
func ParseAccessToken(tokenString string, secret string) (*TokenClaims, error) {
	return ParseAndValidateToken(tokenString, secret, "access")
}

// ParseRefreshToken is a convenience function for parsing refresh tokens
func ParseRefreshToken(tokenString string, secret string) (*TokenClaims, error) {
	return ParseAndValidateToken(tokenString, secret, "refresh")
}

// ParseAndValidateToken parses a JWT token, validates it, and returns the claims
func ParseAndValidateToken(tokenString string, secret string, expectedType string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, inerr.ErrJwtValidation{
				Message: fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]),
			}
		}
		return []byte(secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, inerr.ErrJwtValidation{
					Message: "token has expired",
				}
			case ve.Errors&(jwt.ValidationErrorSignatureInvalid|jwt.ValidationErrorUnverifiable) != 0:
				return nil, inerr.ErrJwtValidation{
					Message: "invalid token signature",
				}
			default:
				return nil, inerr.ErrJwtValidation{
					Message: "invalid token format",
				}
			}
		}
		return nil, err
	}

	if !token.Valid {
		return nil, inerr.ErrJwtValidation{
			Message: "invalid token",
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, inerr.ErrJwtValidation{
			Message: "invalid claims format",
		}
	}

	// Validate token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != expectedType {
		return nil, inerr.ErrJwtValidation{
			Message: fmt.Sprintf("expected %s token, got %s", expectedType, tokenType),
		}
	}

	// Extract claims
	tokenClaims := &TokenClaims{
		UserID:    claims["user_id"].(string),
		TokenType: tokenType,
		ExpiresAt: int64(claims["exp"].(float64)),
		IssuedAt:  int64(claims["iat"].(float64)),
	}

	return tokenClaims, nil
}

func GenerateTokenWithClaims(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
