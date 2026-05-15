package request

type TaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Assignee    string `json:"assignee" binding:"required"`
}
