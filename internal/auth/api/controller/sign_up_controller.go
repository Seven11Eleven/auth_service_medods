package controller

import (
	"net/http"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/Seven11Eleven/auth_service_medods/internal/logger"
	"github.com/Seven11Eleven/auth_service_medods/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SignUpController struct {
	SignUpService models.SignUpService
	Env           *config.Env
}

func (sc *SignUpController) SignUp(c *gin.Context) {
	var req models.SignUpRequest

	err := c.ShouldBind(&req)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка биндинга запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := sc.SignUpService.CheckUsernameExists(c, req.Username)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка проверки существования юзернейма")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username"})
		return
	}

	if exists {
		logger.Logger.Warnf("Юзернейм %s уже существует", req.Username)
		c.JSON(http.StatusConflict, gin.H{"error": models.ErrUsernameExists.Error()})
		return
	}

	userSalt, err := utils.GenerateSalt()
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка генерации соли")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save your password"})
		return
	}

	localSalt := sc.Env.Salt

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка хеширования пароля")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save yoru password"})
		return
	}

	hashash := localSalt + req.Password + userSalt
	logger.Logger.Debugf("userpass: %v, userSalt: %v, localsalt: %v, fullhash: %v", req.Password, userSalt, localSalt, hashash)

	user := &models.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Salt:     userSalt,
	}

	err = sc.SignUpService.RegisterUser(c, user)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка регистрации юзера")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Infof("юзер %s успешно зарегистрирован", user.Username)

	c.JSON(http.StatusCreated, models.SuccessfullResponse{
		Message: "you signed up successfully!",
	})
}