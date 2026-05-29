package repositories

import (
	"context"
	"database/sql"
	"dev/task-management/internal/modules/task/entities"
	"fmt"

	"github.com/google/uuid"
)

type TaskRepository interface {
	FindAll(ctx context.Context, ownerId uuid.UUID) ([]entities.Task, error)
	FindById(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) (*entities.Task, error)
	CreateTask(ctx context.Context, t *entities.Task) error
	UpdateTask(ctx context.Context, id uuid.UUID, task *entities.Task, ownerId uuid.UUID) error
	DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error
	TaskIsExists(ctx context.Context, taskId uuid.UUID) bool
	AssignTask(ctx context.Context, taskId uuid.UUID, assigneeId uuid.UUID) error
}

type taskRepository struct {
	database *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		database: db,
	}
}

func (r *taskRepository) FindAll(ctx context.Context, ownerId uuid.UUID) ([]entities.Task, error) {

	sql := `select  t.task_id, t.title, t.description, t.status, t.assignee_id, t.workspace_id, t.create_at from tasks t
			join workspaces wp on t.workspace_id = wp.workspace_id
			where wp.owner_id = $1 or t.assignee_id = $2
			order by t.create_at desc`

	rows, err := r.database.QueryContext(ctx, sql, ownerId, ownerId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func (r *taskRepository) FindById(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) (*entities.Task, error) {

	sql := `select t.task_id, t.title, t.description, t.status, t.assignee_id, t.workspace_id, t.create_at from tasks t
			join workspaces wp on t.workspace_id = wp.workspace_id
			where (t.task_id = $1) and (wp.owner_id = $2 or t.assignee_id = $3)`

	row := r.database.QueryRowContext(ctx, sql, id, ownerId, ownerId)

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

func (r *taskRepository) UpdateTask(ctx context.Context, id uuid.UUID, task *entities.Task, ownerId uuid.UUID) error {
	sql := `update tasks 
        set title = $1, description = $2, status = $3, assignee_id = $4, workspace_id = $5
        where task_id = $6 and workspace_id in (
            select workspace_id from workspaces where owner_id = $7
        )`
	_, err := r.database.ExecContext(ctx, sql,
		task.Title,
		task.Description,
		task.Status,
		task.Assignee,
		task.Workspace,
		id,
		ownerId)
	return err
}

func (r *taskRepository) DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	sql := `delete from tasks t
			where t.task_id = $1 and exists (
				select 1 from workspaces w 
				where w.workspace_id = t.workspace_id and w.owner_id = $2
			)`

	_, err := r.database.ExecContext(ctx, sql, id, ownerId)
	return err
}

func (r *taskRepository) TaskIsExists(ctx context.Context, taskId uuid.UUID) bool {
	sqlQuery := `select exists (select 1 from tasks where task_id = $1)`
	var exists bool
	result := r.database.QueryRowContext(ctx, sqlQuery, taskId)

	err := result.Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}

func (r *taskRepository) AssignTask(ctx context.Context, taskId uuid.UUID, assigneeId uuid.UUID) error {
	sqlQuery := `update tasks set assignee_id = $1 where task_id = $2`

	result, err := r.database.ExecContext(ctx, sqlQuery, assigneeId, taskId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failed to assign task")
	}
	return nil
}
