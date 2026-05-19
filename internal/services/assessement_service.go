package services

import (
	"dev/task-management/internal/dto/request"
	"dev/task-management/internal/entities"
	"dev/task-management/internal/mapper"
	"dev/task-management/internal/repositories"

	"github.com/google/uuid"
)

type AssessmentService interface {
	GetAssessmentList() ([]*entities.Assessment, error)
	GetAssessment(id uuid.UUID) (*entities.Assessment, error)
	CreateAssessment(ass *request.AssessmentRequest)
	UpdateAssessment(id uuid.UUID, ass *request.AssessmentRequest) error
	DeleteAssessment(id uuid.UUID) error
}

type assessmentService struct {
	assessmentRepository repositories.AssessmentRepository
}

func NewAssessmentService(repo repositories.AssessmentRepository) AssessmentService {
	return &assessmentService{
		assessmentRepository: repo,
	}
}

func (s *assessmentService) GetAssessmentList() ([]*entities.Assessment, error) {
	assessments, err := s.assessmentRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return assessments, nil
}
func (s *assessmentService) GetAssessment(id uuid.UUID) (*entities.Assessment, error) {
	assessment, err := s.assessmentRepository.FindById(id)

	if err != nil {
		return nil, err
	}

	return assessment, nil
}
func (s *assessmentService) CreateAssessment(ass *request.AssessmentRequest) {
	id := uuid.New()
	assessment := mapper.AssessmentRequestToEntity(id, ass)
	s.assessmentRepository.Create(assessment)
}

func (s *assessmentService) UpdateAssessment(id uuid.UUID, ass *request.AssessmentRequest) error {
	assessment := mapper.AssessmentRequestToEntity(id, ass)
	err := s.assessmentRepository.Update(id, assessment)

	if err != nil {
		return err
	}
	return nil
}
func (s *assessmentService) DeleteAssessment(id uuid.UUID) error {
	err := s.assessmentRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
