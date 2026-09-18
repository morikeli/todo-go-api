package repository

import (
	"context"
	"fmt"
	"time"
	"todo-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTask(pool *pgxpool.Pool, title string, description *string, completed bool) (*models.Task, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5000*time.Millisecond)

	defer cancel()

	query := `
		INSERT INTO todos (title, description, completed)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, completed, created_at, updated_at
	`
	var task models.Task
	err := pool.QueryRow(ctx, query, title, description, completed).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil

}
