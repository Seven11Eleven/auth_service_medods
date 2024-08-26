package routes

import (
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/api/controller"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/repository"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/service"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/Seven11Eleven/auth_service_medods/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func NewRefreshTokenRouter(env *config.Env, refreshTokenService utils.JWTUtils, timeout time.Duration, db *pgx.Conn, group *gin.RouterGroup){
	ur := repository.NewUserRepository(db)
	lc := &controller.RefreshTokenController{
		RefreshTokenService: service.NewRefreshTokenService(ur, refreshTokenService,timeout),
		Env: env,
	}
	group.POST("/refresh", lc.RefreshToken)
}