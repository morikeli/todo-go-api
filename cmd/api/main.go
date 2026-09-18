package main

import (
	"log"
	"todo-api/internal/config"
	database "todo-api/internal/db"
	"todo-api/internal/handlers"

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

	router.POST("/task/create", handlers.CreateTaskHandler(pool))
	router.GET("/tasks/all", handlers.GetAllTasksHandler(pool))
	router.GET("/task/:id", handlers.GetTaskByIdHandler(pool))
	router.PUT("/task/:id/update", handlers.UpdateTaskHandler(pool))
	router.DELETE("/task/:id/delete", handlers.DeleteTaskHandler(pool))

	router.Run(":" + cfg.Port)

}