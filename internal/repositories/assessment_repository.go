package repositories

import (
	"dev/task-management/internal/db"
	"dev/task-management/internal/entities"
	"errors"

	"github.com/google/uuid"
)

type AssessmentRepository interface {
	FindAll() ([]*entities.Assessment, error)
	FindById(id uuid.UUID) (*entities.Assessment, error)
	Create(ass *entities.Assessment)
	Update(id uuid.UUID, ass *entities.Assessment) error
	Delete(id uuid.UUID) error
}

type assessmentRepository struct {
}

func NewAssessmentRepository() AssessmentRepository {
	return &assessmentRepository{}
}

func (r *assessmentRepository) FindAll() ([]*entities.Assessment, error) {
	assessments := db.AssessmentMockData
	if len(assessments) == 0 {
		return nil, errors.New("task list is empty.")
	}
	return assessments, nil
}

func (r *assessmentRepository) FindById(id uuid.UUID) (*entities.Assessment, error) {
	for _, ass := range db.AssessmentMockData {
		if ass.Id == id {
			return ass, nil
		}
	}
	return nil, errors.New("task not found")
}

func (r *assessmentRepository) Create(ass *entities.Assessment) {
	db.AssessmentMockData = append(db.AssessmentMockData, ass)
}

func (r *assessmentRepository) Update(id uuid.UUID, ass *entities.Assessment) error {
	for i, assessment := range db.AssessmentMockData {
		if assessment.Id == id {
			db.AssessmentMockData[i] = ass
			return nil
		}
	}
	return errors.New("task not found")
}

func (r *assessmentRepository) Delete(id uuid.UUID) error {
	for i, assessment := range db.AssessmentMockData {
		if assessment.Id == id {
			db.AssessmentMockData = append(db.AssessmentMockData[:i], db.AssessmentMockData[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
