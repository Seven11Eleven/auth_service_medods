package service

import (
	"context"
	"regexp"
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/logger"
)

type signUpService struct {
	userRepository      models.UserRepository
	refreshTokenService models.RefreshTokenService
	contextTimeout      time.Duration
}

// CheckUsernameExists implements models.SignUpService.
func (s *signUpService) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	logger.Logger.Infof("Проверка существования юзернейма: %s", username)
	c, cancel := context.WithTimeout(ctx, time.Duration(s.contextTimeout))
	defer cancel()
	exists, err := s.userRepository.CheckUsernameExists(c, username)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка проверки существования юзернейма")
		return false, err
	}
	logger.Logger.Infof("юзер : статус %s: %v", username, exists)
	return exists, nil
}

// CreateAccessToken implements models.SignUpService.
func (s *signUpService) CreateAccessToken(user *models.User, expired int) (accessToken string, err error) {
	logger.Logger.Infof("Создание акцесс токена для юзера с  ID: %s", user.ID)
	return s.refreshTokenService.CreateAccessToken(user, expired)
}

// CreateRefreshToken implements models.SignUpService.
func (s *signUpService) CreateRefreshToken(user *models.User, expired int) (refreshToken string, err error) {
	logger.Logger.Infof("Создание рефреш токена для юзера с ID: %s", user.ID)
	return s.refreshTokenService.CreateRefreshToken(user, expired)
}

// RegisterUser implements models.SignUpService.
func (s *signUpService) RegisterUser(ctx context.Context, user *models.User) (err error) {
	if !isAlpha(user.Username) {
		logger.Logger.Warnf("Неверный формат юзернейма: %s", user.Username)
		return models.ErrInvalidUsername
	}

	logger.Logger.Infof("Регистрация юзера: %s", user.Username)
	c, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	return s.userRepository.Create(c, user)
}

func NewSignUpService(userRepository models.UserRepository, refreshTokenService models.RefreshTokenService, timeout time.Duration) models.SignUpService {
	return &signUpService{
		userRepository:      userRepository,
		refreshTokenService: refreshTokenService,
		contextTimeout:      timeout,
	}
}

func isAlpha(str string) bool {
	match, _ := regexp.MatchString(`^[A-Za-z\s]+$`, str) //валидация данных
	return match
}
