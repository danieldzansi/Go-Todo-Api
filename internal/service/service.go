package service

import (
	"fmt"
	"log"
	"time"

	models "github.com/danieldzansi/auth-api/internal/models"
	"github.com/danieldzansi/auth-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TodoService interface {
	CreateTodo(userID uuid.UUID, req *models.CreateTodoRequest) (*models.Todo, error)
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(userID, id uuid.UUID) (*models.Todo, error)
	UpdateTodo(userID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error)
	DeleteTodo(userID, id uuid.UUID) (*models.Todo, error)
	GetTodosByUser(userID uuid.UUID) ([]models.Todo, error)
	AttachTags(todo *models.Todo, tags []models.Tag) error
	CompleteOverdueTodos(now time.Time) (int64, error)
}

type UserService interface {
	CreateUser(req *models.User) (*models.User, error)
	Userlogin(req *models.Login) (string, error)
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
	if err := s.repo.CreateTodo(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.repo.GetAllTodos()
}

func (s *todoService) GetTodoByID(userID, id uuid.UUID) (*models.Todo, error) {
	return s.repo.GetTodoByID(userID, id)
}

func (s *todoService) UpdateTodo(userID, id uuid.UUID, req *models.UpdateTodoRequest) (*models.Todo, error) {
	return s.repo.UpdateTodo(userID, id, req)
}

func (s *todoService) DeleteTodo(userID, id uuid.UUID) (*models.Todo, error) {
	return s.repo.DeleteTodo(userID, id)
}
func (s *userService) CreateUser(req *models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}
	if err := s.repo.CreateUser(user); err != nil {
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

func (s *todoService) CompleteOverdueTodos(now time.Time) (int64, error) {
	return s.repo.CompleteOverdueTodos(now)
}

func (s *userService) Userlogin(req *models.Login) (string, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		return "token", err
	}

	return tokenString, nil
}
