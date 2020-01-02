package postgres

import (
	"database/sql"
	"errors"
	"lisheng/todo/pkg/models"
)

//TodoModel to hold db driver
type TodoModel struct {
	DB *sql.DB
}

//TodoSave saves into db and implement TodoStore interface
func (t *TodoModel) TodoSave(title, content string) (int, error) {

	stmt := `insert into todos (title, content, created)
	VALUES($1, $2, current_timestamp) returning id`

	var id int

	err := t.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

//TodoUpdateByID updates and save into db and implement TodoStore interface
func (t *TodoModel) TodoUpdateByID(id int, title, content string) (int, error) {

	stmt := `update todos set title = $1, content = $2 where id = $3 returning id`

	var returnedID int

	err := t.DB.QueryRow(stmt, title, content, id).Scan(&returnedID)
	if err != nil {
		return 0, err
	}
	return int(returnedID), nil
}

//TodoGetByID gets specific todo
func (t *TodoModel) TodoGetByID(id int) (*models.Todo, error) {

	stmt := `select id, title, content, created from todos 
	where id = $1`

	row := t.DB.QueryRow(stmt, id)
	td := &models.Todo{}

	err := row.Scan(&td.ID, &td.Title, &td.Content, &td.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return td, nil
}

//TodoGetAll gets all
func (t *TodoModel) TodoGetAll() ([]*models.Todo, error) {

	stmt := `SELECT id, title, content, created from todos
	order by created desc limit 50`

	rows, err := t.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []*models.Todo{}

	for rows.Next() {
		td := &models.Todo{}
		err = rows.Scan(&td.ID, &td.Title, &td.Content, &td.Created)
		if err != nil {
			return nil, err
		}
		todos = append(todos, td)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

//TodoDeleteByID delete by id
func (t *TodoModel) TodoDeleteByID(id int) error {

	stmt := `delete from todos where id = $1`

	_, err := t.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}