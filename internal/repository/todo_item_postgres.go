package repository

import (
	"fmt"
	"github.com/No1ball/todo-app/internal/todo"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItem(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, input todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description, done) values ($1, $2, $3) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItemQuery, input.Title, input.Description, input.Done)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (item_id, list_id) values ($1, $2)", listItemsTable)
	_, err = tx.Exec(createListItemsQuery, itemId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done from %s ti INNER JOIN %s li on li.item_id = ti.id
    								INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listItemsTable, usersListTable)
	err := r.db.Select(&items, query, listId, userId)
	return items, err
}
