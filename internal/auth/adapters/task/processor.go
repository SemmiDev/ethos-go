package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/internal/auth/domain/gateway"
	"github.com/semmidev/ethos-go/internal/common/assets"
	"github.com/semmidev/ethos-go/internal/common/email"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// TaskProcessor handles processing of auth-related background tasks
type TaskProcessor struct {
	logger logger.Logger
	email  email.Email
}

func NewTaskProcessor(l logger.Logger, email email.Email) *TaskProcessor {
	return &TaskProcessor{
		logger: l,
		email:  email,
	}
}

func (p *TaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload gateway.PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		p.logger.Error(ctx, err, "failed to unmarshal payload")
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	var tpl *template.Template
	tpl, err := template.ParseFS(assets.EmbeddedFiles, assets.EmailVerificationTemplatePath)
	if err != nil {
		p.logger.Error(ctx, err, "failed to parse email template")
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := tpl.ExecuteTemplate(&body, "htmlBody", payload); err != nil {
		p.logger.Error(ctx, err, "failed to execute email template")
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	htmlContent := body.String()
	subject := payload.Subject

	err = p.email.Send(payload.Email, subject, htmlContent, payload)
	if err != nil {
		p.logger.Error(ctx, err, "failed to send verify email")
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	p.logger.Info(ctx, "verify email sent", logger.Field{Key: "email", Value: payload.Email})
	return nil
}

func (p *TaskProcessor) ProcessTaskSendForgotPasswordEmail(ctx context.Context, task *asynq.Task) error {
	var payload gateway.PayloadSendForgotPasswordEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		p.logger.Error(ctx, err, "failed to unmarshal payload")
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	var tpl *template.Template
	tpl, err := template.ParseFS(assets.EmbeddedFiles, assets.EmailForgotPasswordTemplatePath)
	if err != nil {
		p.logger.Error(ctx, err, "failed to parse forgot password email template")
		return fmt.Errorf("failed to parse forgot password email template: %w", err)
	}

	var body bytes.Buffer
	if err := tpl.ExecuteTemplate(&body, "htmlBody", payload); err != nil {
		p.logger.Error(ctx, err, "failed to execute forgot password email template")
		return fmt.Errorf("failed to execute forgot password email template: %w", err)
	}

	htmlContent := body.String()
	subject := payload.Subject

	err = p.email.Send(payload.Email, subject, htmlContent, payload)
	if err != nil {
		p.logger.Error(ctx, err, "failed to send forgot password email")
		return fmt.Errorf("failed to send forgot password email: %w", err)
	}

	p.logger.Info(ctx, "forgot password email sent", logger.Field{Key: "email", Value: payload.Email})
	return nil
}
