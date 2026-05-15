package db

import (
	"dev/task-management/internal/entities"

	"github.com/google/uuid"
)

var TaskDataMock = []*entities.Task{
	{
		Id:          uuid.New(),
		Title:       "Implement Login API",
		Description: "Create REST API for user authentication using JWT",
		Status:      "TODO",
		Assignee:    "Giyuu",
	},
	{
		Id:          uuid.New(),
		Title:       "Setup PostgreSQL Database",
		Description: "Initialize PostgreSQL and create migration scripts",
		Status:      "IN_PROGRESS",
		Assignee:    "Han",
	},
	{
		Id:          uuid.New(),
		Title:       "Create User Entity",
		Description: "Design user entity and validation rules",
		Status:      "DONE",
		Assignee:    "Minh",
	},
	{
		Id:          uuid.New(),
		Title:       "Build Task CRUD API",
		Description: "Implement create, update, delete, and get task endpoints",
		Status:      "TODO",
		Assignee:    "Khoa",
	},
	{
		Id:          uuid.New(),
		Title:       "Add Global Error Handler",
		Description: "Handle application errors consistently",
		Status:      "IN_PROGRESS",
		Assignee:    "Linh",
	},
	{
		Id:          uuid.New(),
		Title:       "Implement Middleware Logger",
		Description: "Log incoming requests and responses",
		Status:      "DONE",
		Assignee:    "Dat",
	},
	{
		Id:          uuid.New(),
		Title:       "Design Database Schema",
		Description: "Create ERD and normalize database tables",
		Status:      "TODO",
		Assignee:    "Phuc",
	},
	{
		Id:          uuid.New(),
		Title:       "Integrate Swagger",
		Description: "Generate API documentation using Swagger",
		Status:      "TODO",
		Assignee:    "An",
	},
	{
		Id:          uuid.New(),
		Title:       "Setup Docker Environment",
		Description: "Dockerize Go application and PostgreSQL",
		Status:      "IN_PROGRESS",
		Assignee:    "Bao",
	},
	{
		Id:          uuid.New(),
		Title:       "Implement Refresh Token",
		Description: "Add refresh token mechanism for authentication",
		Status:      "TODO",
		Assignee:    "Tuan",
	},
	{
		Id:          uuid.New(),
		Title:       "Create Pagination Feature",
		Description: "Add pagination for task listing API",
		Status:      "DONE",
		Assignee:    "Quang",
	},
	{
		Id:          uuid.New(),
		Title:       "Optimize SQL Queries",
		Description: "Improve query performance and indexing",
		Status:      "IN_PROGRESS",
		Assignee:    "Vy",
	},
	{
		Id:          uuid.New(),
		Title:       "Implement Unit Testing",
		Description: "Write unit tests for service layer",
		Status:      "TODO",
		Assignee:    "Long",
	},
	{
		Id:          uuid.New(),
		Title:       "Add Role Authorization",
		Description: "Restrict API access based on user roles",
		Status:      "TODO",
		Assignee:    "Trang",
	},
	{
		Id:          uuid.New(),
		Title:       "Create Config Loader",
		Description: "Load environment variables from .env file",
		Status:      "DONE",
		Assignee:    "Hieu",
	},
	{
		Id:          uuid.New(),
		Title:       "Implement Search API",
		Description: "Allow searching tasks by title and status",
		Status:      "IN_PROGRESS",
		Assignee:    "Nhi",
	},
	{
		Id:          uuid.New(),
		Title:       "Setup CI/CD Pipeline",
		Description: "Automate build and deployment workflow",
		Status:      "TODO",
		Assignee:    "Son",
	},
	{
		Id:          uuid.New(),
		Title:       "Refactor Project Structure",
		Description: "Apply modular monolith architecture",
		Status:      "DONE",
		Assignee:    "Duc",
	},
	{
		Id:          uuid.New(),
		Title:       "Implement File Upload",
		Description: "Allow users to upload avatar images",
		Status:      "TODO",
		Assignee:    "My",
	},
	{
		Id:          uuid.New(),
		Title:       "Add Redis Cache",
		Description: "Cache task list to improve performance",
		Status:      "IN_PROGRESS",
		Assignee:    "Kiet",
	},
}
