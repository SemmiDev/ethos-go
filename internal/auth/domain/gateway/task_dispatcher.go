package gateway

import (
	"context"

	"github.com/google/uuid"
)

type PayloadSendVerifyEmail struct {
	UserID                     uuid.UUID `json:"user_id"`
	Name                       string    `json:"name"`
	Email                      string    `json:"email"`
	VerificationCode           string    `json:"verification_code"`
	VerificationCodeExpiration int       `json:"verification_code_expiration"` // in minutes

	// fill by dispatcher
	From    string `json:"from"`
	Subject string `json:"subject"`
}

type PayloadSendForgotPasswordEmail struct {
	UserID                     uuid.UUID `json:"user_id"`
	Name                       string    `json:"name"`
	Email                      string    `json:"email"`
	VerificationCode           string    `json:"verification_code"`
	VerificationCodeExpiration int       `json:"verification_code_expiration"` // in minutes

	// fill by dispatcher
	From      string `json:"from"`
	Subject   string `json:"subject"`
	ResetLink string `json:"reset_link"`
}

// TaskDispatcher defines the interface for dispatching background tasks
type TaskDispatcher interface {
	DispatchSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail) error
	DispatchSendForgotPasswordEmail(ctx context.Context, payload *PayloadSendForgotPasswordEmail) error
}
