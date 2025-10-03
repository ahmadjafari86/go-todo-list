package repository

import (
	"gorm.io/gorm"

	"github.com/ahmadjafari86/go-todo-list/internal/models"
)

type TodoRepository interface {
	Create(todo *models.Todo) error
	GetAll(ownerID uint) ([]models.Todo, error)
	GetByID(id uint, ownerID uint) (*models.Todo, error)
	Update(todo *models.Todo, ownerID uint) error
	Delete(id uint, ownerID uint) error
}

type GormTodoRepository struct {
	db *gorm.DB
}

func NewGormTodoRepository(db *gorm.DB) TodoRepository {
	return &GormTodoRepository{db: db}
}

func (r *GormTodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *GormTodoRepository) GetAll(ownerID uint) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Where("owner_id = ?", ownerID).Find(&todos).Error
	return todos, err
}

func (r *GormTodoRepository) GetByID(id uint, ownerID uint) (*models.Todo, error) {
	var t models.Todo
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormTodoRepository) Update(todo *models.Todo, ownerID uint) error {
	return r.db.Model(&models.Todo{}).
		Where("id = ? AND owner_id = ?", todo.ID, ownerID).
		Updates(todo).Error
}

func (r *GormTodoRepository) Delete(id uint, ownerID uint) error {
	return r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Todo{}).Error
}
