package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Email    string
	Password string
	FullName string
	CreateAt time.Time
}

func NewUser(id uuid.UUID, email string, password string, fullName string, createAt time.Time) *User {
	return &User{
		Id:       id,
		Email:    email,
		Password: password,
		FullName: fullName,
		CreateAt: createAt,
	}
}
