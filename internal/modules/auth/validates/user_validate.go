package validates

import (
	"context"
	"dev/task-management/internal/modules/auth/repositories"
	"errors"
)

var (
	MIN int = 6
	MAX int = 20
)

func PasswordIsValid(pass string) (bool, error) {
	if pass == "" {
		return false, errors.New("password must be not null.")
	}

	if len(pass) < MIN || len(pass) > MAX {
		return false, errors.New("password must be between 6 and 20 characters long.")
	}

	return true, nil
}

func EmailIsValid(email string, ctx context.Context, repo repositories.UserRepository) (bool, error) {
	user, _ := repo.GetUserByEmail(ctx, email)

	if user != nil {
		return false, errors.New("Email already exists")
	}

	return true, nil

}
