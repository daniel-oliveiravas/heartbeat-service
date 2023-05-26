package heartbeat

import (
	"context"

	"go.uber.org/zap"
)

type Repository interface {
	Upsert(ctx context.Context, heartbeat Heartbeat) error
}

type EventPublisher interface {
	Publish(ctx context.Context, heartbeat Heartbeat) error
}

type Usecase struct {
	Logger         *zap.Logger
	Repository     Repository
	EventPublisher EventPublisher
}

func NewUsecase(logger *zap.Logger, repository Repository, publisher EventPublisher) *Usecase {
	return &Usecase{
		Logger:         logger,
		Repository:     repository,
		EventPublisher: publisher,
	}
}

func (u *Usecase) Beat(ctx context.Context, heartbeat Heartbeat) error {
	if err := u.Repository.Upsert(ctx, heartbeat); err != nil {
		return err
	}

	if err := u.EventPublisher.Publish(ctx, heartbeat); err != nil {
		return err
	}

	return nil
}
