package service

import (
	"context"
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/logger"
	"github.com/Seven11Eleven/auth_service_medods/internal/utils"
	"github.com/google/uuid"
)

type refreshTokenService struct {
	userRepository models.UserRepository
	jwtUtils       utils.JWTUtils
	contextTimeout time.Duration
}

// ExtractIDFromToken implements models.RefreshTokenService.
func (r *refreshTokenService) ExtractIDFromToken(requestedToken string) (uuid.UUID, error) {
	logger.Logger.Debugf("Извлечение ID из токена")
	return r.jwtUtils.ExtractIDFromToken(requestedToken)
}

// IsAuthorized implements models.RefreshTokenService.
func (r *refreshTokenService) IsAuthorized(token string) (bool, error) {
	logger.Logger.Debugf("Проверка авторизации по токену")
	return r.jwtUtils.IsAuthorized(token)
}

// GetUserByEmail implements models.RefreshTokenService.
func (r *refreshTokenService) GetUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	logger.Logger.Infof("Запрос пользователя по email: %s", email)
	return r.userRepository.GetUserByEmail(ctx, email)
}

// ExtractIPFromRefreshToken implements models.RefreshTokenService.
func (r *refreshTokenService) ExtractIPFromRefreshToken(token string, originalToken string) (string, error) {
	logger.Logger.Debugf("Извлечение IP адреса из рефреш токена")
	return r.jwtUtils.ExtractIPFromRefreshToken(token, originalToken)
}

func (r *refreshTokenService) ExtractEmailFromRefreshToken(originalToken string) (string, error) {
	logger.Logger.Debugf("Извлечение email из рефреш токена")
	return r.jwtUtils.ExtractEmailFromRefreshToken(originalToken)
}

// CreateAccessToken implements models.RefreshTokenService.
func (r *refreshTokenService) CreateAccessToken(user *models.User, expired int) (accessToken string, err error) {
	logger.Logger.Infof("Создание акцесс токена для юзера с ID: %s", user.ID)
	return r.jwtUtils.CreateAccessToken(user, expired)
}

// CreateRefreshToken implements models.RefreshTokenService.
func (r *refreshTokenService) CreateRefreshToken(user *models.User, expired int) (refreshToken string, err error) {
	logger.Logger.Infof("Создание рефреш токена для юзера с ID: %s", user.ID)
	return r.jwtUtils.CreateRefreshToken(user, expired)
}

func NewRefreshTokenService(userRepository models.UserRepository, jwtUtils utils.JWTUtils, timeout time.Duration) models.RefreshTokenService {
	return &refreshTokenService{
		userRepository: userRepository,
		jwtUtils:       jwtUtils,
		contextTimeout: timeout,
	}
}
