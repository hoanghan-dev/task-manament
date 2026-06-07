package services

import (
	"context"
	"errors"
	"dev/task-management/internal/modules/auth/dto/request"
	"dev/task-management/internal/modules/auth/dto/response"
	"dev/task-management/internal/modules/auth/mapper"
	"dev/task-management/internal/modules/auth/repositories"
	"dev/task-management/internal/modules/auth/validates"
	wpService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/pkg/apperror"
	"dev/task-management/pkg/utils"

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

	if err := validates.EmailIsValid(userDTO.Email, ctx, s.userRepo); err != nil {
		return nil, err // AppError: Conflict (409) or Internal (500)
	}

	if err := validates.PasswordIsValid(userDTO.Password); err != nil {
		return nil, err // AppError: Validation (400)
	}

	user := mapper.RegisterUserRequestToEntity(userDTO)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 12)

	if err != nil {
		return nil, apperror.Wrap(err, apperror.NewInternal("failed to hash password"))
	}

	user.Password = string(passwordHash)

	errCreate := s.userRepo.CreateUser(ctx, user)

	if errCreate != nil {
		return nil, errCreate // AppError from repository (409 duplicate, 500 internal)
	}

	wpRes, err := s.wpService.CreateWorkspaceDefault(ctx, user.Id, user.FullName)

	if err != nil {
		return nil, err // AppError from workspace service
	}
	return mapper.EntityToUserResponse(user, wpRes), nil
}

func (s *authService) Login(ctx context.Context, userDTO *request.LoginUserRequestDTO) (*response.UserAuthResponseDTO, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, userDTO.Email)
	if err != nil {
		// If user not found → return generic auth error (don't reveal email existence)
		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.Code == apperror.CodeNotFound {
			return nil, apperror.NewUnauthorized("email or password invalid")
		}
		// Real DB error → internal
		return nil, apperror.Wrap(err, apperror.NewInternal("authentication failed"))
	}

	if user == nil {
		return nil, apperror.NewUnauthorized("email or password invalid")
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password))

	if errPass != nil {
		return nil, apperror.NewUnauthorized("email or password invalid")
	}

	token, err := utils.GenerateAccessToken(user.Id, user.Email)

	if err != nil {
		return nil, apperror.Wrap(err, apperror.NewInternal("failed to generate access token"))
	}

	return &response.UserAuthResponseDTO{
		AccessToken:  token,
		RefreshToken: uuid.NewString(),
	}, nil
}
