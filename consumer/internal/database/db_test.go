package database

// import (
// 	//"database/sql"
// 	//"strconv"

// 	"testing"

// 	"consumer/internal/config"
// 	"consumer/internal/models" // Adjust the import path as necessary

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/lib/pq"
// )

// func TestAddOrder(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("failed to create mock db: %v", err)
// 	}
// 	defer db.Close()

// 	postgres := &Postgres{Connection: db}

// 	// Define a sample order based on the provided structure
// 	order := models.Order{
// 		OrderUID:    "b563feb7b2b84b6test",
// 		TrackNumber: "WBILMTESTTRACK",
// 		Entry:       "WBIL",
// 		Delivery: models.Delivery{
// 			Name:    "Test Testov",
// 			Phone:   "+9720000000",
// 			Zip:     "2639809",
// 			City:    "Kiryat Mozkin",
// 			Address: "Ploshad Mira 15",
// 			Region:  "Kraiot",
// 			Email:   "test@gmail.com",
// 		},
// 		Payment: models.Payment{
// 			Transaction:  "b563feb7b2b84b6test",
// 			RequestID:    "",
// 			Currency:     "USD",
// 			Provider:     "wbpay",
// 			Amount:       1817,
// 			PaymentDT:    1637907727,
// 			Bank:         "alpha",
// 			DeliveryCost: 1500,
// 			GoodsTotal:   317,
// 			CustomFee:    0,
// 		},
// 		Items: []models.Item{
// 			{
// 				ChrtID:      9934930,
// 				TrackNumber: "WBILMTESTTRACK",
// 				Price:       453,
// 				Rid:         "ab4219087a764ae0btest",
// 				Name:        "Mascaras",
// 				Sale:        30,
// 				Size:        "0",
// 				TotalPrice:  317,
// 				NmID:        2389212,
// 				Brand:       "Vivienne Sabo",
// 				Status:      202,
// 			},
// 		},
// 		Locale:            "en",
// 		InternalSignature: "",
// 		CustomerID:        "test",
// 		DeliveryService:   "meest",
// 		ShardKey:          "9",
// 		SmID:              99,
// 		DateCreated:       "2021-11-26T06:22:19Z", // Use a Unix timestamp
// 		OofShard:          "1",
// 	}

// 	deliveryID := uuid.New()

// 	queryAddDelivery := `WITH insert_return AS (
//             INSERT INTO  delivery (name, phone, zip, city, address, region,	email)
//             VALUES ($1, $2, $3, $4, $5, $6, $7)
//             RETURNING id
//         )
//         SELECT id FROM insert_return`

// 	// Начало транзакции
// 	mock.ExpectBegin()

// 	// Ожидание запроса на вставку доставки
// 	mock.ExpectQuery(queryAddDelivery).
// 		WithArgs(order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(deliveryID.String()))

// 	mock.ExpectQuery(`WITH insert_return AS \(INSERT INTO payment \(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee\) VALUES \(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7, \\$8, \\$9, \\$10\) RETURNING id\)`).
// 		WithArgs(order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("payment-id"))

// 	mock.ExpectQuery(`WITH insert_return AS \(INSERT INTO items \(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status\) VALUES \(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7, \\$8, \\$9, \\$10, \\$11\) RETURNING id\)`).
// 		WithArgs(order.Items[0].ChrtID, order.Items[0].TrackNumber, order.Items[0].Price, order.Items[0].Rid, order.Items[0].Name, order.Items[0].Sale, order.Items[0].Size, order.Items[0].TotalPrice, order.Items[0].NmID, order.Items[0].Brand, order.Items[0].Status).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("item-id"))

// 	mock.ExpectExec(`INSERT INTO orders \(order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard\)`).
// 		WithArgs(order.OrderUID, order.TrackNumber, order.Entry, "delivery-id", "payment-id", pq.Array([]string{"item-id"}), order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard).
// 		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulating successful insert

// 	// Call the AddOrder function
// 	orderUID, err := postgres.AddOrder(order)

// 	// Assert the expectations
// 	assert.NoError(t, err)
// 	fmt.Println(order.OrderUID, " ", orderUID)
// 	assert.Equal(t, order.OrderUID, orderUID)

// 	// Ensure all expectations were met
// 	err = mock.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }

// func TestAddOrder(t *testing.T) {
// 	// Create a mock database connection
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	// Create a new Postgres instance with the mock connection
// 	pg := NewPostgres(config.Config{
// 		DBUser:     "user",
// 		DBPassword: "password",
// 		DBName:     "dbname",
// 		DBHost:     "localhost",
// 		DBPort:     "5432",
// 		SslMode:    "disable",
// 		KafkaHost:  "",
// 		KafkaPort:  "",
// 		KafkaTopic: "",
// 	})
// 	pg.Connection = db

// 	// Create a test order
// 	order := models.Order{
// 		OrderUID:    "test-order-uid",
// 		TrackNumber: "test-track-number",
// 		Entry:       "test-entry",
// 		Delivery: models.Delivery{
// 			Name:    "test-name",
// 			Phone:   "test-phone",
// 			Zip:     "test-zip",
// 			City:    "test-city",
// 			Address: "test-address",
// 			Region:  "test-region",
// 			Email:   "test-email",
// 		},
// 		Payment: models.Payment{
// 			Transaction:  "test-transaction",
// 			RequestID:    "test-request-id",
// 			Currency:     "test-currency",
// 			Provider:     "test-provider",
// 			Amount:       100,
// 			PaymentDT:    1643723400,
// 			Bank:         "test-bank",
// 			DeliveryCost: 50,
// 			GoodsTotal:   150,
// 			CustomFee:    20,
// 		},
// 		Items: []models.Item{
// 			{
// 				ChrtID:      1,
// 				TrackNumber: "test-track-number-1",
// 				Price:       100,
// 				Rid:         "test-rid-1",
// 				Name:        "test-name-1",
// 				Sale:        10,
// 				Size:        "test-size-1",
// 				TotalPrice:  110,
// 				NmID:        1,
// 				Brand:       "test-brand-1",
// 				Status:      1,
// 			},
// 			{
// 				ChrtID:      2,
// 				TrackNumber: "test-track-number-2",
// 				Price:       200,
// 				Rid:         "test-rid-2",
// 				Name:        "test-name-2",
// 				Sale:        20,
// 				Size:        "test-size-2",
// 				TotalPrice:  220,
// 				NmID:        2,
// 				Brand:       "test-brand-2",
// 				Status:      2,
// 			},
// 		},
// 		Locale:            "en",
// 		InternalSignature: "test-internal-signature",
// 		CustomerID:        "test-customer-id",
// 		DeliveryService:   "test-delivery-service",
// 		ShardKey:          "test-shard-key",
// 		SmID:              1,
// 		DateCreated:       "2022-02-01 12:00:00",
// 		OofShard:          "test-oof-shard",
// 	}

// 	deliveryID := "550e8400-e29b-41d4-a716-446655440000" // Example UUID
// 	mock.ExpectQuery(`WITH insert_return AS (
//             INSERT INTO delivery (name, phone, zip, city, address, region, email)
//             VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7)
//             RETURNING id
//         )
//         SELECT id FROM insert_return`).
// 		WithArgs(order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(deliveryID))

// 	mock.ExpectExec("INSERT INTO payment").WithArgs(order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee).WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectExec("INSERT INTO items").WithArgs(order.Items[0].ChrtID, order.Items[0].TrackNumber, order.Items[0].Price, order.Items[0].Rid, order.Items[0].Name, order.Items[0].Sale, order.Items[0].Size, order.Items[0].TotalPrice, order.Items[0].NmID, order.Items[0].Brand, order.Items[0].Status).WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectExec("INSERT INTO items").WithArgs(order.Items[1].ChrtID, order.Items[1].TrackNumber, order.Items[1].Price, order.Items[1].Rid, order.Items[1].Name, order.Items[1].Sale, order.Items[1].Size, order.Items[1].TotalPrice, order.Items[1].NmID, order.Items[1].Brand, order.Items[1].Status).WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectExec("INSERT INTO orders").WithArgs(order.OrderUID, order.TrackNumber, order.Entry, "delivery-id", "payment-id", pq.Array{"item-id-1", "item-id-2"}, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard).WillReturnResult(sqlmock.NewResult(1, 1))

// 	// Call the AddOrder function
// 	id, err := pg.AddOrder(order)
// 	if err != nil {
// 		t.Errorf("expected no error, but got %v", err)
// 	}
// 	if id != order.OrderUID {
// 		t.Errorf("expected order UID to be %s, but got %s", order.OrderUID, id)
// 	}

// 	// Assert that the expected queries were executed
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
