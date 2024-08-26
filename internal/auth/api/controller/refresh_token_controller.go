package controller

import (
	// "log"
	"net/http"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/Seven11Eleven/auth_service_medods/internal/logger"
	"github.com/Seven11Eleven/auth_service_medods/internal/utils"
	"github.com/gin-gonic/gin"
)

type RefreshTokenController struct {
	RefreshTokenService models.RefreshTokenService
	Env                 *config.Env
}

func (rtc *RefreshTokenController) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	err := c.ShouldBind(&req)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка биндинга запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := rtc.RefreshTokenService.ExtractEmailFromRefreshToken(req.RefreshToken)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка извлечения email из рефреш токена")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Infof("Извлечен email из рефреш токена: %s", email)

	registeredIP, err := rtc.RefreshTokenService.ExtractIPFromRefreshToken(email, req.RefreshToken)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка извлечения IP адреса из рефреш токена")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := rtc.RefreshTokenService.GetUserByEmail(c, email)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка получения пользователя по email")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	currentIP := c.ClientIP()
	logger.Logger.Infof("Текущий IP: %s, Зарегистрированный IP: %s", currentIP, registeredIP)

	if currentIP != registeredIP {
		err := utils.SendWarningEmail(email, currentIP, registeredIP, user.Username)
		if err != nil {
			logger.Logger.WithError(err).Error("Ошибка отправки warning email")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.Logger.Warnf("Обнаружено несовпадение IP адресов для пользователя %s. Предупреждающее письмо отправлено на %s", user.Username, email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Мы обнаружили, что с вашего аккаунта отправлиись запросы с другого айпи адреса, отличного от того с которым вы регистрировались "})
		return
	}

	newAccessToken, err := rtc.RefreshTokenService.CreateAccessToken(user, rtc.Env.AccessTokenExpiryHour)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка создания нового акцес токена")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newRefreshToken, err := rtc.RefreshTokenService.CreateRefreshToken(user, rtc.Env.RefreshTokenExpiryHour)
	if err != nil {
		logger.Logger.WithError(err).Error("ошибка создания нового рефреш токена")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Infof("юзер %s успешно обновил токены", user.Username)

	loginResp := models.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	c.JSON(http.StatusOK, loginResp)
}
