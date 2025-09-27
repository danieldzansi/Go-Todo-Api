package repository

import (
	"fmt"
	"time"

	models "github.com/danieldzansi/auth-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) error
	GetAllTodos() ([]models.Todo, error)
	UpdateTodo(userID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(userID, id uuid.UUID) (*models.Todo, error)
	GetTodosByUser(userID uuid.UUID) ([]models.Todo, error)
	AddTagsToTodo(todo *models.Todo, tags []models.Tag) error
	CompleteOverdueTodos(now time.Time) (int64, error)
}

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
}

type TagRepository interface {
	CreateTag(userID uuid.UUID, req *models.Tag) (*models.Tag, error)
	GetTagsByUser(userID uuid.UUID) ([]models.Tag, error)
	DeleteTag(userID, tagID uuid.UUID) error
	GetTagsByIDs(userID uuid.UUID, ids []uuid.UUID) ([]models.Tag, error)
}

type todoRepository struct {
	db *gorm.DB
}

type userRepository struct {
	db *gorm.DB
}

type tagRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (a *todoRepository) CreateTodo(todo *models.Todo) error {
	now := time.Now()
	todo.ID = uuid.New()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	if err := a.db.Create(todo).Error; err != nil {
		return fmt.Errorf("failed to create todo %w", err)
	}
	return nil
}

func (a *userRepository) CreateUser(user *models.User) error {
	now := time.Now()
	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := a.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (a *todoRepository) GetAllTodos() ([]models.Todo, error) {
	var todos []models.Todo

	if err := a.db.Order("created_at DESC").Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to get todo %w", err)
	}
	return todos, nil
}

func (a *todoRepository) UpdateTodo(userID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error) {
	var todo models.Todo
	if err := a.db.First(&todo, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}

	if err := a.db.Model(&todo).Updates(req).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

func (a *todoRepository) DeleteTodo(userID, id uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	if err := a.db.First(&todo, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	if err := a.db.Delete(&todo).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}
func (a *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := a.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (a *todoRepository) GetTodosByUser(userID uuid.UUID) ([]models.Todo, error) {
	var todos []models.Todo
	if err := a.db.Where("user_id = ?", userID).
		Preload("Tags").
		Order("created_at DESC").
		Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to get todos for user: %w", err)
	}
	return todos, nil
}

func (r *tagRepository) CreateTag(userID uuid.UUID, req *models.Tag) (*models.Tag, error) {
	tag := &models.Tag{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.Create(tag).Error; err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return tag, nil
}

func (r *tagRepository) GetTagsByUser(userID uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) DeleteTag(userID, tagID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", tagID, userID).Delete(&models.Tag{}).Error
}

func (r *tagRepository) GetTagsByIDs(userID uuid.UUID, ids []uuid.UUID) ([]models.Tag, error) {
	if len(ids) == 0 {
		return []models.Tag{}, nil
	}
	var tags []models.Tag
	if err := r.db.Where("user_id = ? AND id IN ?", userID, ids).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	return tags, nil
}

func (a *todoRepository) AddTagsToTodo(todo *models.Todo, tags []models.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	if err := a.db.Model(todo).Association("Tags").Append(tags); err != nil {
		return fmt.Errorf("failed to associate tags: %w", err)
	}
	return nil
}

func (a *todoRepository) CompleteOverdueTodos(now time.Time) (int64, error) {
	tx := a.db.Model(&models.Todo{}).
		Where("due_date IS NOT NULL AND due_date <= ? AND completed = ?", now, false).
		Updates(map[string]interface{}{
			"completed":  true,
			"updated_at": now,
		})
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to complete overdue todos: %w", tx.Error)
	}
	return tx.RowsAffected, nil
}
