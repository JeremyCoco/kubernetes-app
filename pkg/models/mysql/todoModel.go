package mysql

import (
	"database/sql"
	"errors"

	"github.com/JeremyCoco/kubernetes-app/pkg/models"
)

type TodoModel struct {
	DB *sql.DB
}

func (tm *TodoModel) Insert(userId int, todo models.Todo) error {
	tx, err := tm.DB.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec(insertTodo, todo.Title, todo.CreatedAt, todo.Expires)
	if err != nil {
		tx.Rollback()
		return err
	}

	todoId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(insertUserTodo, userId, todoId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (tm *TodoModel) GetAll(userId int) ([]*models.Todo, error) {
	var todos []*models.Todo

	result, err := tm.DB.Query(selectAllTodos, userId)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		todo := &models.Todo{}

		err := result.Scan(&todo.Id, &todo.Title, &todo.Done, &todo.CreatedAt, &todo.Expires)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	err = result.Err()
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (tm *TodoModel) MarkAsComplete(userId, todoId int) error {
	tx, err := tm.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(selectUserTodo, userId, todoId)
	if err != nil {
		tx.Rollback()

		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecords
		}

		return err
	}

	_, err = tx.Exec(updateTodoToDone, todoId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (tm *TodoModel) Delete(userId, todoId int) error {
	tx, err := tm.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteUserTodo, userId, todoId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(deleteTodo, todoId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}
