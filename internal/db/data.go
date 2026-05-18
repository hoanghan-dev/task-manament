package db

import (
	"dev/task-management/internal/entities"

	"github.com/google/uuid"
)

var TaskDataMock = []*entities.Task{
	{
		Id:          uuid.MustParse("d010a8fb-b33d-4311-999b-9150160d994b"),
		Title:       "Implement Login API",
		Description: "Create REST API for user authentication using JWT",
		Status:      "TODO",
		Assignee:    "Giyuu",
	},
	{
		Id:          uuid.MustParse("deb7f6bb-71ee-47d3-9d35-0c79e11b9620"),
		Title:       "Setup PostgreSQL Database",
		Description: "Initialize PostgreSQL and create migration scripts",
		Status:      "IN_PROGRESS",
		Assignee:    "Han",
	},
	{
		Id:          uuid.MustParse("f4042b64-76e4-4c7d-be65-7629bb37eb97"),
		Title:       "Create User Entity",
		Description: "Design user entity and validation rules",
		Status:      "DONE",
		Assignee:    "Minh",
	},
	{
		Id:          uuid.MustParse("389255cd-aca0-49c3-aa3d-2d6fead2e1be"),
		Title:       "Build Task CRUD API",
		Description: "Implement create, update, delete, and get task endpoints",
		Status:      "TODO",
		Assignee:    "Khoa",
	},
	{
		Id:          uuid.MustParse("54d03e2b-02e6-4d61-a359-c32d683012d2"),
		Title:       "Add Global Error Handler",
		Description: "Handle application errors consistently",
		Status:      "IN_PROGRESS",
		Assignee:    "Linh",
	},
	{
		Id:          uuid.MustParse("c197e936-4ccf-4794-86be-74972c63264a"),
		Title:       "Implement Middleware Logger",
		Description: "Log incoming requests and responses",
		Status:      "DONE",
		Assignee:    "Dat",
	},
	{
		Id:          uuid.MustParse("169f78fe-8d99-4631-9daa-3bbb0352ba98"),
		Title:       "Design Database Schema",
		Description: "Create ERD and normalize database tables",
		Status:      "TODO",
		Assignee:    "Phuc",
	},
	{
		Id:          uuid.MustParse("e2e95497-628a-4b63-b945-619deb7ad2db"),
		Title:       "Integrate Swagger",
		Description: "Generate API documentation using Swagger",
		Status:      "TODO",
		Assignee:    "An",
	},
	{
		Id:          uuid.MustParse("fca3ccc8-7f13-42ce-b6e4-a1e3024d748e"),
		Title:       "Setup Docker Environment",
		Description: "Dockerize Go application and PostgreSQL",
		Status:      "IN_PROGRESS",
		Assignee:    "Bao",
	},
	{
		Id:          uuid.MustParse("7fc0b89c-a0c0-49b1-9b06-01d07e062189"),
		Title:       "Implement Refresh Token",
		Description: "Add refresh token mechanism for authentication",
		Status:      "TODO",
		Assignee:    "Tuan",
	},
	{
		Id:          uuid.MustParse("f1e1098a-a5c5-4c80-85a2-6cf61b627a9b"),
		Title:       "Create Pagination Feature",
		Description: "Add pagination for task listing API",
		Status:      "DONE",
		Assignee:    "Quang",
	},
	{
		Id:          uuid.MustParse("c8050e57-228d-431d-8029-b373aa3293b6"),
		Title:       "Optimize SQL Queries",
		Description: "Improve query performance and indexing",
		Status:      "IN_PROGRESS",
		Assignee:    "Vy",
	},
	{
		Id:          uuid.MustParse("ff5d4e10-b64b-46ff-b1ae-bc5616468d40"),
		Title:       "Implement Unit Testing",
		Description: "Write unit tests for service layer",
		Status:      "TODO",
		Assignee:    "Long",
	},
	{
		Id:          uuid.MustParse("20dfac47-e14e-47c5-97cc-1d13a9cacc06"),
		Title:       "Add Role Authorization",
		Description: "Restrict API access based on user roles",
		Status:      "TODO",
		Assignee:    "Trang",
	},
	{
		Id:          uuid.MustParse("62140cc1-4cf2-457c-84d8-033bd08b91b5"),
		Title:       "Create Config Loader",
		Description: "Load environment variables from .env file",
		Status:      "DONE",
		Assignee:    "Hieu",
	},
	{
		Id:          uuid.MustParse("85b74e4f-517f-4a8b-b51f-901dfc298903"),
		Title:       "Implement Search API",
		Description: "Allow searching tasks by title and status",
		Status:      "IN_PROGRESS",
		Assignee:    "Nhi",
	},
	{
		Id:          uuid.MustParse("62bded06-c702-4e96-85ba-9d1a6f4bf86b"),
		Title:       "Setup CI/CD Pipeline",
		Description: "Automate build and deployment workflow",
		Status:      "TODO",
		Assignee:    "Son",
	},
	{
		Id:          uuid.MustParse("ffccfe17-0087-4977-acc2-97154ea0d925"),
		Title:       "Refactor Project Structure",
		Description: "Apply modular monolith architecture",
		Status:      "DONE",
		Assignee:    "Duc",
	},
	{
		Id:          uuid.MustParse("72dd41fc-a57c-493f-aea8-b9239b941373"),
		Title:       "Implement File Upload",
		Description: "Allow users to upload avatar images",
		Status:      "TODO",
		Assignee:    "My",
	},
	{
		Id:          uuid.MustParse("0e13897e-0e04-4b83-a96b-25ba6a680de5"),
		Title:       "Add Redis Cache",
		Description: "Cache task list to improve performance",
		Status:      "IN_PROGRESS",
		Assignee:    "Kiet",
	},
}

var AssessmentMockData = []*entities.Assessment{
	{
		Id:   uuid.MustParse("61af27d6-6f4c-462b-9d15-6bf4998b0795"),
		Name: "Java Core Assessment",
	},
	{
		Id:   uuid.MustParse("e807089c-9946-4e01-b5c6-68da298e0132"),
		Name: "Spring Boot Fundamentals",
	},
	{
		Id:   uuid.MustParse("2f8d2bba-501b-440a-9763-c47372819821"),
		Name: "REST API Design Test",
	},
	{
		Id:   uuid.MustParse("55ba8abd-4a5c-451a-8ec3-cdc605843733"),
		Name: "Database Modeling Quiz",
	},
	{
		Id:   uuid.MustParse("da9a3234-b6aa-4283-a083-6a8702736912"),
		Name: "PostgreSQL Practice Test",
	},
	{
		Id:   uuid.MustParse("d21c097b-f9f7-4b0d-92da-cb7871569164"),
		Name: "Docker Basics Assessment",
	},
	{
		Id:   uuid.MustParse("8d6e4497-7a73-4273-878d-0d77472d4985"),
		Name: "Git & GitHub Evaluation",
	},
	{
		Id:   uuid.MustParse("ff4aaeb8-edd6-43c3-b34e-83d853af8270"),
		Name: "Clean Architecture Review",
	},
	{
		Id:   uuid.MustParse("04d08879-36f9-4823-9b2a-95bafb8aed7d"),
		Name: "Microservices Knowledge Test",
	},
	{
		Id:   uuid.MustParse("3a76ce83-5dd0-4238-b7f4-503f352b0566"),
		Name: "JWT Authentication Quiz",
	},
	{
		Id:   uuid.MustParse("7d762fa5-302c-4676-9926-1694451f363b"),
		Name: "Redis Caching Assessment",
	},
	{
		Id:   uuid.MustParse("a134897b-caaf-4365-94c5-00862192ad3c"),
		Name: "Linux Command Line Test",
	},
	{
		Id:   uuid.MustParse("b948640d-7126-4890-a922-14be4568097a"),
		Name: "Go Language Syntax Quiz",
	},
	{
		Id:   uuid.MustParse("4bacf32b-9f7a-4235-a8ed-3cd7eb5770cb"),
		Name: "Gin Framework Assessment",
	},
	{
		Id:   uuid.MustParse("85d0efec-fa0c-408b-bf91-4fb1da98ccfd"),
		Name: "Concurrency in Go Test",
	},
	{
		Id:   uuid.MustParse("919d5c04-3b41-433e-8ee2-0c8d12755263"),
		Name: "Unit Testing Practice",
	},
	{
		Id:   uuid.MustParse("4d0f3f09-22b6-4f8b-85b6-6c2adf449019"),
		Name: "System Design Fundamentals",
	},
	{
		Id:   uuid.MustParse("ca3e4e38-267a-4a9a-bb17-8aec0f4297d4"),
		Name: "OOP Concepts Evaluation",
	},
	{
		Id:   uuid.MustParse("b36fb2b2-c138-462b-83cd-6adf23d3adac"),
		Name: "Data Structures Assessment",
	},
	{
		Id:   uuid.MustParse("0f9fa15c-0e74-4950-a83d-69fcddbd84b6"),
		Name: "Algorithm Problem Solving Test",
	},
}
