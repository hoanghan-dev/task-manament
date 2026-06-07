package repositories

import (
	"context"
	"database/sql"
	entities "dev/task-management/internal/modules/workspace/entities"
	"dev/task-management/pkg/apperror"

	"github.com/google/uuid"
)

type WorkspaceRepository interface {
	FindByOwnerId(ctx context.Context, ownerId uuid.UUID) (*entities.Workspace, error)
	Create(ctx context.Context, wp *entities.Workspace) (*entities.Workspace, error)
	Update(ctx context.Context, wp *entities.Workspace, wpId uuid.UUID, ownerId uuid.UUID) error
	Delete(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) error
	IsWorkspaceOwnedBy(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) (bool, error)
}

type workspaceRepository struct {
	db *sql.DB
}

func NewWorkspaceRepository(db *sql.DB) WorkspaceRepository {
	return &workspaceRepository{
		db: db,
	}
}

func (r *workspaceRepository) FindByOwnerId(ctx context.Context, ownerId uuid.UUID) (*entities.Workspace, error) {
	sqlQuery := `select workspace_id, name, description, owner_id, create_at from workspaces
			where owner_id = $1`

	// QueryContext (multi-row) never returns sql.ErrNoRows, so we check rows.Next() instead.
	rows, err := r.db.QueryContext(ctx, sqlQuery, ownerId)

	if err != nil {
		return nil, apperror.WrapDBError(err, "workspace")
	}

	defer rows.Close()

	var wp entities.Workspace
	if rows.Next() {

		err := rows.Scan(
			&wp.WorkspaceId,
			&wp.Name,
			&wp.Description,
			&wp.OwnerId,
			&wp.CreateAt,
		)

		if err != nil {
			return nil, apperror.WrapDBError(err, "workspace")
		}
	} else {
		// No rows found → workspace not found for this owner
		return nil, apperror.NewNotFound("workspace")
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapDBError(err, "workspace")
	}

	return &wp, nil
}

func (r *workspaceRepository) Create(ctx context.Context, wp *entities.Workspace) (*entities.Workspace, error) {
	sqlStr := `insert into workspaces (workspace_id,name,description,owner_id,create_at)
			values ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, sqlStr, wp.WorkspaceId, wp.Name, wp.Description, wp.OwnerId, wp.CreateAt)

	if err != nil {
		return nil, apperror.WrapDBError(err, "workspace")
	}

	return wp, nil
}

func (r *workspaceRepository) Update(ctx context.Context, wp *entities.Workspace, wpId uuid.UUID, ownerId uuid.UUID) error {
	sqlStr := `update workspaces 
			set name = $1, description = $2
			where workspace_id = $3 and owner_id = $4`
	result, err := r.db.ExecContext(ctx, sqlStr, wp.Name, wp.Description, wpId, ownerId)
	if err != nil {
		return apperror.WrapDBError(err, "workspace")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("database error"))
	}
	if rowsAffected == 0 {
		return apperror.NewNotFound("workspace")
	}
	return nil
}

func (r *workspaceRepository) Delete(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) error {
	sqlStr := `delete from workspaces where workspace_id = $1 and owner_id = $2`

	result, err := r.db.ExecContext(ctx, sqlStr, wpId, ownerId)

	if err != nil {
		return apperror.WrapDBError(err, "workspace")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("database error"))
	}
	if rowsAffected == 0 {
		return apperror.NewNotFound("workspace")
	}
	return nil
}

func (r *workspaceRepository) IsWorkspaceOwnedBy(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) (bool, error) {
	sqlStr := `select exists (
		select 1 from workspaces where workspace_id = $1 and owner_id = $2
	)`

	var ownered bool

	err := r.db.QueryRowContext(ctx, sqlStr, wpId, ownerId).Scan(&ownered)

	if err != nil {
		return false, apperror.WrapDBError(err, "workspace")
	}

	return ownered, nil
}
