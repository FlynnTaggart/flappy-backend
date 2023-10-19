package handler

import (
	"context"

	"flappy-backend/internal/domain"
)

type Service interface {
	AddRecord(ctx context.Context, user domain.Record) error
	GetRecord(ctx context.Context, nickname string) (int64, error)
	GetTop10Records(ctx context.Context) ([]*domain.Record, error)
}

type Logger interface {
	Debug(msg string, data map[string]interface{})
	Info(msg string, data map[string]interface{})
	Warn(msg string, data map[string]interface{})
	Error(msg string, data map[string]interface{})
}

type addRecordIn struct {
	Nickname string `json:"nickname"`
	Score    int64  `json:"score"`
}

type getRecordIn struct {
	Nickname string `json:"nickname"`
}

type getRecordOut struct {
	Score int64 `json:"score"`
}

type getTop10RecordsOut struct {
	Records []record `json:"records"`
}

type record struct {
	Nickname string `json:"nickname"`
	Score    int64  `json:"score"`
}
