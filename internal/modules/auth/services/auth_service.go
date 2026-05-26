package services

import (
	"context"
	"dev/task-management/internal/modules/auth/dto/request"
	"dev/task-management/internal/modules/auth/dto/response"
	"dev/task-management/internal/modules/auth/mapper"
	"dev/task-management/internal/modules/auth/repositories"
	"dev/task-management/internal/modules/auth/validates"
	wpService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/pkg/utils"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, userDTO *request.RegisterUserRequestDTO) (*response.UserResponseDTO, error)
	Login(ctx context.Context, userDTO *request.LoginUserRequestDTO) (*response.UserAuthResponseDTO, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	wpService wpService.WorkspaceService
}

func NewAuthService(repo repositories.UserRepository, wpService wpService.WorkspaceService) AuthService {
	return &authService{
		userRepo:  repo,
		wpService: wpService,
	}
}

func (s *authService) Register(ctx context.Context, userDTO *request.RegisterUserRequestDTO) (*response.UserResponseDTO, error) {

	emailValid, err := validates.EmailIsValid(userDTO.Email, ctx, s.userRepo)

	if !emailValid {
		return nil, err
	}

	passValid, err := validates.PasswordIsValid(userDTO.Password)

	if !passValid {
		return nil, err
	}
	user := mapper.RegisterUserRequestToEntity(userDTO)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 12)

	if err != nil {
		return nil, err
	}

	user.Password = string(passwordHash)

	errCreate := s.userRepo.CreateUser(ctx, user)

	if errCreate != nil {
		return nil, errCreate
	}

	wpRes, err := s.wpService.CreateWorkspaceDefault(ctx, user.Id, user.FullName)

	if err != nil {
		return nil, err
	}
	return mapper.EntityToUserResponse(user, wpRes), nil
}

func (s *authService) Login(ctx context.Context, userDTO *request.LoginUserRequestDTO) (*response.UserAuthResponseDTO, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, userDTO.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("User not found")
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password))

	if errPass != nil {
		return nil, errors.New("email or password invalid")
	}

	token, err := utils.GenerateAccessToken(user.Id, user.Email)

	if err != nil {
		return nil, err
	}

	return &response.UserAuthResponseDTO{
		AccessToken:  token,
		RefreshToken: uuid.NewString(),
	}, nil
}
