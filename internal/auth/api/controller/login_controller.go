package controller

import (
	"net/http"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/Seven11Eleven/auth_service_medods/internal/logger"
	"github.com/Seven11Eleven/auth_service_medods/internal/utils"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginService models.LoginService
	Env          *config.Env
}

func (lc *LoginController) Login(c *gin.Context) {
	var request models.LoginRequest
	localSalt := lc.Env.Salt

	err := c.ShouldBind(&request)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка бинжа запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := lc.LoginService.GetUserByUsername(c, request.Username)
	if err != nil {
		logger.Logger.WithError(err).Error("Пользователь не найден")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	hashash := localSalt + request.Password + user.Salt
	logger.Logger.Debugf("password: %v, userpass: %v, usersalt: %v, localsalt: %v, fullhash: %v", request.Password, user.Password, user.Salt, localSalt, hashash)

	if err := utils.CompareHashAndPassword(user.Password, request.Password); err != nil {
		logger.Logger.WithError(err).Error("Неверный пароль")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user.IPAddress = c.ClientIP()
	logger.Logger.Infof("АЙпи-адрес юзера: %s", user.IPAddress)

	accessToken, err := lc.LoginService.CreateAccessToken(user, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка создания акцесс токена")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "дебаг1"})
		return
	}

	refreshToken, err := lc.LoginService.CreateRefreshToken(user, lc.Env.RefreshTokenExpiryHour)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка создания рефрещ токена")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "дебаг2"})
		return
	}

	logger.Logger.Infof("юзер %s успешно вошел в систему", user.Username)

	loginResp := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, loginResp)
}