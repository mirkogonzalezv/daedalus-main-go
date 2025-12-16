package usecases

import (
	"daedalus-engine-go/cmd/core/domain/repository"

	"go.uber.org/zap"
)

type SessionUseCase struct {
	repo repository.SessionRepository
	log  *zap.Logger
}

func NewSessionUseCase(repo repository.SessionRepository, log *zap.Logger) *SessionUseCase {
	return &SessionUseCase{repo: repo, log: log}
}
