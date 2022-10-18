package service

import (
	"github.com/No1ball/todo-app/internal/repository"
	"github.com/No1ball/todo-app/internal/todo"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) CreateItem(userId, listId int, input todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(listId, input)
}

func (s *TodoItemService) GetAllItem(userId, listId int) ([]todo.TodoItem, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(userId, listId)
}
