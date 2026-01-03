package email

import (
	"fmt"

	"github.com/semmidev/ethos-go/config"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"gopkg.in/mail.v2"
)

type SMTPClient struct {
	cfg    *config.Config
	logger logger.Logger
	dialer *mail.Dialer
}

func NewSMTPClient(cfg *config.Config, l logger.Logger) (*SMTPClient, error) {
	dialer := mail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS

	return &SMTPClient{
		cfg:    cfg,
		dialer: dialer,
		logger: l,
	}, nil
}

func (s *SMTPClient) Send(recipient, subject string, htmlContent string, data any) error {
	m := mail.NewMessage()

	m.SetHeader("From", m.FormatAddress(s.cfg.SMTPUser, s.cfg.AppName))
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContent)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
