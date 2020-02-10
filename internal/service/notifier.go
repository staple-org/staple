package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go"

	"github.com/staple-org/staple/pkg/config"
)

// Event represents an event that can happen which needs the user's attention via
// prefered notification medium.
type Event string

// Payload is data which the notifier carries.
type Payload string

var (
	// PasswordReset is an event that happens when the user's password is reset.
	PasswordReset Event = "Password Reset"
	// GenerateConfirmCode is an event before password reset which sends a confirm code to the user's email address.
	GenerateConfirmCode Event = "Confirm Code"
)

// Notifier notifies the user of some event.
type Notifier interface {
	Notify(email string, event Event, payload string) error
}

// EmailNotifier is an email based notification entity.
type EmailNotifier struct{}

// NewEmailNotifier creates a new email notifier.
func NewEmailNotifier() EmailNotifier {
	return EmailNotifier{}
}

var (
	domain                = config.Opts.Mailer.Domain
	mgAPIKey              = config.Opts.Mailer.APIKey
	passwordResetTemplate = `Dear %s
Your password has been successfully reset to: %s. Please change as soon as possible.`
	confirmCodeTemplate = `Dear %s
Please enter the following code into the confirm code window: %s`
)

// Notify attempts to send out an email using mailgun contaning the new password.
// Does not need to be a pointer receiver because it isn't storing anything.
func (e EmailNotifier) Notify(email string, event Event, payload string) error {
	sender := fmt.Sprintf("no-reply@%s", domain)
	subject := fmt.Sprintf("[%s] %s Notification", time.Now().Format("2006-01-02"), event)

	if domain == "" && mgAPIKey == "" {
		config.Opts.Logger.Warn().Msg("[WARNING] Mailgun not set up. Falling back to console output...")
		config.Opts.Logger.Info().Str("email", email).Str("subject", subject).Str("payload", payload).Msg("A notification attempt was made for user.")
		return nil
	}

	var body string
	switch event {
	case PasswordReset:
		body = fmt.Sprintf(passwordResetTemplate, email, payload)
	case GenerateConfirmCode:
		body = fmt.Sprintf(confirmCodeTemplate, email, payload)
	}

	mg := mailgun.NewMailgun(domain, mgAPIKey)
	message := mg.NewMessage(sender, subject, body, email)
	_, _, err := mg.Send(message)
	return err
}

// BufferNotifier uses a byte buffer for notifications.
type BufferNotifier struct {
	buffer bytes.Buffer
}

// NewBufferNotifier creates a new notifier.
func NewBufferNotifier() *BufferNotifier {
	return &BufferNotifier{}
}

// Notify uses a buffer to store notifications for a user.
func (b *BufferNotifier) Notify(email string, event Event, payload string) error {
	var body string
	switch event {
	case PasswordReset:
		body = fmt.Sprintf(passwordResetTemplate, email, payload)
	case GenerateConfirmCode:
		body = fmt.Sprintf(confirmCodeTemplate, email, payload)
	}
	b.buffer.WriteString(body)
	return nil
}
