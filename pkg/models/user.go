package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

//UserStore exported
type UserStore interface {
	UserGetByID(int) (*User, error)
	UserGetAll() ([]*User, error)
	UserSave(string, string, string) (int, error)
	UserDeleteByID(int) error
	UserUpdateByID(int, string, int, bool) (int, error)
	Authenticate(string, string) (int, error)
}

var (
	ADMIN     = 0
	MODERATOR = 1
	USER      = 2
)

//User is exported
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Role     int       `json:"role"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

//Validate incoming user request
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.Email, validation.Required, validation.Length(1, 50), is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 50)),
		// validation.Field(&u.Role, validation.Required, is.Int),
	)
}

//ValidateUpdateUser validate update user
func (u User) ValidateUpdateUser() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&u.Role, validation.Required),
		// validation.Field(&u.Active, validation.Required),
	)
}

//Credentials used when logging in
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

//Validate incoming user login request
func (c Credentials) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, validation.Length(1, 50), is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 50)),
	)
}
