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

func TestExtractIDFromToken(t *testing.T) {
	mockJWTUtils := new(mocks.JWTUtils)
	testUUID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockJWTUtils.On("ExtractIDFromToken", mock.Anything).Return(testUUID, nil).Once()

		serv := service.NewRefreshTokenService(nil, mockJWTUtils, time.Second*2)
		id, err := serv.ExtractIDFromToken("someToken")

		assert.NoError(t, err)
		assert.Equal(t, testUUID, id)

		mockJWTUtils.AssertExpectations(t)
	})
}

func TestIsAuthorized(t *testing.T) {
	mockJWTUtils := new(mocks.JWTUtils)

	t.Run("authorized", func(t *testing.T) {
		mockJWTUtils.On("IsAuthorized", mock.Anything).Return(true, nil).Once()

		serv := service.NewRefreshTokenService(nil, mockJWTUtils, time.Second*2)
		isAuthorized, err := serv.IsAuthorized("someToken")

		assert.NoError(t, err)
		assert.True(t, isAuthorized)

		mockJWTUtils.AssertExpectations(t)
	})
}

func TestGetUserByEmail(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockUser := &models.User{
		Email: "vasiliy@example.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepository.On("GetUserByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil).Once()

		serv := service.NewRefreshTokenService(mockUserRepository, nil, time.Second*2)
		user, err := serv.GetUserByEmail(context.Background(), mockUser.Email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, mockUser.Email, user.Email)

		mockUserRepository.AssertExpectations(t)
	})
}

func TestExtractIPFromRefreshToken(t *testing.T) {
	mockJWTUtils := new(mocks.JWTUtils)

	t.Run("success", func(t *testing.T) {
		mockJWTUtils.On("ExtractIPFromRefreshToken", mock.Anything, mock.Anything).Return("127.0.0.1", nil).Once()

		serv := service.NewRefreshTokenService(nil, mockJWTUtils, time.Second*2)
		ip, err := serv.ExtractIPFromRefreshToken("refresh_token", "originalToken")

		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1", ip)

		mockJWTUtils.AssertExpectations(t)
	})
}

func TestExtractEmailFromRefreshToken(t *testing.T) {
	mockJWTUtils := new(mocks.JWTUtils)

	t.Run("success", func(t *testing.T) {
		mockJWTUtils.On("ExtractEmailFromRefreshToken", mock.Anything).Return("vasiliy@example.com", nil).Once()

		serv := service.NewRefreshTokenService(nil, mockJWTUtils, time.Second*2)
		email, err := serv.ExtractEmailFromRefreshToken("refresh_token")

		assert.NoError(t, err)
		assert.Equal(t, "vasiliy@example.com", email)

		mockJWTUtils.AssertExpectations(t)
	})
}



