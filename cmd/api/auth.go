package main

import (
	"errors"
	"lisheng/todo/pkg/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/speps/go-hashids"
)

var (
	ErrNoTokenFound        = errors.New("auth: missing access/refresh token found")
	ErrTokenIsStillValid   = errors.New("auth: token is still valid")
	ErrInvalidMatchedToken = errors.New("auth: invalid access/refresh token")
)

var jwtKey = []byte("todo")

type Claims struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type UserWithToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	*models.User
}

func (app *application) generateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Minute)
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

func (app *application) verifyToken(token string) (int, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !tkn.Valid {
		return 0, nil
	}

	return claims.UserID, nil
}

func (app *application) refreshToken(existingToken, refreshToken string) (string, string, error) {

	//Parse existing refresh token
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(existingToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	v, _ := err.(*jwt.ValidationError)

	if err != nil && v.Errors != jwt.ValidationErrorExpired {
		return "", "", err
	}

	if claims.ExpiresAt > time.Now().Unix() {
		return "", "", ErrTokenIsStillValid
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	//Parse refresh token
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	v, _ = err.(*jwt.ValidationError)

	if err != nil && v.Errors != jwt.ValidationErrorExpired {
		return "", "", err
	}

	user, err := app.user.UserGetByID(claims.UserID)
	if err != nil {
		return "", "", err
	}

	rt, err := app.refreshTokens.RefreshTokenVerify(app.hashAccessToken(existingToken, user), refreshToken, user.ID)
	if err != nil {
		app.infoLog.Println(err)
		return "", "", ErrInvalidMatchedToken
	}

	newRefreshToken := ""

	if claims.ExpiresAt < time.Now().Unix() {
		newrt, err := app.generateRefreshToken(user)
		if err != nil {
			return "", "", err
		}
		if err != nil {
			return "", "", err
		}
		newRefreshToken = newrt
	}

	if newRefreshToken != "" {
		_, err = app.refreshTokens.RefreshTokenUpdateByID(rt.ID, app.hashAccessToken(tokenString, user), newRefreshToken)
	} else {
		_, err = app.refreshTokens.RefreshTokenUpdateByID(rt.ID, app.hashAccessToken(tokenString, user), rt.Token)
	}

	if err != nil {
		return "", "", nil
	}

	return tokenString, newRefreshToken, nil
}

func (app *application) generateRefreshToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(60 * 24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func (app *application) hashAccessToken(existingToken string, user *models.User) string {
	hd := hashids.NewData()
	//get the last part as collision happens if the string is too long
	hd.Salt = strings.Split(existingToken, ".")[2]
	// hd.Salt = existingToken
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{user.ID})
	app.infoLog.Println(hd.Salt)
	app.infoLog.Println(e)
	return e
}
