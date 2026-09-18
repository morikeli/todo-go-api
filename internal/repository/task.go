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

func GetAllTasks(pool *pgxpool.Pool) ([]models.Task, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
		`

	rows, err := pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
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

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskById(pool *pgxpool.Pool, id int) (*models.Task, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`
	var task models.Task
	err := pool.QueryRow(ctx, query, id).Scan(
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

func UpdateTask(pool *pgxpool.Pool, id int, title *string, description *string, completed bool) (*models.Task, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	query := `
		UPDATE todos
		SET title = COALESCE($1, title),
			description = COALESCE($2, description),
			completed = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, title, description, completed, created_at, updated_at
	`
	var task models.Task
	err := pool.QueryRow(ctx, query, title, description, completed, id).Scan(
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

func DeleteTask(pool *pgxpool.Pool, id int) error {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	query := `
		DELETE FROM todos
		WHERE id = $1
	`

	cmd, err := pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("Task with id %d not found!", id)
	}

	return nil	// success
}