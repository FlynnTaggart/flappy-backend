package service

import (
	"context"

	"flappy-backend/internal/domain"
)

type DB interface {
	AddUserRecord(ctx context.Context, user domain.Record) error
	GetUserRecord(ctx context.Context, nickname string) (int64, error)
	GetTop10UserRecords(ctx context.Context) ([]*domain.Record, error)
}
