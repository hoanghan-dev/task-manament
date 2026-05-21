package repositories

import (
	"context"
	"database/sql"
	"dev/task-management/internal/modules/task/entities"

	"github.com/google/uuid"
)

type TaskRepository interface {
	FindAll(ctx context.Context) ([]entities.Task, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.Task, error)
	CreateTask(ctx context.Context, t *entities.Task) error
	UpdateTask(ctx context.Context, id uuid.UUID, task *entities.Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

type taskRepository struct {
	database *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		database: db,
	}
}

func (r *taskRepository) FindAll(ctx context.Context) ([]entities.Task, error) {

	sql := `select task_id, title, description, status, assignee_id, workspace_id, create_at 
			from tasks order by create_at desc`

	rows, err := r.database.QueryContext(ctx, sql)

	if err != nil {
		return nil, err
	}

	tasks := make([]entities.Task, 0)

	for rows.Next() {
		var task entities.Task

		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Assignee,
			&task.Workspace,
			&task.CreateAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *taskRepository) FindById(ctx context.Context, id uuid.UUID) (*entities.Task, error) {

	sql := `select task_id, title, description, status, assignee_id, workspace_id, create_at
			from tasks where task_id = $1`

	row := r.database.QueryRowContext(ctx, sql, id)

	var task entities.Task

	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Assignee,
		&task.Workspace,
		&task.CreateAt,
	)

	if err != nil {
		return nil, err
	}

	return &task, row.Err()
}

func (r *taskRepository) CreateTask(ctx context.Context, t *entities.Task) error {
	sql := `insert into tasks(task_id, title, description, status, assignee_id, workspace_id, create_at)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.database.ExecContext(ctx, sql,
		t.Id,
		t.Title,
		t.Description,
		t.Status,
		t.Assignee,
		t.Workspace,
		t.CreateAt,
	)
	return err
}

func (r *taskRepository) UpdateTask(ctx context.Context, id uuid.UUID, task *entities.Task) error {
	sql := `update tasks 
			set title = $1, description = $2, status = $3, assignee_id = $4, workspace_id = $5
			where task_id = $6`
	_, err := r.database.ExecContext(ctx, sql,
		task.Title,
		task.Description,
		task.Status,
		task.Assignee,
		task.Workspace,
		id)
	return err
}

func (r *taskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	sql := `delete from tasks
			where task_id = $1`

	_, err := r.database.ExecContext(ctx, sql, id)
	return err
}
