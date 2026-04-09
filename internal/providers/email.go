package providers

import (
	"fmt"
	"net/smtp"
	"notification-service/internal/core/models"
)

type EmailNotifier struct {
	host        string // smtp email
	port        string // port number 
	username    string // smtp login (your email address)
	password    string // smtp password
	fromAddress string // sender email : user.tuskira.ai@gmail.com
}

func NewEmailNotifier(host, port, username, password, fromAddress string) *EmailNotifier {
	return &EmailNotifier{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		fromAddress: fromAddress,
	}
}

func (e *EmailNotifier) Send(userID string, message string) error {

	/*
		 For this now i am assuming userID is treated directly as the recipient email.
		 rest for the production i would be fetching the email address using the userID 
		 from database
	*/
	
	toAddress := userID

	/* smtp server authentication  */
	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	subject := "Notification"
	emailBody := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.fromAddress, toAddress, subject, message,
	)

	// smtp mail delivery logic
	addr := fmt.Sprintf("%s:%s", e.host, e.port)
	err := smtp.SendMail(addr, auth, e.fromAddress, []string{toAddress}, []byte(emailBody))
	if err != nil {
		return fmt.Errorf("[EMAIL] failed to send to %s: %w", toAddress, err)
	}

	fmt.Printf("[EMAIL] Successfully sent to %s: %s\n", toAddress, message)
	return nil
}

func (e *EmailNotifier) Supports() models.ChannelType {
	return models.Email
}
