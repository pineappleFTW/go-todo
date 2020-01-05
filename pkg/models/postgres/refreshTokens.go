package postgres

import (
	"database/sql"
	"errors"
	"lisheng/todo/pkg/models"
)

//RefreshTokenModel to hold db driver
type RefreshTokenModel struct {
	DB *sql.DB
}

//RefreshTokenAdd adds new refresh token and identifer to db
func (r *RefreshTokenModel) RefreshTokenAdd(identifier, refreshToken string, userID int) (int, error) {

	stmt := `insert into refresh_tokens (identifier, token, user_id, created, updated)
	VALUES ($1, $2, $3, current_timestamp, current_timestamp) returning id`

	var id int

	err := r.DB.QueryRow(stmt, identifier, refreshToken, userID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

//RefreshTokenUpdateByID updates refresh token and identifer to db
func (r *RefreshTokenModel) RefreshTokenUpdateByID(id int, identifier, refreshToken string) (int, error) {

	stmt := `update refresh_tokens set identifier = $1, token = $2, updated = current_timestamp where id = $3 returning id`

	var returnedID int

	err := r.DB.QueryRow(stmt, identifier, refreshToken, id).Scan(&returnedID)
	if err != nil {
		return 0, err
	}

	return int(returnedID), nil
}

//RefreshTokenDeleteByID delete by id
func (r *RefreshTokenModel) RefreshTokenDeleteByID(id int) error {

	stmt := `delete from refresh_tokens where id = $1`

	_, err := r.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

//RefreshTokenVerify verify if identifier and refresh token matched
func (r *RefreshTokenModel) RefreshTokenVerify(identifier, refreshToken string, userID int) (*models.RefreshToken, error) {

	rt := &models.RefreshToken{}
	rt.User = models.User{}

	stmt := `select id, identifier, token, created, updated from refresh_tokens
	where identifier = $1 and token = $2 and user_id = $3`

	row := r.DB.QueryRow(stmt, identifier, refreshToken, userID)
	err := row.Scan(&rt.ID, &rt.Identifier, &rt.Token, &rt.Created, &rt.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	stmt = `select id, name, email, role, active, created from users where id = $1`
	row = r.DB.QueryRow(stmt, userID)

	err = row.Scan(&rt.User.ID, &rt.User.Name, &rt.User.Email, &rt.User.Role, &rt.User.Active, &rt.User.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return rt, nil
}
