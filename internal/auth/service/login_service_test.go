package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models/mocks"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserByUsername(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockUser := &models.User{
		Username: "Vasiliy",
		Email:    "vasiliy@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepository.On("GetUserByUsername", mock.Anything, mockUser.Username).Return(mockUser, nil).Once()

		serv := service.NewLoginService(mockUserRepository, nil, time.Second*2)
		user, err := serv.GetUserByUsername(context.Background(), mockUser.Username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, mockUser.Username, user.Username)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("not exists", func(t *testing.T) {
		mockUserRepository.On("GetUserByUsername", mock.Anything, mockUser.Username).Return(nil, models.ErrUserNotFound).Once()

		serv := service.NewLoginService(mockUserRepository, nil, time.Second*2)
		user, err := serv.GetUserByUsername(context.Background(), mockUser.Username)

		assert.Error(t, err)
		assert.Nil(t, user)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestCheckUsernameExists(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)

	t.Run("exists", func(t *testing.T) {
		mockUserRepository.On("CheckUsernameExists", mock.Anything, "Diana").Return(true, nil).Once()

		serv := service.NewLoginService(mockUserRepository, nil, time.Second*2)
		exists, err := serv.CheckUsernameExists(context.Background(), "Diana")

		assert.NoError(t, err)
		assert.True(t, exists)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("does not exist", func(t *testing.T) {
		mockUserRepository.On("CheckUsernameExists", mock.Anything, "doesnt existing").Return(false, nil).Once()

		serv := service.NewLoginService(mockUserRepository, nil, time.Second*2)
		exists, err := serv.CheckUsernameExists(context.Background(), "doesnt existing")

		assert.NoError(t, err)
		assert.False(t, exists)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestCreateAccessToken(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockRefreshTokenService := new(mocks.RefreshTokenService)
	mockUser := &models.User{
		ID: uuid.New(),
	}

	t.Run("success", func(t *testing.T) {
		mockRefreshTokenService.On("CreateAccessToken", mockUser, mock.AnythingOfType("int")).Return("access_token", nil).Once()

		serv := service.NewLoginService(mockUserRepository, mockRefreshTokenService, time.Second*2)
		token, err := serv.CreateAccessToken(mockUser, 3600)

		assert.NoError(t, err)
		assert.Equal(t, "access_token", token)

		mockRefreshTokenService.AssertExpectations(t)
	})
}

func TestCreateRefreshToken(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockRefreshTokenService := new(mocks.RefreshTokenService)
	mockUser := &models.User{
		ID: uuid.New(),
	}

	t.Run("success", func(t *testing.T) {
		mockRefreshTokenService.On("CreateRefreshToken", mockUser, mock.AnythingOfType("int")).Return("refresh_token", nil).Once()

		serv := service.NewLoginService(mockUserRepository, mockRefreshTokenService, time.Second*2)
		token, err := serv.CreateRefreshToken(mockUser, 7200)

		assert.NoError(t, err)
		assert.Equal(t, "refresh_token", token)

		mockRefreshTokenService.AssertExpectations(t)
	})
}

func TestRevokeTokens(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockUser := &models.User{
		Email: "vasiliy@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepository.On("DeleteUserRefreshTokenByEmail", mock.Anything, mockUser.Email).Return(nil).Once()

		serv := service.NewLoginService(mockUserRepository, nil, time.Second*2)
		err := serv.RevokeTokens(context.Background(), mockUser)

		assert.NoError(t, err)

		mockUserRepository.AssertExpectations(t)
	})
}
