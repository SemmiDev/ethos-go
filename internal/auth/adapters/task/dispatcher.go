package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/config"
	"github.com/semmidev/ethos-go/internal/auth/domain/gateway"
)

const (
	TaskSendVerifyEmail         = "task:send_verify_email"
	TaskSendForgotPasswordEmail = "task:send_forgot_password_email"

	TaskSendForgotPasswordEmailSubject = "Permintaan Reset Password"
	TaskSendVerifyEmailSubject         = "Verifikasi Email"
)

// AsynqTaskDispatcher implements TaskDispatcher using Asynq
type AsynqTaskDispatcher struct {
	client *asynq.Client
	cfg    *config.Config
}

func NewAsynqTaskDispatcher(cfg *config.Config, client *asynq.Client) *AsynqTaskDispatcher {
	return &AsynqTaskDispatcher{
		client: client,
		cfg:    cfg,
	}
}

var _ gateway.TaskDispatcher = (*AsynqTaskDispatcher)(nil)

func (d *AsynqTaskDispatcher) DispatchSendVerifyEmail(
	ctx context.Context,
	payload *gateway.PayloadSendVerifyEmail,
) error {
	payload.Subject = TaskSendVerifyEmailSubject
	payload.From = d.cfg.AppName

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload)

	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

func (d *AsynqTaskDispatcher) DispatchSendForgotPasswordEmail(
	ctx context.Context,
	payload *gateway.PayloadSendForgotPasswordEmail,
) error {
	payload.Subject = TaskSendForgotPasswordEmailSubject
	payload.From = d.cfg.AppName
	payload.ResetLink = fmt.Sprintf("%s/reset-password?email=%s&code=%s", d.cfg.AppClientURL, payload.Email, payload.VerificationCode)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendForgotPasswordEmail, jsonPayload)

	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}
