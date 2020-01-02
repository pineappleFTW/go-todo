package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"lisheng/todo/pkg/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createTodo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	t := models.Todo{}
	json.NewDecoder(r.Body).Decode(&t)
	err := t.Validate()
	if err != nil {

		// app.clientError(w, err, http.StatusBadRequest)
		app.generateFailResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err := app.todo.TodoSave(t.Title, t.Content)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	td, err := app.todo.TodoGetByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, td)

}

func (app *application) showTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	td, err := app.todo.TodoGetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.generateFailResponse(w, models.ErrNoRecord.Error(), http.StatusBadRequest)
			return
		}
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, td)
}

func (app *application) showAllTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	todos, err := app.todo.TodoGetAll()
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, todos)
}

func (app *application) updateTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	t := models.Todo{}
	json.NewDecoder(r.Body).Decode(&t)

	err = t.Validate()
	if err != nil {
		app.generateFailResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err = app.todo.TodoUpdateByID(id, t.Title, t.Content)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	td, err := app.todo.TodoGetByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	app.generateSuccessResponse(w, td)
}

func (app *application) deleteTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.todo.TodoDeleteByID(id)
	if err != nil {
		app.generateErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	app.generateSuccessResponse(w, nil)
}
