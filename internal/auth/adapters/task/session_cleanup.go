package task

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// TaskSessionCleanup is the unique identifier for the session cleanup task
const TaskSessionCleanup = "auth:session:cleanup"

// NewSessionCleanupTask creates a new task for session cleanup.
func NewSessionCleanupTask() *asynq.Task {
	return asynq.NewTask(TaskSessionCleanup, nil)
}

// SessionCleanupProcessor handles the execution of session cleanup.
type SessionCleanupProcessor struct {
	sessionRepo session.Repository
	log         logger.Logger
}

// NewSessionCleanupProcessor creates a new processor instance with required dependencies.
func NewSessionCleanupProcessor(
	sessionRepo session.Repository,
	log logger.Logger,
) *SessionCleanupProcessor {
	return &SessionCleanupProcessor{
		sessionRepo: sessionRepo,
		log:         log,
	}
}

// ProcessTask implements the asynq.Handler interface.
func (p *SessionCleanupProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	p.log.Info(ctx, "starting session cleanup processor",
		logger.Field{Key: "task_id", Value: t.ResultWriter().TaskID()},
	)

	deletedCount, err := p.sessionRepo.DeleteExpired(ctx)
	if err != nil {
		p.log.Error(ctx, err, "failed to cleanup expired sessions")
		return err
	}

	if deletedCount > 0 {
		p.log.Info(ctx, "session cleanup completed",
			logger.Field{Key: "deleted_count", Value: deletedCount},
		)
	} else {
		p.log.Debug(ctx, "no expired sessions found")
	}

	return nil
}
