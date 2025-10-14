package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/danieldzansi/auth-api/internal/database"
	"github.com/danieldzansi/auth-api/internal/handlers"
	models "github.com/danieldzansi/auth-api/internal/models"
	"github.com/danieldzansi/auth-api/internal/repository"
	"github.com/danieldzansi/auth-api/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := database.ConnectGorm()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Todo{}, &models.Tag{}, &models.TodoTag{}); err != nil {
		log.Fatal("failed to auto-migrate database schema:", err)
	}
	todoRepo := repository.NewTodoRepository(db)
	todoSvc := service.NewTodoService(todoRepo)
	tagRepo := repository.NewTagRepository(db)
	tagSvc := service.NewTagService(tagRepo)
	tagHandler := handlers.NewTagHandler(tagSvc)
	todoHandler := handlers.NewTodoHandler(todoSvc, tagSvc)

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"}, // your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Todo API is running",
			})
		})

		users := api.Group("/users")
		userlogin := api.Group("/userlogin")
		{

			users.POST("/", userHandler.CreateUser)
			userlogin.POST("/", userHandler.Userlogin)
			auth := users.Group("/")
			auth.Use(handlers.AuthMiddleware())

			userTodos := users.Group("/todos")
			userTodos.Use(handlers.AuthMiddleware())
			{
				userTodos.POST("/", todoHandler.CreateTodo)             // POST   /api/v1/users/:user_id/todos
				userTodos.GET("/", todoHandler.GetAllTodos)             // GET    /api/v1/users/todos
				userTodos.GET("/:id", todoHandler.GetTodoByID)          // GET    /api/v1/users/todos/:id
				userTodos.PATCH("/:id", todoHandler.UpdateTodoHandler)  // PATCH  /api/v1/users/todos/:id
				userTodos.DELETE("/:id", todoHandler.DeleteTodoHandler) // DELETE /api/v1/users/todos/:id
			}

			intervalStr := os.Getenv("DUE_CHECK_INTERVAL")
			if intervalStr == "" {
				intervalStr = "1m"
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
			userTags.Use(handlers.AuthMiddleware())
			{
				userTags.POST("/", tagHandler.CreateTag)      // POST   /api/v1/users/:user_id/tags
				userTags.GET("/", tagHandler.GetTagsByUser)   // GET    /api/v1/users/:user_id/tags
				userTags.DELETE("/:id", tagHandler.DeleteTag) // DELETE /api/v1/users/:user_id/tags/:id
			}
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
