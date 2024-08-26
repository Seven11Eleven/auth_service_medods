package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/api/routes"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/repository"
	"github.com/Seven11Eleven/auth_service_medods/internal/auth/service"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/Seven11Eleven/auth_service_medods/internal/database"
	"github.com/Seven11Eleven/auth_service_medods/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type App struct {
	Server *gin.Engine
	DB     *pgx.Conn
	Env    *config.Env
}

func NewApp(ctx context.Context) (*App, error) {
	env := config.NewEnv()

	db := database.NewPostgreSQLConnection(env)

	ginEngine := gin.Default()

	userRepo := repository.NewUserRepository(db)
	jwtUtils := utils.NewJWTUtils(env, userRepo)

	refreshTokenService := service.NewRefreshTokenService(userRepo, jwtUtils, 10*time.Second)

	routes.SetupRoutes(*env, refreshTokenService, 10*time.Second, db, ginEngine)

	return &App{
		Server: ginEngine,
		DB:     db,
		Env:    env,
	}, nil
}

func (a *App) Run() error {
	port := a.Env.ServerPort
	if port == "" {
		fmt.Println("походу что то не загрузилось")
		port = "8080" //I'm a genius
	}

	return a.Server.Run(fmt.Sprintf(":%s", port))
}

func (a *App) Close() {
	database.ClosePostgreSQLConnection(a.DB)
}
