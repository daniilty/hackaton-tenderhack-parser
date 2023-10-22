package logs

import (
	"context"
	"tenderhack-parser/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type Logs struct {
	db *sqlx.DB
}

func NewLogs(db *sqlx.DB) *Logs {
	return &Logs{db}
}

func (l *Logs) UpsertLog(ctx context.Context, log *models.Log) error {
	const q = `insert into logs_data(id, ts, category_id, data) values(:id, :ts, :category_id, :data)`

	_, err := l.db.NamedExecContext(ctx, q, log)

	return err
}

func (l *Logs) Resolve(ctx context.Context, id int) error {
	const q = `update logs_data set resolved_at=now() where id=$1`

	_, err := l.db.ExecContext(ctx, q, id)

	return err
}

func (l *Logs) Unresolve(ctx context.Context, id int) error {
	const q = `update logs_data set resolved_at=null where id=$1`

	_, err := l.db.ExecContext(ctx, q, id)

	return err
}

func (l *Logs) GetResolved(ctx context.Context, from time.Time, to time.Time) ([]*models.Log, error) {
	const q = `select id, ts, category_id, data, resolved_at from logs_data where ts >= $1 AND ts <= $2 AND resolved_at is not null ORDER BY ts asc`

	logs := []*models.Log{}
	err := l.db.SelectContext(ctx, &logs, q, from, to)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (l *Logs) GetUnresolved(ctx context.Context, from time.Time, to time.Time) ([]*models.Log, error) {
	const q = `select id, ts, category_id, data, resolved_at from logs_data where ts >= $1 AND ts <= $2 AND resolved_at is null ORDER BY ts asc`

	logs := []*models.Log{}
	err := l.db.SelectContext(ctx, &logs, q, from, to)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
