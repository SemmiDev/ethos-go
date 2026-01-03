package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/auth/adapters/google"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

type GetGoogleAuthURLQuery struct {
	State string
}

type GetGoogleAuthURLHandler decorator.QueryHandler[GetGoogleAuthURLQuery, string]

type getGoogleAuthURLHandler struct {
	googleService *google.Service
}

func NewGetGoogleAuthURLHandler(
	googleService *google.Service,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetGoogleAuthURLHandler {
	return decorator.ApplyQueryDecorators[GetGoogleAuthURLQuery, string](
		getGoogleAuthURLHandler{googleService: googleService},
		log,
		metricsClient,
	)
}

func (h getGoogleAuthURLHandler) Handle(ctx context.Context, query GetGoogleAuthURLQuery) (string, error) {
	return h.googleService.GetLoginURL(query.State), nil
}
