package repository

import (
	"context"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *pgx.Conn
}

// GetUserByEmail implements models.UserRepository.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	logger.Logger.Infof("Запрос юзера по имейлу: %s", email)
	query := `
			SELECT 
					u.id, u.username, u.password, u.salt, u.refresh_token
			FROM
					users u
			WHERE   
					u.email = $1
`

	rows, err := u.db.Query(ctx, query, email)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в GetUserByEmail")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID uuid.UUID
		var username, userPass, userSalt, userRefreshToken string

		err := rows.Scan(&userID, &username, &userPass, &userSalt, &userRefreshToken)
		if err != nil {
			logger.Logger.WithError(err).Error("Ошибка сканирования строки в GetUserByEmail")
			return nil, err
		}

		if user == nil {
			user = &models.User{
				ID:           userID,
				Username:     username,
				Email:        email,
				Salt:         userSalt,
				RefreshToken: userRefreshToken,
			}
		}
	}

	if err = rows.Err(); err != nil {
		logger.Logger.WithError(err).Error("Ошибка после обработки строк в GetUserByEmail")
		return nil, err
	}

	return user, nil
}

// CheckRefreshTokenExists implements models.UserRepository.
func (u *userRepository) CheckRefreshTokenExists(ctx context.Context, hashedToken string) (bool, error) {
	logger.Logger.Infof("проверка сущестует ли рефшен токен")
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE refresh_token = $1)"

	var exists bool
	err := u.db.QueryRow(ctx, query, hashedToken).Scan(&exists)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в CheckRefreshTokenExists")
		return false, err
	}

	logger.Logger.Infof("Результат проверки существования рефреш токена: %v", exists)
	return exists, nil
}

// GetUserByUsername implements models.UserRepository.
func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	logger.Logger.Infof("Запрос юзера по юзернейму: %s", username)
	query := `
		SELECT 
			u.id, u.password, u.email, u.salt, u.refresh_token
		FROM
			users u
		WHERE
			u.username = $1
	`

	rows, err := u.db.Query(ctx, query, username)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в GetUserByUsername")
		return nil, err
	}
	defer rows.Close()

	var user *models.User

	for rows.Next() {
		var userID uuid.UUID
		var userPass, userSalt, userEmail, userRefreshToken string

		err := rows.Scan(&userID, &userPass, &userEmail, &userSalt, &userRefreshToken)
		if err != nil {
			logger.Logger.WithError(err).Error("Ошибка сканирования в GetUserByUsername")
			return nil, err
		}

		if user == nil {
			user = &models.User{
				ID:           userID,
				Username:     username,
				Password:     userPass,
				Email:        userEmail,
				Salt:         userSalt,
				RefreshToken: userRefreshToken,
			}
		}
	}

	if err = rows.Err(); err != nil {
		logger.Logger.WithError(err).Error("Ошибка после обработки строк в GetUserByUsername")
		return nil, err
	}

	return user, nil
}

// DeleteUserRefreshTokenByEmail implements models.UserRepository.
func (u *userRepository) DeleteUserRefreshTokenByEmail(ctx context.Context, email string) (err error) {
	logger.Logger.Infof("Удаление рефреш токена по имейлу: %s", email)
	query := `
	UPDATE users
	SET refresh_token = NULL
	WHERE email = $1
	`
	_, err = u.db.Exec(ctx, query, email)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в DeleteUserRefreshTokenByEmail")
		return err
	}

	logger.Logger.Info("Рефреш токен успешно удален")
	return nil
}

// GetRefreshToken implements models.UserRepository.
func (u *userRepository) GetRefreshToken(ctx context.Context, email string) (hashedToken string, err error) {
	logger.Logger.Infof("Получение рефреш токена по имейлу: %s", email)
	query := ` 
	SELECT refresh_token 
	FROM users
	WHERE email = $1;
`

	var hashedRefreshToken string
	err = u.db.QueryRow(ctx, query, email).Scan(&hashedRefreshToken)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в GetRefreshToken")
		return "", err
	}

	logger.Logger.Info("Рефреш токен успешно получен")
	return hashedRefreshToken, nil
}

// SaveRefreshToken implements models.UserRepository.
func (u *userRepository) SaveRefreshToken(id uuid.UUID, hashedToken string) error {
	logger.Logger.Infof("Сохранение рефреш токена для юзера с ID: %s", id)
	query := `
		UPDATE users
		SET refresh_token = $1
		WHERE id = $2;
	`

	_, err := u.db.Exec(context.Background(), query, hashedToken, id)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в SaveRefreshToken")
		return err
	}

	logger.Logger.Info("Рефреш токен успешно сохранен")
	return nil
}

// CheckUsernameExists implements models.UserRepository.
func (u *userRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	logger.Logger.Infof("Проверка существования юзернейма: %s", username)
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"

	var exists bool
	err := u.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в CheckUsernameExists")
		return false, err
	}

	logger.Logger.Infof("статус юзернейма: %v", exists)
	return exists, nil
}

// Create implements models.UserRepository.
func (u *userRepository) Create(ctx context.Context, user *models.User) error {
	logger.Logger.Infof("Создание юзера: %s", user.Username)
	query := `
		INSERT INTO users (id, username, email, password, refresh_token, salt, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := u.db.Exec(ctx, query, user.ID, user.Username, user.Email, user.Password, user.RefreshToken, user.Salt, user.CreatedAt)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка запроса в Create (user repository)")
		return err
	}

	logger.Logger.Info("Пользователь успешно создан")
	return nil
}

func NewUserRepository(db *pgx.Conn) models.UserRepository {
	return &userRepository{db: db}
}
