package service

import (
	"log"
	"time"

	models "github.com/danieldzansi/auth-api/internal/models"
	"github.com/danieldzansi/auth-api/internal/repository"
	"github.com/google/uuid"
)

type TodoService interface {
	CreateTodo(userID uuid.UUID, req *models.CreateTodoRequest) (*models.Todo, error)
	GetAllTodos() ([]models.Todo, error)
	UpdateTodo(userID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(userID, id uuid.UUID) (*models.Todo, error)
	GetTodosByUser(userID uuid.UUID) ([]models.Todo, error)
	AttachTags(todo *models.Todo, tags []models.Tag) error
	CompleteOverdueTodos(now time.Time) (int64, error)
}

type UserService interface {
	CreateUser(req *models.User) (*models.User, error)
}

type TagService interface {
	CreateTag(userID uuid.UUID, req *models.Tag) (*models.Tag, error)
	GetTagsByUser(userID uuid.UUID) ([]models.Tag, error)
	DeleteTag(userID, tagID uuid.UUID) error
	GetTagsByIDs(userID uuid.UUID, ids []uuid.UUID) ([]models.Tag, error)
}

type todoService struct {
	repo repository.TodoRepository
}

type tagService struct {
	repo repository.TagRepository
}

func NewTodoService(r repository.TodoRepository) TodoService {
	return &todoService{repo: r}
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func NewTagService(r repository.TagRepository) TagService {
	return &tagService{repo: r}
}

func (s *todoService) CreateTodo(userID uuid.UUID, req *models.CreateTodoRequest) (*models.Todo, error) {
	log.Printf("Creating todo with userID=%s", userID.String())
	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
		UserID:      userID,
	}
	// Create the todo first
	if err := s.repo.CreateTodo(todo); err != nil {
		return nil, err
	}

	// If there are tag IDs, we need access to tag repository functions. We only have TodoRepository here.
	// Strategy: type assert underlying repository if it also implements TagRepository methods is not possible cleanly.
	// Instead, association will be handled at handler level OR via a separate service. To keep it simple without changing constructor wiring,
	// we will perform tag association in a follow-up step in the handler using a new TagService passed there.
	return todo, nil
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.repo.GetAllTodos()
}

func (s *todoService) UpdateTodo(userID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error) {
	return s.repo.UpdateTodo(userID, id, req)
}

func (s *todoService) DeleteTodo(userID, id uuid.UUID) (*models.Todo, error) {
	return s.repo.DeleteTodo(userID, id)
}
func (s *userService) CreateUser(req *models.User) (*models.User, error) {
	user := &models.User{
		Username: req.Username,
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *todoService) GetTodosByUser(userID uuid.UUID) ([]models.Todo, error) {
	return s.repo.GetTodosByUser(userID)
}

func (s *tagService) CreateTag(userID uuid.UUID, req *models.Tag) (*models.Tag, error) {
	tag := &models.Tag{
		Name:   req.Name,
		UserID: userID,
	}
	createdTag, err := s.repo.CreateTag(userID, tag)
	if err != nil {
		return nil, err
	}
	return createdTag, nil
}

func (s *tagService) GetTagsByUser(userID uuid.UUID) ([]models.Tag, error) {
	return s.repo.GetTagsByUser(userID)
}

func (s *tagService) DeleteTag(userID, tagID uuid.UUID) error {
	return s.repo.DeleteTag(userID, tagID)
}

func (s *tagService) GetTagsByIDs(userID uuid.UUID, ids []uuid.UUID) ([]models.Tag, error) {
	return s.repo.GetTagsByIDs(userID, ids)
}

// AttachTags associates tags with a todo (no-op if tags slice empty)
func (s *todoService) AttachTags(todo *models.Todo, tags []models.Tag) error {
	if len(tags) == 0 {
		return nil
	}
	if err := s.repo.AddTagsToTodo(todo, tags); err != nil {
		log.Printf("failed attaching %d tags to todo %s: %v", len(tags), todo.ID, err)
		return err
	}
	return nil
}

// CompleteOverdueTodos sets completed=true for all todos whose due_date has passed
func (s *todoService) CompleteOverdueTodos(now time.Time) (int64, error) {
	return s.repo.CompleteOverdueTodos(now)
}
