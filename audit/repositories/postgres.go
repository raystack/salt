package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/odpf/salt/audit"
)

type AuditPostgresModel struct {
	Timestamp time.Time      `db:"timestamp"`
	Action    string         `db:"action"`
	Actor     string         `db:"actor"`
	Data      types.JSONText `db:"data"`
	Metadata  types.JSONText `db:"metadata"`
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (r *PostgresRepository) DB() *sql.DB {
	return r.db
}

func (r *PostgresRepository) Init(ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS audit_logs (
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			action TEXT NOT NULL,
			actor TEXT NOT NULL,
			data JSONB NOT NULL,
			metadata JSONB NOT NULL
		);

		CREATE INDEX IF NOT EXISTS audit_logs_timestamp_idx ON audit_logs (timestamp);
		CREATE INDEX IF NOT EXISTS audit_logs_action_idx ON audit_logs (action);
		CREATE INDEX IF NOT EXISTS audit_logs_actor_idx ON audit_logs (actor);
	`
	if _, err := r.db.ExecContext(ctx, sql); err != nil {
		return fmt.Errorf("migrating audit model to postgres db: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Insert(ctx context.Context, l *audit.Log) error {
	data, err := json.Marshal(l.Data)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}
	metadata, err := json.Marshal(l.Metadata)
	if err != nil {
		return fmt.Errorf("marshaling metadata: %w", err)
	}
	m := &AuditPostgresModel{
		Timestamp: l.Timestamp,
		Action:    l.Action,
		Actor:     l.Actor,
		Data:      data,
		Metadata:  metadata,
	}

	if _, err := r.db.ExecContext(ctx, "INSERT INTO audit_logs (timestamp, action, actor, data, metadata) VALUES ($1, $2, $3, $4, $5)", m.Timestamp, m.Action, m.Actor, m.Data, m.Metadata); err != nil {
		return fmt.Errorf("inserting to db: %w", err)
	}

	return nil
}
