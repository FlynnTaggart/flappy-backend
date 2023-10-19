package service

import (
	"context"

	"flappy-backend/internal/domain"
)

type FlappyService struct {
	DB DB
}

func NewFlappyService(DB DB) *FlappyService {
	return &FlappyService{
		DB: DB,
	}
}

func (s *FlappyService) AddRecord(ctx context.Context, user domain.Record) error {
	return s.DB.AddUserRecord(ctx, user)
}

func (s *FlappyService) GetRecord(ctx context.Context, nickname string) (int64, error) {
	return s.DB.GetUserRecord(ctx, nickname)
}

func (s *FlappyService) GetTop10Records(ctx context.Context) ([]*domain.Record, error) {
	return s.DB.GetTop10UserRecords(ctx)
}
