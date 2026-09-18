package handlers

import (
	"net/http"
	"strconv"
	repo "todo-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Completed   bool    `json:"is_completed"`
}

func CreateTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateTaskRequest

		// Check if the request body is valid JSON and bind it to the input struct
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task, err := repo.CreateTask(pool, input.Title, input.Description, input.Completed)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, task)

	}
}
