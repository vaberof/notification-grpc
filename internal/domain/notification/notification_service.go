package notification

import (
	"context"
	"github.com/vaberof/notification-grpc/pkg/email"
	"github.com/vaberof/notification-grpc/pkg/logging/logs"
	"log/slog"
)

type NotificationService interface {
	SendEmail(ctx context.Context, to string, emailType string, subject string, body map[string]string) error
}

type SmtpConfig struct {
	SenderEmail string
	Password    string
}

type notificationServiceImpl struct {
	config *SmtpConfig

	logger *slog.Logger
}

func NewNotificationService(config *SmtpConfig, logs *logs.Logs) NotificationService {
	logger := logs.WithName("domain.notification.service")
	return &notificationServiceImpl{config: config, logger: logger}
}

func (n *notificationServiceImpl) SendEmail(ctx context.Context, to string, emailType string, subject string, body map[string]string) error {
	const operation = "SendEmail"

	log := n.logger.With(
		slog.String("operation", operation),
		slog.String("to", to))

	err := email.SendEmail(
		&email.SmtpConfig{Sender: n.config.SenderEmail, Password: n.config.Password},
		&email.SendEmailRequest{To: []string{to}, Type: emailType, Body: body, Subject: subject})
	if err != nil {
		log.Error("failed to send email", err)

		return err
	}

	log.Info("email has been sent successfully")

	return nil
}
