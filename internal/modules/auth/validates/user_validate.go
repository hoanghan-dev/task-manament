package validates

import "errors"

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
