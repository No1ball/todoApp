package repository

import (
	"fmt"
	"github.com/No1ball/todo-app/internal/todo"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable)

	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListTable)

	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	getAllListQuery := fmt.Sprintf("select tl.id, tl.title, tl.description from %s tl INNER JOIN %s ul on tl.id = ul.list_id where ul.user_id = $1", todoListTable, usersListTable)
	err := r.db.Select(&lists, getAllListQuery, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId, id int) (todo.TodoList, error) {
	var list todo.TodoList

	getAllListQuery := fmt.Sprintf(`select tl.id, tl.title, tl.description from %s tl
                                       INNER JOIN %s ul on tl.id = ul.list_id where ul.user_id = $1 AND ul.list_id = $2`,
		todoListTable, usersListTable)
	err := r.db.Get(&list, getAllListQuery, userId, id)
	return list, err
}

func (r *TodoListPostgres) Delete(userId, id int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul where ul.list_id = tl.id AND ul.user_id = $1 AND ul.list_id = $2 ",
		todoListTable, usersListTable)
	_, err := r.db.Exec(query, userId, id)
	return err
}

func (r *TodoListPostgres) Update(userId, id int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE ul.list_id = tl.id AND ul.list_id = $%d AND ul.user_id = $%d",
		todoListTable, setQuery, usersListTable, argId, argId+1)
	args = append(args, id, userId)
	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)
	_, err := r.db.Exec(query, args...)
	return err
}
