package services_test

// import (
// 	"dev/task-management/internal/modules/task/dto/request"
// 	"dev/task-management/internal/modules/task/entities"
// 	"dev/task-management/internal/modules/task/repositories"
// 	"dev/task-management/internal/modules/task/services"
// 	"testing"

// 	"github.com/google/uuid"
// )

// func TestGetAllTask(t *testing.T) {

// 	repo := repositories.NewTaskRepository()
// 	service := services.NewTaskService(repo)

// 	tasks, err := service.GetAllTask()

// 	if len(tasks) != 20 {
// 		t.Errorf("expected 20, got %d", len(tasks))
// 	}

// 	if err != nil {
// 		t.Errorf("expected error = nil, got error = %s", err.Error())
// 	}

// }

// func TestGetTask(t *testing.T) {
// 	id := uuid.MustParse("d010a8fb-b33d-4311-999b-9150160d994b")

// 	repo := repositories.NewTaskRepository()
// 	service := services.NewTaskService(repo)

// 	task, err := service.GetTask(id)

// 	if err != nil {
// 		t.Errorf("expected error = nil, got error = %s", err.Error())
// 	}

// 	if task.Title != "Implement Login API" {
// 		t.Errorf("expected title = Implement Login API, got etititletlerror = %s", task.Title)
// 	}

// }

// func TestCreateTask(t *testing.T) {

// 	tests := []struct {
// 		name        string
// 		taskRequest request.TaskRequest
// 		expected    entities.Task
// 	}{
// 		{
// 			name: "Create task successfully",
// 			taskRequest: request.TaskRequest{
// 				Title:       "Integrate Swagger",
// 				Description: "Generate API documentation using Swagger cdnnnn",
// 				Status:      "TODO",
// 				Assignee:    "An",
// 			},
// 			expected: entities.Task{
// 				Title:       "Integrate Swagger",
// 				Description: "Generate API documentation using Swagger cdnnnn",
// 				Status:      "TODO",
// 				Assignee:    "An",
// 			},
// 		},
// 	}

// 	repo := repositories.NewTaskRepository()
// 	service := services.NewTaskService(repo)

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			result, err := service.CreateTask(&test.taskRequest)
// 			if err != nil {
// 				t.Errorf("expected error = nil, got error = %s", err.Error())
// 			}

// 			if result.Title != test.expected.Title {
// 				t.Errorf("expected title = %s, got title = %s", test.expected.Title, result.Title)
// 			}

// 			if result.Description != test.expected.Description {
// 				t.Errorf("expected description = %s, got description = %s", test.expected.Description, result.Description)
// 			}

// 			if result.Status != test.expected.Status {
// 				t.Errorf("expected status = %s, got status = %s", test.expected.Status, result.Status)
// 			}

// 			if result.Assignee != test.expected.Assignee {
// 				t.Errorf("expected assignee = %s, got assignee = %s", test.expected.Assignee, result.Assignee)
// 			}

// 			if result.Id == uuid.Nil {
// 				t.Error("expected id is not nil")
// 			}
// 		})
// 	}
// }

// func TestUpdateTask(t *testing.T) {

// 	tests := []struct {
// 		name        string
// 		Id          uuid.UUID
// 		taskRequest request.TaskRequest
// 		expected    error
// 	}{
// 		{
// 			name: "Update task successfully",
// 			Id:   uuid.MustParse("e2e95497-628a-4b63-b945-619deb7ad2db"),
// 			taskRequest: request.TaskRequest{
// 				Title:       "Integrated Swagger!!",
// 				Description: "Generate API documentation using Swagger",
// 				Status:      "TODO",
// 				Assignee:    "An",
// 			},
// 			expected: nil,
// 		},
// 	}

// 	repo := repositories.NewTaskRepository()
// 	service := services.NewTaskService(repo)

// 	for _, test := range tests {
// 		t.Run(
// 			test.name, func(t *testing.T) {
// 				err := service.UpdateTask(test.Id, &test.taskRequest)
// 				if err != test.expected {
// 					t.Errorf("expected error = %v, got error = %v", test.expected, err)
// 				}
// 			},
// 		)
// 	}

// }
