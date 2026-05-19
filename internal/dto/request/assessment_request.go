package request

type AssessmentRequest struct {
	Name string `json:"name" binding:"required"`
}
