package service

import (
	"github.com/No1ball/todo-app/internal/models"
	"github.com/No1ball/todo-app/internal/repository"
	"github.com/No1ball/todo-app/internal/todo"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(user models.SignInInput) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, id int) (todo.TodoList, error)
	Delete(userId, id int) error
	Update(userId, id int, input todo.UpdateListInput) error
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
	}
}
