package service

import (
	"context"
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/logger"
)

type loginService struct {
	userRepository      models.UserRepository
	refreshTokenService models.RefreshTokenService
	contextTimeout      time.Duration
}

// GetUserByUsername implements models.LoginService.
func (l *loginService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	logger.Logger.Infof("Запрос юзера по юзернейму: %s", username)
	user, err := l.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка получения юзера по юзернейму")
		return nil, err
	}
	return user, nil
}

// CheckUsernameExists implements models.LoginService.
func (l *loginService) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	c, cancel := context.WithTimeout(ctx, time.Duration(l.contextTimeout))
	defer cancel()
	exists, err := l.userRepository.CheckUsernameExists(c, username)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка проверки существования юзернейма")
		return false, err
	}
	logger.Logger.Infof("Проверка существования юзернейма %s: %v", username, exists)
	return exists, nil
}

// CreateAccessToken implements models.LoginService.
func (l *loginService) CreateAccessToken(user *models.User, expired int) (accessToken string, err error) {
	logger.Logger.Infof("Создание акцесс токена: %s", user.ID)
	return l.refreshTokenService.CreateAccessToken(user, expired)
}

// CreateRefreshToken implements models.LoginService.
func (l *loginService) CreateRefreshToken(user *models.User, expired int) (refreshToken string, err error) {
	logger.Logger.Infof("Создание рефрещ токена для юзера с айди: %s", user.ID)
	return l.refreshTokenService.CreateRefreshToken(user, expired)
}

// RevokeTokens implements models.LoginService.
func (l *loginService) RevokeTokens(ctx context.Context, user *models.User) error {
	logger.Logger.Infof("обнуление токеноав : %s", user.ID)
	err := l.userRepository.DeleteUserRefreshTokenByEmail(ctx, user.Email)
	if err != nil {
		logger.Logger.WithError(err).Error("oшибка")
		return err
	}
	return nil
}

func NewLoginService(userRepository models.UserRepository, refreshTokenService models.RefreshTokenService, timeout time.Duration) models.LoginService {
	return &loginService{
		userRepository:      userRepository,
		refreshTokenService: refreshTokenService,
		contextTimeout:      timeout,
	}
}
