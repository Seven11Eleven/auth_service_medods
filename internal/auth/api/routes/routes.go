package routes

import (
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes( env config.Env, refreshTokenService models.RefreshTokenService, timeout time.Duration, db *pgx.Conn, gin *gin.Engine){
	publicRouter := gin.Group("")

	NewSignUpRouter(&env, refreshTokenService, timeout, db, publicRouter)
	NewLoginRouter(&env, refreshTokenService, timeout, db, publicRouter)
	NewRefreshTokenRouter(&env, refreshTokenService ,timeout, db, publicRouter)

}