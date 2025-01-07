package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ziliscite/messaging-app/pkg/must"
	"os"
	"strconv"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrFailedSigning = errors.New("failed creating token")
var ErrParsingToken = errors.New("failed parsing token")
var ErrExpiredToken = errors.New("token expired")

type CustomClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

// Create token
//
// Returns token string, expiration, and error
func Create(id uint, exp int64, email, secretKey string) (string, time.Time, error) {
	now := time.Now()
	expAt := now.Add(time.Duration(exp) * time.Minute)

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    must.MustEnv(os.Getenv("APP_NAME")),
			Subject:   fmt.Sprintf("%d", id),
		},
		Email: email,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := t.SignedString([]byte(secretKey))
	if err != nil {
		return "", expAt, ErrFailedSigning
	}

	return tokenStr, expAt, nil
}

// Validate token
//
// Returns user id and email
func Validate(tokenStr, secretKey string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", ErrParsingToken
	}

	if !token.Valid {
		return 0, "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return 0, "", ErrInvalidToken
	}

	// Check if the token has expired
	if claims.ExpiresAt.Before(time.Now()) {
		return 0, "", ErrExpiredToken
	}

	// Validate the issuer
	if claims.Issuer != must.MustEnv(os.Getenv("APP_NAME")) {
		return 0, "", ErrInvalidToken
	}

	// Parse the subject (user ID) from the token
	id, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return 0, "", ErrInvalidToken
	}

	return uint(id), claims.Email, nil
}
