package providers

import (
	"context"
	"fmt"
	"net/smtp"
	"notification-service/internal/core/models"
)

type EmailNotifier struct {
	host        string
	port        string
	username    string
	password    string
	fromAddress string
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

func (e *EmailNotifier) Send(ctx context.Context, userID string, message string) error {
	type result struct{ err error }
	ch := make(chan result, 1)

	go func() {
		toAddress := userID
		auth := smtp.PlainAuth("", e.username, e.password, e.host)
		subject := "Notification"
		emailBody := fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			e.fromAddress, toAddress, subject, message,
		)
		addr := fmt.Sprintf("%s:%s", e.host, e.port)
		err := smtp.SendMail(addr, auth, e.fromAddress, []string{toAddress}, []byte(emailBody))
		if err != nil {
			ch <- result{fmt.Errorf("[EMAIL] failed to send to %s: %w", toAddress, err)}
			return
		}
		fmt.Printf("[EMAIL] Successfully sent to %s: %s\n", toAddress, message)
		ch <- result{nil}
	}()

	select {
	case r := <-ch:
		return r.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *EmailNotifier) Supports() models.ChannelType {
	return models.Email
}
