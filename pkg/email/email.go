package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

const (
	VerificationEmail = "verification_email"
)

type SendEmailRequest struct {
	To      []string
	Type    string
	Body    map[string]string
	Subject string
}

func SendEmail(config *SmtpConfig, req *SendEmailRequest) error {
	from := config.Sender
	to := req.To

	password := config.Password

	var body bytes.Buffer

	templatePath := getTemplatePath(req.Type)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	t.Execute(&body, req.Body)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: %s\n", req.Subject)
	msg := []byte(subject + mime + body.String())

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}

func getTemplatePath(emailType string) string {
	switch emailType {
	case VerificationEmail:
		return "./templates/verification_email.html"
	default:
		return ""
	}
}
