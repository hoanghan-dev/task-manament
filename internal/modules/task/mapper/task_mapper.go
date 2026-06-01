package mapper

import (
	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/dto/response"
	"dev/task-management/internal/modules/task/entities"
	"time"

	"github.com/google/uuid"
)

func TaskRequestToEntity(id uuid.UUID, dto *request.CreateTaskRequest, createAt time.Time) *entities.Task {
	return entities.NewTask(id, dto.Title, dto.Description, dto.Status, dto.Assignee, dto.Workspace, createAt)
}

func TaskUpdateToEntity(t *request.UpdateTaskRequest) *entities.Task {
	return &entities.Task{
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Workspace:   t.Workspace,
	}
}

func EntityToTaskResponse(t *entities.Task) *response.TaskResponse {
	return &response.TaskResponse{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Assignee:    t.Assignee,
		Workspace:   t.Workspace,
		CreateAt:    t.CreateAt,
	}
}
