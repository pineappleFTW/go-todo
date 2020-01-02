package postgres

import (
	"database/sql"
	"errors"
	"lisheng/todo/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

//UserModel to hold db driver
type UserModel struct {
	DB *sql.DB
}

//UserSave create new user and save to db
func (u *UserModel) UserSave(name, email, password string) (int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `insert into users (name, email, hashed_password, role, created)
	VALUES($1, $2, $3, $4, current_timestamp) returning id `

	var id int
	err = u.DB.QueryRow(stmt, name, email, string(hashedPassword), models.USER).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

//UserUpdateByID update users by id (except password)
func (u *UserModel) UserUpdateByID(id int, name string, role int, active bool) (int, error) {

	stmt := `update users set name = $1, role = $2, active = $3 where id = $4 returning id`

	var returnedID int

	err := u.DB.QueryRow(stmt, name, role, active, id).Scan(&returnedID)
	if err != nil {
		return 0, err
	}

	return int(returnedID), nil
}

//UserGetByID get users by id
func (u *UserModel) UserGetByID(id int) (*models.User, error) {
	stmt := `select id, name, email, role, active, created from users where id = $1`

	row := u.DB.QueryRow(stmt, id)
	user := &models.User{}

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Active, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

//UserGetAll get all users
func (u *UserModel) UserGetAll() ([]*models.User, error) {

	stmt := `select id, name, email, role, active, created from users
	order by created desc limit 50`

	rows, err := u.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Active, &user.Created)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

//UserDeleteByID delete user by id
func (u *UserModel) UserDeleteByID(id int) error {

	stmt := `delete from users where id = $1`

	_, err := u.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

//UserIsPasswordMatched get hashedpassword of user
func (u *UserModel) UserIsPasswordMatched(id int, password string) (bool, error) {
	stmt := `select hashed_password from users where id = $1`
	row := u.DB.QueryRow(stmt, id)
	var hashedPassword []byte
	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, models.ErrInvalidCredentials
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, models.ErrInvalidCredentials
		}
		return false, err
	}

	return true, nil
}
