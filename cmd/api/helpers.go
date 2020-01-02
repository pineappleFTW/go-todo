package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, err error, status int) {
	if err != nil {
		http.Error(w, err.Error(), status)
	} else {
		http.Error(w, http.StatusText(status), status)
	}
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, nil, http.StatusNotFound)
}

func (app *application) generateSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	r := response{
		Status:  Success,
		Data:    data,
		Message: "",
		Code:    http.StatusOK,
	}
	res, err := json.Marshal(r)

	if err != nil {
		app.generateFailResponse(w, err, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s", res)
}

func (app *application) generateFailResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	r := response{
		Status:  Fail,
		Data:    data,
		Message: "",
		Code:    status,
	}
	res, _ := json.Marshal(r)

	http.Error(w, string(res), r.Code)
}

func (app *application) generateErrorResponse(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	r := response{
		Status:  Error,
		Data:    nil,
		Message: err.Error(),
		Code:    status,
	}

	res, _ := json.Marshal(r)

	http.Error(w, string(res), r.Code)
}
