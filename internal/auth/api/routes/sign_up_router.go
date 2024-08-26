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

func NewSignUpRouter(env *config.Env, refreshTokenService models.RefreshTokenService, timeout time.Duration, db *pgx.Conn, group *gin.RouterGroup){
	ur := repository.NewUserRepository(db)
	lc := &controller.SignUpController{
		SignUpService: service.NewSignUpService(ur, refreshTokenService, timeout),
		Env: env,
	}
	group.POST("/signup", lc.SignUp)
}