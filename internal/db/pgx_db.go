package db

import (
	"context"
	"fmt"

	"flappy-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxDB struct {
	*pgxpool.Pool
}

func NewPgxDB(pool *pgxpool.Pool) *PgxDB {
	return &PgxDB{
		Pool: pool,
	}
}

func (p *PgxDB) AddUserRecord(ctx context.Context, record domain.Record) error {
	_, err := p.Exec(ctx, `insert into records as r (nickname, score) values ($1, $2)
		on conflict (nickname) do 
		update set score = EXCLUDED.score
		where r.score < EXCLUDED.score`,
		record.Nickname,
		record.Score)

	if err != nil {
		return fmt.Errorf("db: AddUserRecord: error query: %w", err)
	}

	return nil
}

func (p *PgxDB) GetUserRecord(ctx context.Context, nickname string) (int64, error) {
	record := &domain.Record{}
	err := p.QueryRow(ctx, `select score from records where nickname = $1`, nickname).
		Scan(&record.Score)
	if err != nil {
		return 0, fmt.Errorf("db: GetUserRecord: error query: %w", err)
	}

	return record.Score, nil
}

func (p *PgxDB) GetTop10UserRecords(ctx context.Context) ([]*domain.Record, error) {
	rows, err := p.Query(ctx, `select * from records order by score desc limit 10`)
	if err != nil {
		return nil, fmt.Errorf("db: GetTop10UserRecords: error query: %w", err)
	}

	res := make([]*domain.Record, 0, 10)

	for rows.Next() {
		record := &domain.Record{}

		err = rows.Scan(&record.Nickname, &record.Score)

		if err != nil {
			return nil, fmt.Errorf("db: GetTop10UserRecords: error parsing query: %w", err)
		}

		res = append(res, record)

		if len(res) == 10 {
			break
		}
	}

	return res, nil
}
