package utils

import (
	"fmt"
	"net/smtp"

	"github.com/Seven11Eleven/auth_service_medods/internal/config"
)

func SendWarningEmail(userEmail, userIP, registeredIP, username string) error {
	env := config.NewEnv()
	smtpHost := env.SmtpHost
	smtpPort := env.SmtpPort
	senderEmail := env.SmtpEmail
	senderPassword := env.SmtpPassword

	subject := "email warning!"
	body := fmt.Sprintf("Дорогой %v,\n\nМы обнаружили попытку войти в ваш аккаунт с другого айпи адреса, нежели который вы нам давали.\n\nЗАрегистрированый айпи адрес : %s\nАйпи адрес с которого была попытка войти в аккаунт: %s\n\nПОжалуйста, есди это были не вы, то измените пароль.\n\nВСего наилучшего ,\nSeven11Eleven & Co", username, registeredIP, userIP)
	message := []byte("Subject: " + subject + "\r\n" +
		"From: " + senderEmail + "\r\n" +
		"To: " + userEmail + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{userEmail}, message)
	if err != nil {
		return err
	}

	return nil
}
