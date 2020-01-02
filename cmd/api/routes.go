package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/todo", app.createTodo)
	router.GET("/todo/:id", app.showTodo)
	router.GET("/todo", app.requireAuthentication(app.showAllTodos))
	router.PATCH("/todo/:id", app.updateTodo)
	router.DELETE("/todo/:id", app.deleteTodo)

	router.POST("/user", app.createUser)
	router.GET("/user/:id", app.showUser)
	router.GET("/users", app.showAllUsers)
	router.PATCH("/user/:id", app.updateUser)
	router.DELETE("/user/:id", app.deleteUser)

	return standardMiddleware.Then(router)
}
