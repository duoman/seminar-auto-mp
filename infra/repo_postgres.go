package infra

import (
	"database/sql"
	"errors"
	"fmt"

	"marketplace/domain"

	"github.com/lib/pq"
)

type PostgresListingRepository struct {
	db *sql.DB
}

func NewPostgresListingRepository(db *sql.DB) *PostgresListingRepository {
	return &PostgresListingRepository{db: db}
}

func (r *PostgresListingRepository) Close() error { return r.db.Close() }

func (r *PostgresListingRepository) Save(l *domain.Listing) error {
	if l == nil {
		return errors.New("nil listing")
	}
	if l.Car == nil || l.Seller == nil {
		return errors.New("listing must have car and seller")
	}

	query := `
		INSERT INTO listings (
			id, car_id, car_brand, car_model, car_year, car_mileage,
			seller_id, seller_name, seller_type,
			title, description, photos, price, status,
			has_warranty, has_test_drive, region, commission
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18
		)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			photos = EXCLUDED.photos,
			price = EXCLUDED.price,
			status = EXCLUDED.status,
			has_warranty = EXCLUDED.has_warranty,
			has_test_drive = EXCLUDED.has_test_drive,
			region = EXCLUDED.region,
			commission = EXCLUDED.commission
	`
	_, err := r.db.Exec(query,
		l.ID,
		l.Car.ID, l.Car.Brand, l.Car.Model, l.Car.Year, l.Car.Mileage,
		l.Seller.ID, l.Seller.Name, string(l.Seller.Type),
		l.Title, l.Description, pq.Array(l.Photos), l.Price, string(l.Status),
		l.HasWarranty, l.HasTestDrive, l.Region, l.Commission,
	)
	if err != nil {
		return fmt.Errorf("save listing: %w", err)
	}
	return nil
}

func (r *PostgresListingRepository) FindByID(id string) (*domain.Listing, error) {
	query := `
		SELECT
			id, car_id, car_brand, car_model, car_year, car_mileage,
			seller_id, seller_name, seller_type,
			title, description, photos, price, status,
			has_warranty, has_test_drive, region, commission
		FROM listings WHERE id = $1
	`
	return scanListing(r.db.QueryRow(query, id))
}

func (r *PostgresListingRepository) All() ([]*domain.Listing, error) {
	query := `
		SELECT
			id, car_id, car_brand, car_model, car_year, car_mileage,
			seller_id, seller_name, seller_type,
			title, description, photos, price, status,
			has_warranty, has_test_drive, region, commission
		FROM listings ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}
	defer rows.Close()

	var out []*domain.Listing
	for rows.Next() {
		l, err := scanListing(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

type rowScanner interface{ Scan(dest ...any) error }

func scanListing(s rowScanner) (*domain.Listing, error) {
	var (
		l          domain.Listing
		car        domain.Car
		seller     domain.Seller
		photos     []string
		sellerType string
		status     string
	)
	err := s.Scan(
		&l.ID,
		&car.ID, &car.Brand, &car.Model, &car.Year, &car.Mileage,
		&seller.ID, &seller.Name, &sellerType,
		&l.Title, &l.Description, pq.Array(&photos), &l.Price, &status,
		&l.HasWarranty, &l.HasTestDrive, &l.Region, &l.Commission,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("listing not found")
		}
		return nil, fmt.Errorf("scan listing: %w", err)
	}
	l.Car = &car
	l.Seller = &seller
	seller.Type = domain.SellerType(sellerType)
	l.Status = domain.ListingStatus(status)
	l.Photos = photos
	return &l, nil
}
