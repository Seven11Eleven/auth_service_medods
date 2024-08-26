package utils

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/auth/models"
	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTUtils interface {
	CreateAccessToken(user *models.User, expired int) (string, error)
	CreateRefreshToken(user *models.User, expired int) (string, error)
	ExtractIPFromRefreshToken(email string, originalToken string) (string, error)
	ExtractIDFromToken(requestedToken string) (uuid.UUID, error)
	ExtractEmailFromRefreshToken(originalToken string) (string, error)
	IsAuthorized(token string) (bool, error)
}

type jwtUtilsImp struct {
	env            *config.Env
	userRepository models.UserRepository
}

// ExtractEmailFromRefreshToken implements JWTUtils.
func (j *jwtUtilsImp) ExtractEmailFromRefreshToken(originalToken string) (string, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(originalToken)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 token: %w", err)
	}

	var payload models.CustomRefreshClaims
	err = json.Unmarshal(tokenBytes, &payload)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling token JSON: %w", err)
	}
	if payload.Email == "" {
		return "", fmt.Errorf("email is empty in the token payload")
	}

	return payload.Email, nil
}

// ExtractIDFromToken implements JWTUtils.
func (j *jwtUtilsImp) ExtractIDFromToken(requestedToken string) (uuid.UUID, error) {
	token, err := jwt.Parse(requestedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неправильный метод подписи: %v", token.Header["alg"])
		}
		return []byte(j.env.JWTSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка при разборе токена: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("невалидный токне")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("не удалось извлечь ID из токена")
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("невалидный формат: %v", err)
	}

	return parsedUUID, nil
}

// IsAuthorized implements JWTUtils.
func (j *jwtUtilsImp) IsAuthorized(token string) (bool, error) {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error!!: %v", token.Header["alg"])
		}
		return []byte(j.env.JWTSecret), nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateAccessToken implements JWTUtils.
func (j *jwtUtilsImp) CreateAccessToken(user *models.User, expired int) (string, error) {
	claims := &models.JWTCustomClaims{
		Username:  user.Username,
		ID:        user.ID,
		IPAddress: user.IPAddress,
		Email:     user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expired))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) //Access токен тип JWT, алгоритм SHA512 HS512 == SHA512
	accessToken, err := token.SignedString([]byte(j.env.JWTSecret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// CreateRefreshToken implements JWTUtils.
func (j *jwtUtilsImp) CreateRefreshToken(user *models.User, expired int) (string, error) {

	tokenPayload := models.CustomRefreshClaims{
		ID:        user.ID,
		Email:     user.Email,
		IPAddress: user.IPAddress,
		CreatedAt: time.Now(),
	}

	tokenPayloadBytes, err := json.Marshal(tokenPayload)
	if err != nil {
		return "", err
	}

	encodedToken := base64.URLEncoding.EncodeToString(tokenPayloadBytes)
	fmt.Printf("Encoded token: %s\n", encodedToken)

	shaHash := sha256.New() //sha256 было решено использовать, чтобы понизить кол-во байтов до 32, ибо опять таки выходила ошибка при дебаге "error": "bcrypt: password length exceeds 72 bytes", sha256 позволяет сохранять токены с большим пейлоадом, как у меня в bcrypt, у которого как оказалось есть лимит в 72 байта
	shaHash.Write([]byte(encodedToken))
	shaHashedToken := shaHash.Sum(nil)
	shaHashedTokenStr := hex.EncodeToString(shaHashedToken)
	fmt.Printf("SHA-256 hashed token: %s\n", shaHashedTokenStr)

	bcryptHashedToken, err := bcrypt.GenerateFromPassword([]byte(shaHashedTokenStr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	fmt.Printf("Bcrypt hashed token: %s\n", string(bcryptHashedToken))

	err = j.userRepository.SaveRefreshToken(user.ID, string(bcryptHashedToken))
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}

// ExtractIPFromToken implements JWTUtils.
func (j *jwtUtilsImp) ExtractIPFromRefreshToken(email string, originalToken string) (string, error) {
	shaHash := sha256.New()
	shaHash.Write([]byte(originalToken))
	shaHashedToken := shaHash.Sum(nil)
	shaHashedTokenStr := hex.EncodeToString(shaHashedToken)

	hashedRefreshToken, err := j.userRepository.GetRefreshToken(context.Background(), email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedRefreshToken), []byte(shaHashedTokenStr))
	if err != nil {
		return "", err
	}

	tokenBytes, err := base64.URLEncoding.DecodeString(originalToken)
	if err != nil {
		return "", err
	}
	var payload models.CustomRefreshClaims
	err = json.Unmarshal(tokenBytes, &payload)
	if err != nil {
		return "", err
	}

	return payload.IPAddress, nil
}

func NewJWTUtils(env *config.Env, userRepository models.UserRepository) JWTUtils {
	return &jwtUtilsImp{
		env:            env,
		userRepository: userRepository,
	}
}
