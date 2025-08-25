package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"trade_company/internal/config"
	"trade_company/internal/models"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(config *config.Config) *EmailService {
	return &EmailService{
		config: config,
	}
}

// GenerateVerificationToken generates a random verification token
func (es *EmailService) GenerateVerificationToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GeneratePasswordResetToken generates a random password reset token
func (es *EmailService) GeneratePasswordResetToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SendVerificationEmail sends an email verification email
func (es *EmailService) SendVerificationEmail(user *models.User, verificationToken string) error {
	// In development, just log the email
	if es.config.AppEnv == "development" {
		es.logEmail(user.Email, "Verify Your Email - Business Exchange",
			es.generateVerificationEmailText(user.FirstName, verificationToken))
		return nil
	}

	// TODO: Implement SendGrid integration
	// For now, just log the email
	es.logEmail(user.Email, "Verify Your Email - Business Exchange",
		es.generateVerificationEmailText(user.FirstName, verificationToken))
	return nil
}

// SendPasswordResetEmail sends a password reset email
func (es *EmailService) SendPasswordResetEmail(user *models.User, resetToken string) error {
	// In development, just log the email
	if es.config.AppEnv == "development" {
		es.logEmail(user.Email, "Reset Your Password - Business Exchange",
			es.generatePasswordResetEmailText(user.FirstName, resetToken))
		return nil
	}

	// TODO: Implement SendGrid integration
	// For now, just log the email
	es.logEmail(user.Email, "Reset Your Password - Business Exchange",
		es.generatePasswordResetEmailText(user.FirstName, resetToken))
	return nil
}

// SendLeadNotification sends a notification to a seller about a new lead
func (es *EmailService) SendLeadNotification(seller *models.User, lead *models.Lead) error {
	subject := fmt.Sprintf("New Lead: %s", lead.Subject)

	// In development, just log the email
	if es.config.AppEnv == "development" {
		es.logEmail(seller.Email, subject,
			es.generateLeadNotificationText(seller.FirstName, lead))
		return nil
	}

	// TODO: Implement SendGrid integration
	// For now, just log the email
	es.logEmail(seller.Email, subject,
		es.generateLeadNotificationText(seller.FirstName, lead))
	return nil
}

// logEmail logs email content in development mode
func (es *EmailService) logEmail(to, subject, textContent string) {
	fmt.Printf("=== EMAIL LOG ===\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("Text Content:\n%s\n", textContent)
	fmt.Printf("================\n")
}

// generateVerificationEmailText generates text content for verification email
func (es *EmailService) generateVerificationEmailText(firstName, verificationToken string) string {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", es.config.AppName, verificationToken)

	return fmt.Sprintf(`Welcome to Business Exchange!

Hi %s,

Thank you for signing up! Please verify your email address by visiting this link:

%s

This link will expire in 24 hours.

Best regards,
The Business Exchange Team`, firstName, verificationURL)
}

// generatePasswordResetEmailText generates text content for password reset email
func (es *EmailService) generatePasswordResetEmailText(firstName, resetToken string) string {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", es.config.AppName, resetToken)

	return fmt.Sprintf(`Reset Your Password

Hi %s,

We received a request to reset your password. Visit this link to create a new password:

%s

If you didn't request this, you can safely ignore this email.

This link will expire in 30 minutes.

Best regards,
The Business Exchange Team`, firstName, resetURL)
}

// generateLeadNotificationText generates text content for lead notification
func (es *EmailService) generateLeadNotificationText(firstName string, lead *models.Lead) string {
	return fmt.Sprintf(`New Lead Received!

Hi %s,

You have received a new lead from a potential buyer:

Subject: %s
From: %s %s
Message: %s
Contact Phone: %s

Log in to your dashboard to respond to this lead.

Best regards,
The Business Exchange Team`, firstName, lead.Subject, lead.Sender.FirstName, lead.Sender.LastName, lead.Message, lead.ContactPhone)
}
