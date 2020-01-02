package main

import (
	"lisheng/todo/pkg/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("todo")

type Claims struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type UserWithToken struct {
	Token string `json:"token"`
	*models.User
}

func (app *application) generateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (app *application) verifyToken(token string) (bool, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}

	if !tkn.Valid {
		return false, nil
	}

	return true, nil
}

func (app *application) refreshToken(existingToken string) (string, error) {

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(existingToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", err
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return "", err
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
