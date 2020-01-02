package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

//TodoStore interface (easier to mock test)
type TodoStore interface {
	TodoGetByID(int) (*Todo, error)
	TodoGetAll() ([]*Todo, error)
	TodoSave(string, string) (int, error)
	TodoDeleteByID(int) error
	TodoUpdateByID(int, string, string) (int, error)
}

//Todo is exported
type Todo struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

//Validate incoming request
func (t Todo) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&t.Content, validation.Required, validation.Length(1, 100)),
	)
}
