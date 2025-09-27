package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type Todo struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Title       string     `json:"title" db:"title" binding:"required"`
	Description string     `json:"description" db:"description"`
	Completed   bool       `json:"completed" db:"completed"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Tags        []Tag      `json:"tags" gorm:"many2many:todo_tags;"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username" binding:"required"`
	Todos     []Todo    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Tag struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"uniqueIndex:idx_user_tag,priority:2;not null"`
	Todos     []Todo    `json:"todos" gorm:"many2many:todo_tags;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoTag struct {
	TodoID uuid.UUID `json:"todo_id" db:"todo_id"`
	TagID  uuid.UUID `json:"tag_id" db:"tag_id"`
}

type CreateTodoRequest struct {
	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	Title       string      `json:"title" binding:"required"`
	Description string      `json:"description" binding:"required"`
	Completed   bool        `json:"completed "`
	TagIDs      []uuid.UUID `json:"tag_ids"` // Accepts string or array of UUIDs
}

type TodoResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UpdateTodoRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Completed   *bool      `json:"completed"`
	UserID      *string    `json:"user_id"`
}

func (Todo) TableName() string {
	return "my_todos"
}
