package main

import (
	"fmt"

	"marketplace-anemic/domain"
	"marketplace-anemic/service"
)

func main() {
	fmt.Println("=== Анемичная модель: инварианты не защищены ===")
	fmt.Println()

	svc := service.New()

	// Проблема 1. Отрицательная цена — возможно
	l1 := svc.Create("L-1", "CAR-1", "Tesla", -1000)
	fmt.Printf("❌ Создано объявление с ценой %d\n", l1.Price)
	fmt.Println("   Money(-1000) — можно. Никто не проверит.")
	fmt.Println()

	// Проблема 2. Публичные поля — Status можно поставить вручную
	l2 := svc.Create("L-2", "CAR-2", "BMW", 5_000_000)
	l2.Status = "active" // ← в обход Publish
	fmt.Printf("❌ Статус выставлен напрямую: %s (без publish)\n", l2.Status)
	fmt.Println("   Поле публичное — любой код может это сделать.")
	fmt.Println()

	// Проблема 3. Продажа без публикации
	l3 := svc.Create("L-3", "CAR-3", "Audi", 3_000_000)
	l3.Status = "active" // ← «читерство»
	_ = svc.Sell(l3)     // продали
	_ = svc.Sell(l3)     // продали ещё раз? нет: проверка есть, но
	// если бы кто-то напрямую: l3.Status = "sold" — можно было бы бесконечно
	fmt.Printf("❌ Sold listing: %s\n", l3.Status)
	fmt.Println()

	// Проблема 4. 100500 фото — возможно
	l4 := svc.Create("L-4", "CAR-4", "Kia", 1_500_000)
	for i := 0; i < 100; i++ {
		l4.Photos = append(l4.Photos, fmt.Sprintf("photo-%d", i))
	}
	fmt.Printf("❌ Фото в объявлении: %d (лимит 10 — где-то в голове)\n", len(l4.Photos))
	fmt.Println()

	// Проблема 5. Комиссия 150% — возможно
	l5 := svc.Create("L-5", "CAR-5", "Lada", 800_000)
	l5.Commission = 150
	fmt.Printf("❌ Комиссия: %d%% (должна быть 0-100)\n", l5.Commission)
	fmt.Println()

	// Проблема 6. Публикация без фото — «прошла», если не через svc
	l6 := svc.Create("L-6", "CAR-6", "Ford", 2_000_000)
	l6.Status = "pending"
	// прошли мимо svc.Publish → фото нет, но статус active
	l6.Status = "active"
	fmt.Printf("❌ Публикация без фото (обход svc.Publish): status=%s, photos=%d\n",
		l6.Status, len(l6.Photos))
	fmt.Println()

	// Проблема 7. Тест требует сервис, потому что логика в нём
	fmt.Println("❌ Тест 'нельзя опубликовать без фото' требует:")
	fmt.Println("   - создать Listing через svc.Create")
	fmt.Println("   - вызвать svc.Publish")
	fmt.Println("   - если svc.Publish забыли — тест не поймает (логика в сервисе)")
	fmt.Println()

	// А вот что нельзя:
	_ = domain.Listing{} // ← можно, пустой Listing с полями по умолчанию

	fmt.Println("=== Итог ===")
	fmt.Println("Все правила — в сервисе. Поля публичные.")
	fmt.Println("Любой код может нарушить инварианты напрямую.")
	fmt.Println()
	fmt.Println("Сравни с marketplace/domain/sales/listing.go:")
	fmt.Println("  - поля приватные")
	fmt.Println("  - только конструктор NewListing")
	fmt.Println("  - только методы Publish(), MarkSold(), AttachDelivery()")
	fmt.Println("  - Money, DiscountPercent, ListingID — валидируют себя")
}
