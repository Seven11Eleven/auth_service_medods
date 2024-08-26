package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ServerPort             string `mapstructure:"PORT"`
	DBName                 string `mapstructure:"DATABASE_NAME"`
	DBPass                 string `mapstructure:"DATABASE_PASSWORD"`
	DBUser                 string `mapstructure:"DATABASE_USER"`
	DBHost                 string `mapstructure:"DATABASE_HOST"`
	DBPort                 string `mapstructure:"DATABASE_PORT"`
	JWTSecret              string `mapstructure:"JWT_SECRET"`
	Salt                   string `mapstructure:"SALT"`
	SmtpHost               string `mapstructure:"SMTP_HOST"`
	SmtpPort               string `mapstructure:"SMTP_PORT"`
	SmtpEmail              string `mapstructure:"SMTP_EMAIL"`
	SmtpPassword           string `mapstructure:"SMTP_PASS"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_EXPIRY_HOUR"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile("/home/seveneleven/go/src/github.com/Seven11Eleven/auth_service_medods/.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalf("не получается загрузить с файла окружений: %v", err)
	}
	return &env
}
