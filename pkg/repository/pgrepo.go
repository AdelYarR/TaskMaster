package repository

import (
	"TaskMaster/pkg/models"
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewRepo(s string, log *slog.Logger) *PGRepo {
	pool, err := pgxpool.Connect(context.Background(), s)
	if err != nil {
		log.Error("Failed to create repo")
	}

	return &PGRepo{
		pool:   pool,
		logger: log,
	}
}

func (p *PGRepo) SignUp(email string, password string) (int, error) {
	var userID int

	err := p.pool.QueryRow(context.Background(),
		`INSERT INTO users (email, password) 
		VALUES ($1, $2)
		RETURNING id;`,
		email,
		password,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (p *PGRepo) SignIn(email string) (int, string, error) {
	var user_id int
	var hash string
	err := p.pool.QueryRow(context.Background(),
		`SELECT id, password FROM users WHERE email = $1;`,
		email,
	).Scan(&user_id, &hash)
	if err != nil {
		return 0, "", err
	}

	return user_id, hash, nil
}

func (p *PGRepo) GetTasks(userID int) ([]models.Task, error) {
	rows, err := p.pool.Query(context.Background(),
		`SELECT title, description, status, priority, due_date, created_at FROM tasks
	WHERE user_id = $1;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task

		err := rows.Scan(&t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (p *PGRepo) AddTask(task models.Task, userID int) error {
	_, err := p.pool.Exec(context.Background(),
		`INSERT INTO tasks (title, description, status, priority, due_date, created_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		time.Now(),
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}
