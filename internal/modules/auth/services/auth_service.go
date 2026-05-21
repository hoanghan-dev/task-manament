package services

import (
	"context"
	"dev/task-management/internal/modules/auth/dto/request"
	"dev/task-management/internal/modules/auth/dto/response"
	"dev/task-management/internal/modules/auth/mapper"
	"dev/task-management/internal/modules/auth/repositories"
	"dev/task-management/internal/modules/auth/validates"
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
	userRepo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: repo,
	}
}

func (s *authService) Register(ctx context.Context, userDTO *request.RegisterUserRequestDTO) (*response.UserResponseDTO, error) {
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

	return mapper.EntityToUserResponse(user), nil
}

func (s *authService) Login(ctx context.Context, userDTO *request.LoginUserRequestDTO) (*response.UserAuthResponseDTO, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, userDTO.Email)
	if err != nil {
		return nil, err
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
