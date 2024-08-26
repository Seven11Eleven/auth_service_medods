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



func TestRegisterUser(t *testing.T) {
	mockUserRepository := new(mocks.UserRepository)
	mockUser := &models.User{
		ID:       uuid.New(),
		Username: "ValidUser",
	}

	t.Run("successful registration", func(t *testing.T) {
		mockUserRepository.On("Create", mock.Anything, mockUser).Return(nil).Once()

		serv := service.NewSignUpService(mockUserRepository, nil, time.Second*2)
		err := serv.RegisterUser(context.Background(), mockUser)

		assert.NoError(t, err)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("invalid username", func(t *testing.T) {
		invalidUser := &models.User{
			ID:       uuid.New(),
			Username: "Invalid User 123!", // Некорректное имя пользователя
		}

		serv := service.NewSignUpService(mockUserRepository, nil, time.Second*2)
		err := serv.RegisterUser(context.Background(), invalidUser)

		assert.Error(t, err)
		assert.Equal(t, models.ErrInvalidUsername, err)
	})

	t.Run("registration failure", func(t *testing.T) {
		mockUserRepository.On("Create", mock.Anything, mockUser).Return(models.ErrUsernameExists).Once()

		serv := service.NewSignUpService(mockUserRepository, nil, time.Second*2)
		err := serv.RegisterUser(context.Background(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, models.ErrUsernameExists, err)

		mockUserRepository.AssertExpectations(t)
	})
}
