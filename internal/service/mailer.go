package service

import (
	"fmt"
	"os"
	"time"

	mailgun "github.com/mailgun/mailgun-go"
)

var (
	domain      = os.Getenv("MG_DOMAIN")
	mgAPIKey    = os.Getenv("MG_API_KEY")
	msgTemplate = `Dear %s
Your password has been successfully reset to: %s. Please change as soon as possible.`
)

// SendResetPasswordEmail attempts to send out an email using mailgun contaning the new password.
func SendResetPasswordEmail(email, newPassword string) error {
	mg := mailgun.NewMailgun(domain, mgAPIKey)

	sender := fmt.Sprintf("no-reply@%s", domain)
	subject := fmt.Sprintf("[%s] IdleRPG Notifier", time.Now().Format("2006-01-02"))
	body := fmt.Sprintf(msgTemplate, email, newPassword)

	message := mg.NewMessage(sender, subject, body, email)
	_, _, err := mg.Send(message)
	return err
}
