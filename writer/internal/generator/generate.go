package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"writer/internal/models"
)

// Функция для генерации заказа
func GenerateOrder() ([]byte, error) {
	order := models.Order{
		OrderUID:          generateOrderUID(),
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Delivery:          models.Delivery{Name: "Test Testov", Phone: "+9720000000", Zip: "2639809", City: "Kiryat Mozkin", Address: "Ploshad Mira 15", Region: "Kraiot", Email: "test@gmail.com"},
		Payment:           models.Payment{Transaction: generateOrderUID(), RequestID: "", Currency: "USD", Provider: "wbpay", Amount: 1817, PaymentDT: 1637907727, Bank: "alpha", DeliveryCost: 1500, GoodsTotal: 317, CustomFee: 0},
		Items:             []models.Item{models.Item{ChrtID: 9934930, TrackNumber: "WBILMTESTTRACK", Price: 453, Rid: "ab4219087a764ae0btest", Name: "Mascaras", Sale: 30, Size: "0", TotalPrice: 317, NmID: 2389212, Brand: "Vivienne Sabo", Status: 202}},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          "1",
	}

	// Преобразуем заказ в JSON
	jsonData, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return jsonData, nil
}

// Генерация уникального идентификатора для заказа
func generateOrderUID() string {
	return fmt.Sprintf("order-%d", rand.Int())
}

//mockgen -source=internal/transport/generateOrder/generate.go -destination=internal/transport/generateOrder/mocks/mock_generate.go -package=mocks
