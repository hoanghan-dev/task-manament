package services

import (
	"context"
	"dev/task-management/internal/modules/auth/dto/request"
	"dev/task-management/internal/modules/auth/dto/response"
	"dev/task-management/internal/modules/auth/mapper"
	"dev/task-management/internal/modules/auth/repositories"
	"dev/task-management/internal/modules/auth/validates"
	wpService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/pkg/apperror"
	"dev/task-management/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, userDTO *request.RegisterUserRequestDTO) (*response.UserResponseDTO, error)
	Login(ctx context.Context, userDTO *request.LoginUserRequestDTO) (*response.UserAuthResponseDTO, error)
	Logout(userId uuid.UUID) error
	RefreshToken(ctx context.Context, token string) (*response.UserAuthResponseDTO, error)
}

type CacheService interface {
	Get(key string, dest any) error
	Set(key string, value any, ttl time.Duration) error
	Clear(pattern string) error
}

type authService struct {
	userRepo  repositories.UserRepository
	wpService wpService.WorkspaceService
	redis     CacheService
}

func NewAuthService(repo repositories.UserRepository, wpService wpService.WorkspaceService, redis CacheService) AuthService {
	return &authService{
		userRepo:  repo,
		wpService: wpService,
		redis:     redis,
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

	token, err := utils.GenerateAccessToken(user.Id, user.Email, 5*time.Minute)

	if err != nil {
		return nil, apperror.Wrap(err, apperror.NewInternal("failed to generate access token"))
	}

	refreshToken, err := utils.GenerateAccessToken(user.Id, user.Email, 7*24*time.Hour)

	if err != nil {
		return nil, apperror.Wrap(err, apperror.NewInternal("failed to generate refresh token"))
	}

	authKey := "auth:refresh_token:user_id" + refreshToken
	err = s.redis.Set(authKey, refreshToken, 7*24*time.Hour)

	if err != nil {
		return nil, apperror.Wrap(err, apperror.NewInternal("failed to set refresh token into redis"))
	}

	return &response.UserAuthResponseDTO{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (*response.UserAuthResponseDTO, error) {

	userId, err := utils.GetUserId(token)

	if err != nil {
		return nil, apperror.NewUnauthorized("Invalid or expired refresh token")
	}

	user, err := s.userRepo.GetUserById(ctx, userId)
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

	authKey := "auth:refresh_token:user_id" + userId.String()
	var refreshToken string
	redisErr := s.redis.Get(authKey, &refreshToken)

	if redisErr != nil {
		return nil, apperror.NewUnauthorized("Invalid or expired refresh token")
	}

	if refreshToken != token {
		return nil, apperror.NewUnauthorized("Invalid or expired refresh token")
	}

	err = s.redis.Clear(authKey)

	if err != nil {
		return nil, apperror.Wrap(err, apperror.NewInternal("failed to delete refresh token in redis"))
	}

	accessToken, accessTokenErr := utils.GenerateAccessToken(userId, user.Email, 5*time.Minute)

	if accessTokenErr != nil {
		return nil, apperror.Wrap(accessTokenErr, apperror.NewInternal("failed to generate access token"))
	}

	refreshToken, refreshTokenErr := utils.GenerateAccessToken(user.Id, user.Email, 7*24*time.Hour)

	if refreshTokenErr != nil {
		return nil, apperror.Wrap(refreshTokenErr, apperror.NewInternal("failed to generate refresh token"))
	}

	err = s.redis.Set(authKey, refreshToken, 7*24*time.Hour)

	if err != nil {
		return nil, apperror.Wrap(err, apperror.NewInternal("failed to set refresh token into redis"))
	}

	return &response.UserAuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Logout(userId uuid.UUID) error {
	authKey := "auth:refresh_token:user_id:" + userId.String()

	err := s.redis.Clear(authKey)

	if err != nil {
		return apperror.Wrap(err, apperror.NewInternal("failed to delete refresh token in redis"))
	}
	return nil
}
