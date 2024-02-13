package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("oahdoaisdh@!7332asdhoh")

type UserAuth struct{}

type Credential struct {
	UserID int `json:"id"`
}

type Claims struct {
	Credential
	jwt.RegisteredClaims
}

func GetToken(credential Credential) (token string, err error) {
	expirationTime := time.Now().Add(time.Hour * 24) // one day expiration
	claims := Claims{
		Credential: credential,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString(secretKey)
	return
}

func VerifyToken(token string) (Claims, error) {
	var claims Claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return claims, err
	}

	if !jwtToken.Valid {
		return claims, errors.New("your token is not valid")
	}
	return claims, nil
}
