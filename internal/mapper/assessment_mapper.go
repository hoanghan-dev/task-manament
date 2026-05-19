package mapper

import (
	"dev/task-management/internal/dto/request"
	"dev/task-management/internal/entities"

	"github.com/google/uuid"
)

func AssessmentRequestToEntity(id uuid.UUID, ar *request.AssessmentRequest) *entities.Assessment {
	return entities.NewAssessment(id, ar.Name)
}
