package validates

import (
	"context"
	"dev/task-management/internal/modules/auth/repositories"
	"dev/task-management/pkg/apperror"
	"errors"
	"log"
)

var (
	MIN int = 6
	MAX int = 20
)

func PasswordIsValid(pass string) error {
	if pass == "" {
		return apperror.NewValidation("password must be not null")
	}

	if len(pass) < MIN || len(pass) > MAX {
		return apperror.NewValidation("password must be between 6 and 20 characters long")
	}

	return nil
}

func EmailIsValid(email string, ctx context.Context, repo repositories.UserRepository) error {
	user, err := repo.GetUserByEmail(ctx, email)

	if err != nil {
		// If user is not found, that means the email is available → valid
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.Code == apperror.CodeNotFound {
			return nil
		}
		// Real DB error → propagate as internal error
		log.Println("Check email: ", err)
		return apperror.Wrap(err, apperror.NewInternal("failed to check email availability"))
	}

	// User found → email already taken
	if user != nil {
		return apperror.NewConflict("email already exists")
	}

	return nil
}
