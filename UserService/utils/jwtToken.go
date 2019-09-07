package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Error notes
var (
	TokenExpired     = errors.New(" Token is expired")
	TokenNotValidYet = errors.New(" Token not active yet")
	TokenMalformed   = errors.New(" That's not even a token")
	TokenInvalid     = errors.New(" Couldn't handle this token:")
)

// Payload, contains some custom information
type CustomJWTClaims struct {
	Id int64 `json:"user_id"`
	jwt.StandardClaims
}

// Create JWT token
func CreateJWTToken(claims CustomJWTClaims, salt []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(salt)
}

// Parse JWT token
func ParseJWTToken(tokenString string, salt []byte) (*CustomJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return salt, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	} else if claims, ok := token.Claims.(*CustomJWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// Refresh JWT token
func RefreshJWTToken(tokenString string, salt []byte, survivalTime time.Duration) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims( tokenString, &CustomJWTClaims{},
		func(token *jwt.Token) (interface{}, error) { return salt, nil } )

	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomJWTClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(survivalTime).Unix()
		return CreateJWTToken(*claims, salt)
	}
	return "", TokenInvalid
}
