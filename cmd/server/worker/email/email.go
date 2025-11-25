package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
)

//go:embed templates/*.html
var templateFS embed.FS

type SMTPService struct {
	config config.SMTP
}

func NewSMTPService(cfg config.SMTP) *SMTPService {
	return &SMTPService{
		config: cfg,
	}
}

func (s *SMTPService) Send(from string, to []string, subject, templateName string, data any) error {
	// Parse template from embedded FS
	t, err := template.ParseFS(templateFS, "templates/"+templateName)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	auth := smtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := fmt.Appendf(nil, "To: %s\r\nSubject: %s\r\n%s\r\n%s", strings.Join(to, ", "), subject, headers, body.String())

	if err := smtp.SendMail(s.config.DSN(), auth, from, to, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
