package infra

import (
	"database/sql"
	"errors"
	"fmt"

	"marketplace/domain/catalog"
	"marketplace/domain/sales"
	"marketplace/domain/shared"

	"github.com/lib/pq"
)

type PostgresRepo struct{ db *sql.DB }

func NewPostgresRepo(db *sql.DB) *PostgresRepo { return &PostgresRepo{db: db} }
func (r *PostgresRepo) Close() error           { return r.db.Close() }

func (r *PostgresRepo) Save(l *sales.Listing) error {
	return r.saveListingTx(r.db, l)
}

func (r *PostgresRepo) SaveWithOutbox(l *sales.Listing, evt sales.OutboxEvent) error {
	tx, err := r.db.Begin()
	if err != nil { return err }
	defer tx.Rollback()

	if err := r.saveListingTx(tx, l); err != nil { return err }

	_, err = tx.Exec(
		`INSERT INTO outbox (id, event_type, payload) VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO NOTHING`,
		evt.ID, evt.Type, evt.Payload,
	)
	if err != nil { return fmt.Errorf("save outbox: %w", err) }

	return tx.Commit()
}

type execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func (r *PostgresRepo) saveListingTx(db execer, l *sales.Listing) error {
	q := `
		INSERT INTO listings (
			id, car_id, car_brand, car_model, car_year, car_mileage,
			seller_id, seller_name, seller_type, title, description, photos,
			price_amount, status, has_warranty, has_test_drive, region,
			commission, delivery_id
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			delivery_id = EXCLUDED.delivery_id,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			photos = EXCLUDED.photos,
			commission = EXCLUDED.commission,
			has_warranty = EXCLUDED.has_warranty,
			has_test_drive = EXCLUDED.has_test_drive,
			region = EXCLUDED.region`
	_, err := db.Exec(q,
		l.ID().String(),
		l.Car().ID().String(), l.Car().Brand(), l.Car().Model(),
		l.Car().Year().Value(), l.Car().Mileage().Value(),
		l.Seller().ID().String(), l.Seller().Name(), string(l.Seller().Type()),
		l.Title(), l.Description(), pq.Array(l.Photos()),
		l.Price().Amount(), string(l.Status()),
		l.HasWarranty(), l.HasTestDrive(), l.Region(),
		l.Commission().Value(), nullableString(l.DeliveryID()),
	)
	return err
}

func (r *PostgresRepo) FindByID(id sales.ListingID) (*sales.Listing, error) {
	q := `
		SELECT id, car_id, car_brand, car_model, car_year, car_mileage,
			seller_id, seller_name, seller_type, title, description, photos,
			price_amount, status, has_warranty, has_test_drive, region,
			commission, COALESCE(delivery_id, '')
		FROM listings WHERE id = $1`
	return scan(r.db.QueryRow(q, id.String()))
}

func (r *PostgresRepo) All() ([]*sales.Listing, error) {
	q := `
		SELECT id, car_id, car_brand, car_model, car_year, car_mileage,
			seller_id, seller_name, seller_type, title, description, photos,
			price_amount, status, has_warranty, has_test_drive, region,
			commission, COALESCE(delivery_id, '')
		FROM listings ORDER BY id`
	rows, err := r.db.Query(q)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*sales.Listing
	for rows.Next() {
		l, err := scan(rows)
		if err != nil { return nil, err }
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *PostgresRepo) AttachDelivery(id sales.ListingID, deliveryID string) error {
	_, err := r.db.Exec(`UPDATE listings SET delivery_id = $1 WHERE id = $2`, deliveryID, id.String())
	return err
}

type scanner interface{ Scan(...any) error }

func scan(s scanner) (*sales.Listing, error) {
	var (
		id, carID, brand, model, sellerID, sellerName, sellerType string
		title, desc, status, region, deliveryID                    string
		carYear, carMileage, commission                            int
		priceAmount                                                int64
		photos                                                     []string
		warranty, testDrive                                        bool
	)
	err := s.Scan(&id, &carID, &brand, &model, &carYear, &carMileage,
		&sellerID, &sellerName, &sellerType, &title, &desc, pq.Array(&photos),
		&priceAmount, &status, &warranty, &testDrive, &region,
		&commission, &deliveryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, errors.New("not found") }
		return nil, err
	}
	cID, _ := catalog.NewCarID(carID)
	y, _ := catalog.NewYear(carYear)
	m, _ := catalog.NewMileage(carMileage)
	p, _ := shared.NewMoney(priceAmount, "RUB")
	car, _ := catalog.NewCar(cID, brand, model, y, m, p)
	sID, _ := sales.NewSellerID(sellerID)
	seller, _ := sales.NewSeller(sID, sellerName, sales.SellerType(sellerType))
	lID, _ := sales.NewListingID(id)
	d, _ := sales.NewDiscountPercent(commission)
	return sales.Restore(lID, car, seller, title, desc, photos, p,
		sales.ListingStatus(status), warranty, testDrive, region, d, deliveryID), nil
}

func nullableString(s string) any {
	if s == "" { return nil }
	return s
}
