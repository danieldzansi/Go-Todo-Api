package handlers

import (
	"log"
	"net/http"

	"github.com/danieldzansi/auth-api/internal/models"
	"github.com/danieldzansi/auth-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TodoHandler struct {
	svc    service.TodoService
	tagSvc service.TagService
}

func NewTodoHandler(s service.TodoService, tagSvc service.TagService) *TodoHandler {
	return &TodoHandler{svc: s, tagSvc: tagSvc}
}

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

// CreateUser registers a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// CreateTodo creates a todo for a specific user
// Route: POST /users/:user_id/todos
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.svc.CreateTodo(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If tag IDs were provided, fetch and associate them now.
	if len(req.TagIDs) > 0 {
		tags, err := h.tagSvc.GetTagsByIDs(userID, req.TagIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or unauthorized tag IDs"})
			return
		}
		// Validate all requested tags were found
		if len(tags) != len(req.TagIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "one or more tag IDs not found"})
			return
		}
		if err := h.svc.AttachTags(todo, tags); err != nil {
			log.Printf("error attaching tags to todo %s: %v", todo.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to attach tags"})
			return
		}
		// Populate tags in response immediately
		todo.Tags = tags
	}

	c.JSON(http.StatusCreated, todo)
}

// GetAllTodos retrieves all todos (with their tags) for a specific user
// Route: GET /users/:user_id/todos
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	todos, err := h.svc.GetTodosByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

// UpdateTodoHandler updates a user's todo
// Route: PATCH /users/:user_id/todos/:id
func (h *TodoHandler) UpdateTodoHandler(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo ID"})
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.svc.UpdateTodo(userID, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": todo})
}

// DeleteTodoHandler deletes a user's todo
// Route: DELETE /users/:user_id/todos/:id
func (h *TodoHandler) DeleteTodoHandler(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo ID"})
		return
	}

	todo, err := h.svc.DeleteTodo(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Todo deleted successfully",
		"data":    todo,
	})
}

type TagHandler struct {
	svc service.TagService
}

func NewTagHandler(s service.TagService) *TagHandler {
	return &TagHandler{svc: s}
}

// CreateTag creates a tag for a specific user
// Route: POST /users/:user_id/tags
func (h *TagHandler) CreateTag(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.Tag
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.svc.CreateTag(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTagsByUser retrieves all tags for a specific user
// Route: GET /users/:user_id/tags
func (h *TagHandler) GetTagsByUser(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tags, err := h.svc.GetTagsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// DeleteTag deletes a tag for a specific user
// Route: DELETE /users/:user_id/tags/:id
func (h *TagHandler) DeleteTag(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	tagIDParam := c.Param("id")
	tagID, err := uuid.Parse(tagIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag ID"})
		return
	}

	err = h.svc.DeleteTag(userID, tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tag deleted successfully",
	})
}
