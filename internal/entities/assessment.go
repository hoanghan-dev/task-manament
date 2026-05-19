package entities

import "github.com/google/uuid"

type Assessment struct {
	Id   uuid.UUID
	Name string
}

func NewAssessment(id uuid.UUID, name string) *Assessment {
	return &Assessment{
		Id:   id,
		Name: name,
	}
}
