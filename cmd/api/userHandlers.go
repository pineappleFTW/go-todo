package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"lisheng/todo/pkg/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) loginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := models.Credentials{}
	json.NewDecoder(r.Body).Decode(&c)

	err := c.Validate()
	if err != nil {
		app.generateFailResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err := app.user.Authenticate(c.Email, c.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.generateErrorResponse(w, models.ErrInvalidCredentials, http.StatusInternalServerError)
			return
		}
		app.generateFailResponse(w, models.ErrInvalidCredentials, http.StatusInternalServerError)
		return
	}

	user, err := app.user.UserGetByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	accessToken, err := app.generateToken(user)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	refreshToken, err := app.generateRefreshToken(user)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	_, err = app.refreshTokens.RefreshTokenAdd(app.hashAccessToken(accessToken, user), refreshToken, user.ID)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	ut := &UserWithToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	app.generateSuccessResponse(w, ut)
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)

	err := u.Validate()
	if err != nil {
		app.generateFailResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err := app.user.UserSave(u.Name, u.Email, string(u.Password))
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	user, err := app.user.UserGetByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	accessToken, err := app.generateToken(user)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	refreshToken, err := app.generateRefreshToken(user)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	_, err = app.refreshTokens.RefreshTokenAdd(app.hashAccessToken(accessToken, user), refreshToken, user.ID)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	ut := &UserWithToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	app.generateSuccessResponse(w, ut)
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accessToken := r.Header.Get("Authorization")
	refreshToken := r.Header.Get("RefreshToken")
	if accessToken == "" || refreshToken == "" {
		app.generateErrorResponse(w, ErrNoTokenFound, http.StatusInternalServerError)
		return
	}

	newAccessToken, newRefreshToken, err := app.refreshToken(accessToken, refreshToken)
	if err != nil {
		app.infoLog.Println(err)
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	res := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}{
		newAccessToken,
		newRefreshToken,
	}

	app.generateSuccessResponse(w, res)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)

	err := u.Validate()
	if err != nil {
		app.generateFailResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err := app.user.UserSave(u.Name, u.Email, string(u.Password))
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	user, err := app.user.UserGetByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, user)
}

func (app *application) showUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	user, err := app.user.UserGetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.generateFailResponse(w, models.ErrNoRecord.Error(), http.StatusBadRequest)
			return
		}
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, user)
}

func (app *application) showAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := app.user.UserGetAll()
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, users)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	_, err = app.user.UserGetByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	receivedUser := models.User{}
	json.NewDecoder(r.Body).Decode(&receivedUser)

	err = receivedUser.ValidateUpdateUser()
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	id, err = app.user.UserUpdateByID(id, receivedUser.Name, receivedUser.Role, receivedUser.Active)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	u, err := app.user.UserGetByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, u)
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	_, err = app.user.UserGetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.generateFailResponse(w, models.ErrNoRecord.Error(), http.StatusBadRequest)
			return
		}
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	err = app.user.UserDeleteByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, nil)
}
