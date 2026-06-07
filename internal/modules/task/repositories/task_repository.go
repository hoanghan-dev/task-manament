package repositories

import (
	"context"
	"database/sql"
	"dev/task-management/internal/modules/task/entities"
	"dev/task-management/pkg/apperror"

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
	UpdateTaskStatus(ctx context.Context, taskId uuid.UUID, status string) error
	TaskExistsInWorkspace(ctx context.Context, taskId uuid.UUID, workspaceId uuid.UUID) bool
	CanUserAccessTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (bool, error)
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

	sqlStr := `select  t.task_id, t.title, t.description, t.status, t.assignee_id, t.workspace_id, t.create_at from tasks t
			join workspaces wp on t.workspace_id = wp.workspace_id
			where wp.owner_id = $1 or t.assignee_id = $2
			order by t.create_at desc`

	rows, err := r.database.QueryContext(ctx, sqlStr, ownerId, ownerId)

	if err != nil {
		return nil, apperror.WrapDBError(err, "task")
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
			return nil, apperror.WrapDBError(err, "task")
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *taskRepository) FindById(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) (*entities.Task, error) {

	sqlStr := `select t.task_id, t.title, t.description, t.status, t.assignee_id, t.workspace_id, t.create_at from tasks t
			join workspaces wp on t.workspace_id = wp.workspace_id
			where (t.task_id = $1) and (wp.owner_id = $2 or t.assignee_id = $3)`

	row := r.database.QueryRowContext(ctx, sqlStr, id, ownerId, ownerId)

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
		return nil, apperror.WrapDBError(err, "task")
	}

	return &task, nil
}

func (r *taskRepository) CreateTask(ctx context.Context, t *entities.Task) error {
	sqlStr := `insert into tasks(task_id, title, description, status, assignee_id, workspace_id, create_at)
			values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.database.ExecContext(ctx, sqlStr,
		t.Id,
		t.Title,
		t.Description,
		t.Status,
		t.Assignee,
		t.Workspace,
		t.CreateAt,
	)

	if err != nil {
		return apperror.WrapDBError(err, "task")
	}
	return nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, id uuid.UUID, task *entities.Task, ownerId uuid.UUID) error {
	sqlStr := `update tasks 
        set title = $1, description = $2, status = $3, workspace_id = $4
        where task_id = $5 and workspace_id in (
            select workspace_id from workspaces where owner_id = $6
        )`
	result, err := r.database.ExecContext(ctx, sqlStr,
		task.Title,
		task.Description,
		task.Status,
		task.Workspace,
		id,
		ownerId)

	if err != nil {
		return apperror.WrapDBError(err, "task")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("database error"))
	}
	if rowsAffected == 0 {
		return apperror.NewNotFound("task")
	}
	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	sqlStr := `delete from tasks t
			where t.task_id = $1 and exists (
				select 1 from workspaces w 
				where w.workspace_id = t.workspace_id and w.owner_id = $2
			)`

	result, err := r.database.ExecContext(ctx, sqlStr, id, ownerId)

	if err != nil {
		return apperror.WrapDBError(err, "task")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("database error"))
	}
	if rowsAffected == 0 {
		return apperror.NewNotFound("task")
	}
	return nil
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

func (r *taskRepository) TaskExistsInWorkspace(ctx context.Context, taskId uuid.UUID, workspaceId uuid.UUID) bool {
	sqlQuery := `select exists (select 1 from tasks where task_id = $1 and workspace_id = $2)`
	var exists bool
	result := r.database.QueryRowContext(ctx, sqlQuery, taskId, workspaceId)

	err := result.Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}

func (r *taskRepository) AssignTask(ctx context.Context, taskId uuid.UUID, assigneeId uuid.UUID) error {
	sqlQuery := `update tasks set assignee_id = $1 where task_id = $2 and assignee_id != $3`

	result, err := r.database.ExecContext(ctx, sqlQuery, assigneeId, taskId, assigneeId)

	if err != nil {
		return apperror.WrapDBError(err, "task")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("database error"))
	}

	if rowsAffected == 0 {
		return apperror.NewConflict("task is already assigned to this user or task not found")
	}
	return nil
}

func (r *taskRepository) UpdateTaskStatus(ctx context.Context, taskId uuid.UUID, status string) error {
	sqlQuery := `update tasks set status = $1 where task_id = $2`

	result, err := r.database.ExecContext(ctx, sqlQuery, status, taskId)

	if err != nil {
		return apperror.WrapDBError(err, "task")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("database error"))
	}

	if rowsAffected == 0 {
		return apperror.NewNotFound("task")
	}
	return nil
}

// CanUserAccessTask checks if the user is the workspace owner or the task assignee
func (r *taskRepository) CanUserAccessTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (bool, error) {
	sqlQuery := `SELECT EXISTS (
		SELECT 1 FROM tasks t
		JOIN workspaces w ON t.workspace_id = w.workspace_id
		WHERE t.task_id = $1
		AND (w.owner_id = $2 OR t.assignee_id = $3)
	)`

	var canAccess bool
	err := r.database.QueryRowContext(ctx, sqlQuery, taskId, userId, userId).Scan(&canAccess)
	if err != nil {
		return false, apperror.WrapDBError(err, "task")
	}

	return canAccess, nil
}
