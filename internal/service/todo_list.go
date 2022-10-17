package service

import (
	"github.com/No1ball/todo-app/internal/repository"
	"github.com/No1ball/todo-app/internal/todo"
)

type todoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *todoListService {
	return &todoListService{repo: repo}
}

func (s *todoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *todoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *todoListService) GetById(userId, id int) (todo.TodoList, error) {
	return s.repo.GetById(userId, id)
}
