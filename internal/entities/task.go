package entities

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewTask(id int, title string, description string) *Task {
	return &Task{
		Id:          id,
		Title:       title,
		Description: description,
	}
}
