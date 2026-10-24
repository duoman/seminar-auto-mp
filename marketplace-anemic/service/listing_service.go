package service

import (
	"errors"

	"marketplace-anemic/domain"
)

// Вся бизнес-логика в сервисе. Каждая функция — «набор правил»,
// которые легко забыть в другом сервисе.
type ListingService struct{}

func New() *ListingService { return &ListingService{} }

// Create — создаёт объявление. Валидации нет.
func (s *ListingService) Create(id, carID, title string, price int) *domain.Listing {
	return &domain.Listing{
		ID:     id,
		CarID:  carID,
		Title:  title,
		Price:  price,
		Status: "draft",
	}
}

// Publish — «бизнес-логика». Проверки есть, но их можно обойти,
// если другой код меняет поле Status напрямую.
func (s *ListingService) Publish(l *domain.Listing) error {
	if l.Status != "pending" {
		return errors.New("только из pending можно опубликовать")
	}
	if l.Title == "" {
		return errors.New("нужен заголовок")
	}
	if len(l.Photos) == 0 {
		return errors.New("нужно хотя бы одно фото")
	}
	l.Status = "active"
	return nil
}

// Sell — «продать». Правило — только из active.
func (s *ListingService) Sell(l *domain.Listing) error {
	if l.Status != "active" {
		return errors.New("продать можно только active")
	}
	l.Status = "sold"
	return nil
}
