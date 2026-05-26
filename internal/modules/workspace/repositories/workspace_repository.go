package repositories

import (
	"context"
	"database/sql"
	entities "dev/task-management/internal/modules/workspace/entities"
	"errors"

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

	rows, err := r.db.QueryContext(ctx, sqlQuery, ownerId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("workspace not found")
		}
		return nil, err
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
			return nil, err
		}
	}

	return &wp, rows.Err()
}
func (r *workspaceRepository) Create(ctx context.Context, wp *entities.Workspace) (*entities.Workspace, error) {
	sql := `insert into workspaces (workspace_id,name,description,owner_id,create_at)
			values ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, sql, wp.WorkspaceId, wp.Name, wp.Description, wp.OwnerId, wp.CreateAt)

	if err != nil {
		return nil, err
	}

	return wp, nil
}
func (r *workspaceRepository) Update(ctx context.Context, wp *entities.Workspace, wpId uuid.UUID, ownerId uuid.UUID) error {
	sql := `update workspaces 
			set name = $1, description = $2
			where workspace_id = $3 and owner_id = $4`
	_, err := r.db.ExecContext(ctx, sql, wp.Name, wp.Description, wpId, ownerId)
	if err != nil {
		return err
	}
	return nil
}
func (r *workspaceRepository) Delete(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) error {
	sql := `delete from workspaces where workspace_id = $1 and owner_id = $2`

	_, err := r.db.ExecContext(ctx, sql, wpId, ownerId)

	return err
}

func (r *workspaceRepository) IsWorkspaceOwnedBy(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) (bool, error) {
	sql := `select exists (
		select 1 from workspaces where workspace_id = $1 and owner_id = $2
	)`

	var ownered bool

	err := r.db.QueryRowContext(ctx, sql, wpId, ownerId).Scan(&ownered)

	if err != nil {
		return false, err
	}

	return ownered, nil
}
