package mapper

import (
	"dev/task-management/internal/dto/request"
	"dev/task-management/internal/entities"

	"github.com/google/uuid"
)

func TaskRequestToEntity(id uuid.UUID, dto *request.TaskRequest) *entities.Task {
	return entities.NewTask(id, dto.Title, dto.Description, dto.Status, dto.Assignee)
}
