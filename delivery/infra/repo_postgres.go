package infra

import (
	"database/sql"
	"errors"

	"delivery/domain"
	"delivery/ports"
)

type PostgresRepo struct{ db *sql.DB }
func NewPostgresRepo(db *sql.DB) *PostgresRepo { return &PostgresRepo{db: db} }
func (r *PostgresRepo) Close() error           { return r.db.Close() }

func (r *PostgresRepo) Save(d *domain.Delivery) error { return saveTx(r.db, d) }

func (r *PostgresRepo) SaveWithOutbox(d *domain.Delivery, eventID, eventType string, payload []byte) error {
	tx, err := r.db.Begin()
	if err != nil { return err }
	defer tx.Rollback()
	if err := saveTx(tx, d); err != nil { return err }
	_, err = tx.Exec(`INSERT INTO outbox (id, event_type, payload) VALUES ($1,$2,$3) ON CONFLICT (id) DO NOTHING`,
		eventID, eventType, payload)
	if err != nil { return err }
	return tx.Commit()
}

type execer interface{ Exec(string, ...any) (sql.Result, error) }

func saveTx(db execer, d *domain.Delivery) error {
	_, err := db.Exec(`INSERT INTO deliveries (id, listing_id, customer_id, type, from_city, to_city, status, price, distance_km)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status`,
		d.ID(), d.ListingID(), d.CustomerID(), string(d.Type()),
		d.From().City, d.To().City, string(d.Status()), d.Price(), d.DistanceKM())
	return err
}

func (r *PostgresRepo) FindByID(id string) (*domain.Delivery, error) {
	q := `SELECT id, listing_id, customer_id, type, from_city, to_city, status, price, distance_km
		FROM deliveries WHERE id = $1`
	return scan(r.db.QueryRow(q, id))
}

func (r *PostgresRepo) All() ([]*domain.Delivery, error) {
	q := `SELECT id, listing_id, customer_id, type, from_city, to_city, status, price, distance_km
		FROM deliveries ORDER BY created_at DESC`
	rows, err := r.db.Query(q)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*domain.Delivery
	for rows.Next() {
		d, err := scan(rows)
		if err != nil { return nil, err }
		out = append(out, d)
	}
	return out, rows.Err()
}

type scanner interface{ Scan(...any) error }

func scan(s scanner) (*domain.Delivery, error) {
	var id, listingID, customerID, typ, fromCity, toCity, status string
	var price, distance int
	err := s.Scan(&id, &listingID, &customerID, &typ, &fromCity, &toCity, &status, &price, &distance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, errors.New("not found") }
		return nil, err
	}
	return domain.Restore(id, listingID, customerID, domain.DeliveryType(typ),
		domain.Point{City: fromCity}, domain.Point{City: toCity},
		domain.DeliveryStatus(status), price, distance), nil
}

type OutboxPostgres struct{ db *sql.DB }
func NewOutboxRepo(db *sql.DB) *OutboxPostgres { return &OutboxPostgres{db: db} }

func (r *OutboxPostgres) FetchUnpublished(limit int) ([]ports.OutboxEvent, error) {
	rows, err := r.db.Query(`SELECT id, event_type, payload FROM outbox WHERE published_at IS NULL ORDER BY created_at LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ports.OutboxEvent
	for rows.Next() {
		var e ports.OutboxEvent
		if err := rows.Scan(&e.ID, &e.Type, &e.Payload); err != nil { return nil, err }
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *OutboxPostgres) MarkPublished(id string) error {
	_, err := r.db.Exec(`UPDATE outbox SET published_at = NOW() WHERE id = $1`, id)
	return err
}
