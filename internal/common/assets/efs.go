package assets

import "embed"

//go:embed "template"
var EmbeddedFiles embed.FS

const (
	EmailVerificationTemplatePath   = "template/email-verification.tmpl"
	EmailForgotPasswordTemplatePath = "template/email-forgot-password.tmpl"
)
