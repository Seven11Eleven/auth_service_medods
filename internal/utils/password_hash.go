package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
)

const saltSize = 16

func GenerateSalt() (string, error) {
	salt := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func HashPassword(password string) (string, error) {
	if len(password) > 72 { //72 потому что во время дебага я столкнулся с такой проблемой "error": "bcrypt: password length exceeds 72 bytes", решено было просто не давать сделать паорль длинее 72 байтов
		log.Printf("Pass len: %d", len(password))
		return "", bcrypt.ErrPasswordTooLong
	}
	log.Printf("hashpass: %s", password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	log.Printf("Generated hash password: %s", string(hashedPassword))
	return string(hashedPassword), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	log.Printf("pass: %s\n hash pass: %s", password, hashedPassword)
	if len(password) > 72 {
		log.Printf("len pass: %d", len(password))
		return bcrypt.ErrPasswordTooLong
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type PasswordHasher interface {
	GenerateSalt() (string, error)
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

type hashUtilsImpl struct{}

func (h *hashUtilsImpl) GenerateSalt() (string, error) {
	return GenerateSalt()
}

func (h *hashUtilsImpl) HashPassword(password string) (string, error) {
	return HashPassword(password)
}

func (h *hashUtilsImpl) CompareHashAndPassword(hashedPassword, password string) error {
	return CompareHashAndPassword(hashedPassword, password)
}

func NewHashUtils() PasswordHasher {
	return &hashUtilsImpl{}
}
