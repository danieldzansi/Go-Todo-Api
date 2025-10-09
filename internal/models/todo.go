// package models

// import (
// 	"errors"
// 	"time"

// 	"github.com/google/uuid"
// )

// var ErrUserNotFound = errors.New("user not found")

// type Todo struct {
// 	ID          uuid.UUID  `json:"id" db:"id"`
// 	Title       string     `json:"title" db:"title" binding:"required"`
// 	Description string     `json:"description" db:"description"`
// 	Completed   bool       `json:"completed" db:"completed"`
// 	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
// 	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
// 	Tags        []Tag      `json:"tags" gorm:"many2many:todo_tags;"`
// 	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
// 	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
// }

// type User struct {
// 	ID        uuid.UUID `json:"id" db:"id"`
// 	Username  string    `json:"username" db:"username" binding:"required"`
// 	Todos     []Todo    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
// 	CreatedAt time.Time `json:"created_at" db:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
// }

// type Tag struct {
// 	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// 	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
// 	Name      string    `json:"name" gorm:"uniqueIndex:idx_user_tag,priority:2;not null"`
// 	Todos     []Todo    `json:"todos" gorm:"many2many:todo_tags;"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

// type TodoTag struct {
// 	TodoID uuid.UUID `json:"todo_id" db:"todo_id"`
// 	TagID  uuid.UUID `json:"tag_id" db:"tag_id"`
// }

// type CreateTodoRequest struct {
// 	UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
// 	Title       string      `json:"title" binding:"required"`
// 	Description string      `json:"description" binding:"required"`
// 	Completed   bool        `json:"completed "`
// 	TagIDs      []uuid.UUID `json:"tag_ids"` // Accepts string or array of UUIDs
// }

// type TodoResponse struct {
// 	Success bool        `json:"success"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data"`
// }

// type UpdateTodoRequest struct {
// 	Title       *string    `json:"title"`
// 	Description *string    `json:"description"`
// 	DueDate     *time.Time `json:"due_date"`
// 	Completed   *bool      `json:"completed"`
// 	UserID      *string    `json:"user_id"`
// }

//	func (Todo) TableName() string {
//		return "my_todos"
//	}
package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username  string    `json:"username" db:"username" binding:"required" gorm:"not null;uniqueIndex:uni_users_username"`
	Password  string    `json:"password" db:"password" binding:"required" `
	Todos     []Todo    `json:"todos,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"not null;default:now()"`
}

type Login struct {
	ID       uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string    `json:"username" db:"username" binding:"required"`
	Password string    `json:"password" db:"password" binding:"required" `
}

func (User) TableName() string {
	return "users"
}

type Todo struct {
	ID          uuid.UUID  `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title       string     `json:"title" db:"title" binding:"required" gorm:"not null"`
	Description string     `json:"description" db:"description"`
	Completed   bool       `json:"completed" db:"completed" gorm:"not null;default:false"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date" gorm:"type:timestamptz"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id" gorm:"type:uuid;not null"`
	Tags        []Tag      `json:"tags,omitempty" gorm:"many2many:todo_tags;"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"not null;default:now()"`
}

func (Todo) TableName() string {
	return "my_todos"
}

type Tag struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" db:"user_id" gorm:"type:uuid;not null"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Name      string    `json:"name" db:"name" gorm:"not null;uniqueIndex:uni_tags_user_id_name,composite:user_id"`
	Todos     []Todo    `json:"todos,omitempty" gorm:"many2many:todo_tags;"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"not null;default:now()"`
}

func (Tag) TableName() string {
	return "tags"
}

type TodoTag struct {
	TodoID uuid.UUID `json:"todo_id" db:"todo_id" gorm:"type:uuid;primaryKey"`
	TagID  uuid.UUID `json:"tag_id" db:"tag_id" gorm:"type:uuid;primaryKey"`
}

func (TodoTag) TableName() string {
	return "todo_tags"
}

type CreateTodoRequest struct {
	// UserID      uuid.UUID   `json:"user_id" binding:"required"`
	Title       string      `json:"title" binding:"required"`
	Description string      `json:"description"`
	Completed   bool        `json:"completed"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
}

type TodoResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type UpdateTodoRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Completed   *bool      `json:"completed"`
}
