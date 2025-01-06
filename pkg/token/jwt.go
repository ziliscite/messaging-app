package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrFailedSigning = errors.New("failed creating token")
var ErrParsingToken = errors.New("failed parsing token")
var ErrExpiredToken = errors.New("token expired")

// Create token
//
// Returns token string, expiration, and error
func Create(id uint, exp int64, email, secretKey string) (string, time.Time, error) {
	expAt := time.Now().Add(time.Duration(exp) * time.Minute)

	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    id,
			"email": email,
			"exp":   expAt,
		},
	)

	tokenStr, err := t.SignedString([]byte(secretKey))
	if err != nil {
		return "", expAt, ErrFailedSigning
	}

	return tokenStr, expAt, nil
}

// Validate token
//
// Returns user id and username
func Validate(tokenStr, secretKey string) (uint, string, error) {
	key := []byte(secretKey)
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return 0, "", ErrParsingToken
	}

	if !token.Valid {
		return 0, "", ErrInvalidToken
	}

	expClaim, ok := claims["exp"].(float64)
	if !ok {
		return 0, "", ErrInvalidToken
	}

	if int64(expClaim) < time.Now().Unix() {
		return 0, "", ErrExpiredToken
	}

	return uint(claims["id"].(float64)), claims["username"].(string), nil
}
