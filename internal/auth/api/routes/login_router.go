package routes

import (
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/api/controller"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/repository"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/service"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func NewLoginRouter(env *config.Env, refreshTokenService models.RefreshTokenService, timeout time.Duration, db *pgx.Conn, group *gin.RouterGroup){
	ur := repository.NewUserRepository(db)
	lc := &controller.LoginController{
		LoginService: service.NewLoginService(ur,refreshTokenService, timeout),
		Env: env,
	}
	group.POST("/login", lc.Login)
	group.GET("/tokenByGUID", lc.TokenByGUID)
}