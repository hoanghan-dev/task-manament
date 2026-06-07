package repositories

import (
	"context"
	"database/sql"
	"dev/task-management/internal/modules/comment/entities"
	"dev/task-management/pkg/apperror"

	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entities.Comment) error
	FindByTaskId(ctx context.Context, taskId uuid.UUID) ([]entities.Comment, error)
	FindById(ctx context.Context, commentId uuid.UUID) (*entities.Comment, error)
	Delete(ctx context.Context, commentId uuid.UUID) error
	TaskExists(ctx context.Context, taskId uuid.UUID) (bool, error)
	CanUserAccessTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (bool, error)
	FindCommentReceiversByTaskId(ctx context.Context, taskId uuid.UUID) ([]uuid.UUID, error)
}

type commentRepository struct {
	database *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{
		database: db,
	}
}

func (r *commentRepository) Create(ctx context.Context, comment *entities.Comment) error {
	sqlQuery := `INSERT INTO comments (comment_id, task_id, user_id, content, create_at)
				VALUES ($1, $2, $3, $4, $5)`

	_, err := r.database.ExecContext(ctx, sqlQuery,
		comment.CommentId,
		comment.TaskId,
		comment.UserId,
		comment.Content,
		comment.CreateAt,
	)

	if err != nil {
		return apperror.WrapDBError(err, "comment")
	}
	return nil
}

func (r *commentRepository) FindByTaskId(ctx context.Context, taskId uuid.UUID) ([]entities.Comment, error) {
	sqlQuery := `SELECT comment_id, task_id, user_id, content, create_at, update_at
				FROM comments
				WHERE task_id = $1
				ORDER BY create_at ASC`

	rows, err := r.database.QueryContext(ctx, sqlQuery, taskId)
	if err != nil {
		return nil, apperror.WrapDBError(err, "comment")
	}
	defer rows.Close()

	comments := make([]entities.Comment, 0)

	for rows.Next() {
		var comment entities.Comment

		err := rows.Scan(
			&comment.CommentId,
			&comment.TaskId,
			&comment.UserId,
			&comment.Content,
			&comment.CreateAt,
			&comment.UpdateAt,
		)

		if err != nil {
			return nil, apperror.WrapDBError(err, "comment")
		}

		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *commentRepository) FindById(ctx context.Context, commentId uuid.UUID) (*entities.Comment, error) {
	sqlQuery := `SELECT comment_id, task_id, user_id, content, create_at, update_at
				FROM comments
				WHERE comment_id = $1`

	row := r.database.QueryRowContext(ctx, sqlQuery, commentId)

	var comment entities.Comment

	err := row.Scan(
		&comment.CommentId,
		&comment.TaskId,
		&comment.UserId,
		&comment.Content,
		&comment.CreateAt,
		&comment.UpdateAt,
	)

	if err != nil {
		return nil, apperror.WrapDBError(err, "comment")
	}

	return &comment, nil
}

func (r *commentRepository) Delete(ctx context.Context, commentId uuid.UUID) error {
	sqlQuery := `DELETE FROM comments WHERE comment_id = $1`

	result, err := r.database.ExecContext(ctx, sqlQuery, commentId)
	if err != nil {
		return apperror.WrapDBError(err, "comment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("database error"))
	}

	if rowsAffected == 0 {
		return apperror.NewNotFound("comment")
	}

	return nil
}

func (r *commentRepository) TaskExists(ctx context.Context, taskId uuid.UUID) (bool, error) {
	sqlQuery := `SELECT EXISTS (SELECT 1 FROM tasks WHERE task_id = $1)`

	var exists bool
	err := r.database.QueryRowContext(ctx, sqlQuery, taskId).Scan(&exists)
	if err != nil {
		return false, apperror.WrapDBError(err, "task")
	}

	return exists, nil
}

// CanUserAccessTask checks if the user is the workspace owner or the task assignee
func (r *commentRepository) CanUserAccessTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (bool, error) {
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

// FindCommentReceiversByTaskId returns deduplicated user IDs of workspace owner and task assignee
func (r *commentRepository) FindCommentReceiversByTaskId(ctx context.Context, taskId uuid.UUID) ([]uuid.UUID, error) {
	sqlQuery := `SELECT DISTINCT u.user_id FROM (
		SELECT w.owner_id AS user_id FROM tasks t
		JOIN workspaces w ON t.workspace_id = w.workspace_id
		WHERE t.task_id = $1
		UNION
		SELECT t.assignee_id AS user_id FROM tasks t
		WHERE t.task_id = $2 AND t.assignee_id IS NOT NULL
	) u`

	rows, err := r.database.QueryContext(ctx, sqlQuery, taskId, taskId)
	if err != nil {
		return nil, apperror.WrapDBError(err, "comment")
	}
	defer rows.Close()

	userIds := make([]uuid.UUID, 0)

	for rows.Next() {
		var userId uuid.UUID
		err := rows.Scan(&userId)
		if err != nil {
			return nil, apperror.WrapDBError(err, "comment")
		}
		userIds = append(userIds, userId)
	}

	return userIds, rows.Err()
}
