package repositories

import (
	"context"
	"database/sql"
	"dev/task-management/internal/modules/auth/entities"
	"errors"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	UserIsExists(ctx context.Context, userId uuid.UUID) bool
}

type userRepository struct {
	database *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		database: db,
	}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `select user_id, email, password_hash, full_name, create_at
			from users where email = $1`

	row := r.database.QueryRowContext(ctx, query, email)

	var user entities.User

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.CreateAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
	sql := `insert into users (user_id, email, password_hash, full_name, create_at)
			values ($1, $2, $3, $4, $5)`

	_, err := r.database.ExecContext(ctx, sql, user.Id, user.Email, user.Password, user.FullName, user.CreateAt)
	return err
}

func (r *userRepository) UserIsExists(ctx context.Context, userId uuid.UUID) bool {
	sqlQuery := `select exists (select 1 from users where user_id = $1)`

	var exists bool

	result := r.database.QueryRowContext(ctx, sqlQuery, userId)

	err := result.Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}
