package infra

import (
	"database/sql"

	"marketplace/domain/sales"
)

type OutboxRepo struct{ db *sql.DB }

func NewOutboxRepo(db *sql.DB) *OutboxRepo { return &OutboxRepo{db: db} }

func (r *OutboxRepo) FetchUnpublished(limit int) ([]sales.OutboxEvent, error) {
	rows, err := r.db.Query(
		`SELECT id, event_type, payload FROM outbox
		 WHERE published_at IS NULL ORDER BY created_at LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()

	var out []sales.OutboxEvent
	for rows.Next() {
		var e sales.OutboxEvent
		if err := rows.Scan(&e.ID, &e.Type, &e.Payload); err != nil { return nil, err }
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *OutboxRepo) MarkPublished(id string) error {
	_, err := r.db.Exec(`UPDATE outbox SET published_at = NOW() WHERE id = $1`, id)
	return err
}
