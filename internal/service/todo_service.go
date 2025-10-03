package service

import (
	"github.com/ahmadjafari86/go-todo-list/internal/models"
	"github.com/ahmadjafari86/go-todo-list/internal/repository"
)

type TodoService interface {
	CreateTodo(todo *models.Todo, ownerID uint) error
	ListTodos(ownerID uint) ([]models.Todo, error)
	GetTodo(id, ownerID uint) (*models.Todo, error)
	UpdateTodo(todo *models.Todo, ownerID uint) error
	ToggleComplete(id, ownerID uint) (*models.Todo, error)
	DeleteTodo(id, ownerID uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(todo *models.Todo, ownerID uint) error {
	todo.OwnerID = ownerID
	return s.repo.Create(todo)
}

func (s *todoService) ListTodos(ownerID uint) ([]models.Todo, error) {
	return s.repo.GetAll(ownerID)
}

func (s *todoService) GetTodo(id, ownerID uint) (*models.Todo, error) {
	return s.repo.GetByID(id, ownerID)
}

func (s *todoService) UpdateTodo(todo *models.Todo, ownerID uint) error {
	return s.repo.Update(todo, ownerID)
}

func (s *todoService) ToggleComplete(id, ownerID uint) (*models.Todo, error) {
	t, err := s.repo.GetByID(id, ownerID)
	if err != nil {
		return nil, err
	}
	t.Completed = !t.Completed
	if err := s.repo.Update(t, ownerID); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *todoService) DeleteTodo(id, ownerID uint) error {
	return s.repo.Delete(id, ownerID)
}
