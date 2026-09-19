package main

import (
	"log"
	"todo-api/internal/config"
	database "todo-api/internal/db"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	var cfg *config.Config
	var err error

	cfg, err = config.Load()

	if err != nil {
		log.Fatal("[ERROR]: Unable to load config: ", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("[ERROR]: Failed to connect to database: ", err)

	}

	defer pool.Close()
	
	var router *gin.Engine = gin.Default()
	
	router.SetTrustedProxies(nil)
	
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H {
			"message": "Looks like we're live and running!!",
			"status": "success",
			"database": "connected",
		})
	})

	router.POST("/auth/signup", handlers.SignupHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protectedRouters := router.Group("/task")
	protectedRouters.Use(middleware.AuthMiddleware(cfg))
	protectedRouters.POST("", handlers.CreateTaskHandler(pool))
	protectedRouters.GET("/all", handlers.GetAllTasksHandler(pool))
	protectedRouters.GET("/:id", handlers.GetTaskByIdHandler(pool))
	protectedRouters.PUT("/:id", handlers.UpdateTaskHandler(pool))
	protectedRouters.DELETE("/:id", handlers.DeleteTaskHandler(pool))


	// Test route
	router.GET("/test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())

	router.Run(":" + cfg.Port)

}