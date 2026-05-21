package repositories

import (
	"context"
	"database/sql"
	"dev/task-management/internal/modules/auth/entities"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
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
	sql := `select user_id, email, password_hash, full_name, create_at
			from users where email = $1`

	row := r.database.QueryRowContext(ctx, sql, email)

	var user entities.User

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.CreateAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, row.Err()
}

func (r *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
	sql := `insert into users (user_id, email, password_hash, full_name, create_at)
			values ($1, $2, $3, $4, $5)`

	_, err := r.database.ExecContext(ctx, sql, user.Id, user.Email, user.Password, user.FullName, user.CreateAt)
	return err
}
