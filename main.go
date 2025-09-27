// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"

// 	"github.com/danieldzansi/auth-api/internal/database"
// 	"github.com/danieldzansi/auth-api/internal/handlers"
// 	"github.com/danieldzansi/auth-api/internal/repository"
// 	"github.com/danieldzansi/auth-api/internal/service"
// )

// func main() {

// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using system environment variables")
// 	}

// 	db, err := database.ConnectGorm()
// 	if err != nil {
// 		log.Fatal("failed to connect to database:", err)
// 	}

// 	todoRepo := repository.NewTodoRepository(db)
// 	todoSvc := service.NewTodoService(todoRepo)
// 	todoHandler := handlers.NewTodoHandler(todoSvc)

// 	userRepo := repository.NewUserRepository(db)
// 	userSvc := service.NewUserService(userRepo)
// 	userHandler := handlers.NewUserHandler(userSvc)

// 	gin.SetMode(os.Getenv("GIN_MODE"))
// 	router := gin.Default()
// 	router.Use(gin.Logger())
// 	router.Use(gin.Recovery())

// 	api := router.Group("/api/v1")
// 	{
// 		api.GET("/health", func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{
// 				"status":  "ok",
// 				"message": "Todo API is running",
// 			})
// 		})

// 		todos := api.Group("/todos")
// 		{
// 			todos.POST("/", todoHandler.CreateTodo)
// 			todos.GET("/", todoHandler.GetAllTodos)
// 			todos.PATCH("/:id", todoHandler.UpdateTodoHandler)
// 			todos.DELETE("/:id", todoHandler.DeleteTodoHandler)

// 		}
// 		users := api.Group("/users")
// 		{
// 			users.POST("/", userHandler.CreateUser) // POST /api/v1/users
// 		}

// 	}

// 	port := os.Getenv("SERVER_PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	log.Printf("Server starting on port %s", port)
// 	if err := router.Run(":" + port); err != nil {
// 		log.Fatal("Failed to start server:", err)
// 	}
// }

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/danieldzansi/auth-api/internal/database"
	"github.com/danieldzansi/auth-api/internal/handlers"
	models "github.com/danieldzansi/auth-api/internal/models"
	"github.com/danieldzansi/auth-api/internal/repository"
	"github.com/danieldzansi/auth-api/internal/service"
)

func main() {
	// --- Load environment variables ---
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// --- Database connection ---
	db, err := database.ConnectGorm()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// --- Ensure required extension and tables exist ---
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Printf("warning: failed to ensure uuid-ossp extension: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Todo{}, &models.Tag{}, &models.TodoTag{}); err != nil {
		log.Fatal("failed to auto-migrate database schema:", err)
	}

	// --- Repository / Service / Handler wiring ---
	todoRepo := repository.NewTodoRepository(db)
	todoSvc := service.NewTodoService(todoRepo)
	// Tag service must be created before wiring todo handler now that it needs tag association
	tagRepo := repository.NewTagRepository(db)
	tagSvc := service.NewTagService(tagRepo)
	tagHandler := handlers.NewTagHandler(tagSvc)
	todoHandler := handlers.NewTodoHandler(todoSvc, tagSvc)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	// tagRepo / tagSvc already initialized above

	// --- Gin setup ---
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// --- API routes ---
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Todo API is running",
			})
		})

		// User routes
		users := api.Group("/users")
		{
			// Create a user
			users.POST("/", userHandler.CreateUser)

			// Nested todos routes for a specific user
			userTodos := users.Group("/:user_id/todos")
			{
				userTodos.POST("/", todoHandler.CreateTodo)             // POST   /api/v1/users/:user_id/todos
				userTodos.GET("/", todoHandler.GetAllTodos)             // GET    /api/v1/users/:user_id/todos
				userTodos.PATCH("/:id", todoHandler.UpdateTodoHandler)  // PATCH  /api/v1/users/:user_id/todos/:id
				userTodos.DELETE("/:id", todoHandler.DeleteTodoHandler) // DELETE /api/v1/users/:user_id/todos/:id
			}

			// --- Background scheduler to complete overdue todos ---
			intervalStr := os.Getenv("DUE_CHECK_INTERVAL")
			if intervalStr == "" {
				intervalStr = "1m" // default to every minute
			}
			interval, err := time.ParseDuration(intervalStr)
			if err != nil {
				log.Printf("Invalid DUE_CHECK_INTERVAL '%s', falling back to 1m", intervalStr)
				interval = time.Minute
			}
			ticker := time.NewTicker(interval)
			go func() {
				log.Printf("Overdue checker running every %s", interval)
				for range ticker.C {
					now := time.Now().UTC()
					n, err := todoSvc.CompleteOverdueTodos(now)
					if err != nil {
						log.Printf("Overdue checker error: %v", err)
						continue
					}
					if n > 0 {
						log.Printf("Overdue checker: marked %d todo(s) completed", n)
					}
				}
			}()

			userTags := users.Group("/:user_id/tags")
			{
				userTags.POST("/", tagHandler.CreateTag)      // POST   /api/v1/users/:user_id/tags
				userTags.GET("/", tagHandler.GetTagsByUser)   // GET    /api/v1/users/:user_id/tags
				userTags.DELETE("/:id", tagHandler.DeleteTag) // DELETE /api/v1/users/:user_id/tags/:id
			}
		}
	}

	// --- Start server ---
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
