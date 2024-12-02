package transport

import (
	"testing"
	"writer/internal/models"

	publ "writer/internal/publisher/mocks"

	"github.com/fasthttp/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestSetHandler_Set(t *testing.T) {
	type mockBehavior func(r *publ.MockInterfaceKafkaClient, ord models.Order)
	testTable := []struct {
		name                 string       // имя теста
		inputBody            string       //тело запроса
		inputOrder           models.Order //структура заказа, которую мы пеердаем в метод сервиса
		mockBehavior         mockBehavior
		expectedStatusCode   int    // ожидаемый статус код
		expectedResponseBody string //тело ответа
	}{
		{ // тест на успешный вызов SendOrderToKafka
			name: "Ok",
			inputBody: `
					{
					"order_uid": "b563feb7b2b84b6test",
					"track_number": "WBILMTESTTRACK",
					"entry": "WBIL",
					"delivery": {
						"name": "Test Testov",
						"phone": "+9720000000",
						"zip": "2639809",
						"city": "Kiryat Mozkin",
						"address": "Ploshad Mira 15",
						"region": "Kraiot",
						"email": "test@gmail.com"
					},
					"payment": {
						"transaction": "b563feb7b2b84b6test",
						"request_id": "",
						"currency": "USD",
						"provider": "wbpay",
						"amount": 1817,
						"payment_dt": 1637907727,
						"bank": "alpha",
						"delivery_cost": 1500,
						"goods_total": 317,
						"custom_fee": 0
					},
					"items": [
						{
						"chrt_id": 9934930,
						"track_number": "WBILMTESTTRACK",
						"price": 453,
						"rid": "ab4219087a764ae0btest",
						"name": "Mascaras",
						"sale": 30,
						"size": "0",
						"total_price": 317,
						"nm_id": 2389212,
						"brand": "Vivienne Sabo",
						"status": 202
						}
					],
					"locale": "en",
					"internal_signature": "",
					"customer_id": "test",
					"delivery_service": "meest",
					"shardkey": "9",
					"sm_id": 99,
					"date_created": "2021-11-26T06:22:19Z",
					"oof_shard": "1"
					}`,
			inputOrder: models.Order{
				OrderUID:    "b563feb7b2b84b6test",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
				Delivery: models.Delivery{
					Name:    "Test Testov",
					Phone:   "+9720000000",
					Zip:     "2639809",
					City:    "Kiryat Mozkin",
					Address: "Ploshad Mira 15",
					Region:  "Kraiot",
					Email:   "test@gmail.com",
				},
				Payment: models.Payment{
					Transaction:  "b563feb7b2b84b6test",
					RequestID:    "",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       1817,
					PaymentDT:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
					CustomFee:    0,
				},
				Items: []models.Item{
					{
						ChrtID:      9934930,
						TrackNumber: "WBILMTESTTRACK",
						Price:       453,
						Rid:         "ab4219087a764ae0btest",
						Name:        "Mascaras",
						Sale:        30,
						Size:        "0",
						TotalPrice:  317,
						NmID:        2389212,
						Brand:       "Vivienne Sabo",
						Status:      202,
					},
				},
				Locale:            "en",
				InternalSignature: "",
				CustomerID:        "test",
				DeliveryService:   "meest",
				ShardKey:          "9",
				SmID:              99,
				DateCreated:       "2021-11-26T06:22:19Z",
				OofShard:          "1",
			},
			mockBehavior: func(r *publ.MockInterfaceKafkaClient, ord models.Order) {
				r.EXPECT().SendOrderToKafka("records", ord).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		// добавить тест, когда сенд ордер возвращает ошибку
		// {
		// 	name:           "SendOrderToKafka error",
		// 	queryArgs:       "1",
		// 	inputQuantity: 1,
		// 	expectedStatusCode: fasthttp.StatusBadRequest,
		// 	mockBehavior: func(r *mocks.MockInterfaceKafkaClient, ord models.Order, q int) {
		// 		r.EXPECT().SendOrderToKafka("records", ord).Times(q).Return(errors.New("kafka error"))
		// 	},
		// },
		{
			name: "methodnotallowed",
			inputBody: `
				{
					"order_uid": "b563feb7b2b84b6test",
					"track_number": "WBILMTESTTRACK",
					"entry": "WBIL",
					"delivery": {
						"name": "Test Testov",
						"phone": "+9720000000",
						"zip": "2639809",
						"city": "Kiryat Mozkin",
						"address": "Ploshad Mira 15",
						"region": "Kraiot",
						"email": "test@gmail.com"
					},
					"payment": {
						"transaction": "b563feb7b2b84b6test",
						"request_id": "",
						"currency": "USD",
						"provider": "wbpay",
						"amount": 1817,
						"payment_dt": 1637907727,
						"bank": "alpha",
						"delivery_cost": 1500,
						"goods_total": 317,
						"custom_fee": 0
					},
					"items": [
						{
						"chrt_id": 9934930,
						"track_number": "WBILMTESTTRACK",
						"price": 453,
						"rid": "ab4219087a764ae0btest",
						"name": "Mascaras",
						"sale": 30,
						"size": "0",
						"total_price": 317,
						"nm_id": 2389212,
						"brand": "Vivienne Sabo",
						"status": 202
						}
					],
					"locale": "en",
					"internal_signature": "",
					"customer_id": "test",
					"delivery_service": "meest",
					"shardkey": "9",
					"sm_id": 99,
					"date_created": "2021-11-26T06:22:19Z",
					"oof_shard": "1"
				}`,
			inputOrder: models.Order{
				OrderUID:    "b563feb7b2b84b6test",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
				Delivery: models.Delivery{
					Name:    "Test Testov",
					Phone:   "+9720000000",
					Zip:     "2639809",
					City:    "Kiryat Mozkin",
					Address: "Ploshad Mira 15",
					Region:  "Kraiot",
					Email:   "test@gmail.com",
				},
				Payment: models.Payment{
					Transaction:  "b563feb7b2b84b6test",
					RequestID:    "",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       1817,
					PaymentDT:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
					CustomFee:    0,
				},
				Items: []models.Item{
					{
						ChrtID:      9934930,
						TrackNumber: "WBILMTESTTRACK",
						Price:       453,
						Rid:         "ab4219087a764ae0btest",
						Name:        "Mascaras",
						Sale:        30,
						Size:        "0",
						TotalPrice:  317,
						NmID:        2389212,
						Brand:       "Vivienne Sabo",
						Status:      202,
					},
				},
				Locale:            "en",
				InternalSignature: "",
				CustomerID:        "test",
				DeliveryService:   "meest",
				ShardKey:          "9",
				SmID:              99,
				DateCreated:       "2021-11-26T06:22:19Z",
				OofShard:          "1",
			},
			mockBehavior: func(r *publ.MockInterfaceKafkaClient, ord models.Order) {
				//r.EXPECT().SendOrderToKafka("records", ord).Return(nil).Times(1)
			},
			expectedStatusCode:   405,
			expectedResponseBody: "",
		},
		{
			name: "equeljson",
			inputBody: `
				{
					"order_uid": "",
					"track_number": "WBILMTESTTRACK",
					"entry": "WBIL",
					"delivery": {
						"name": "Test Testov",
						"phone": "+9720000000",
						"zip": "2639809",
						"city": "Kiryat Mozkin",
						"address": "Ploshad Mira 15",
						"region": "Kraiot",
						"email": "test@gmail.com"
					},
					"payment": {
						"transaction": "b563feb7b2b84b6test",
						"request_id": "",
						"currency": "USD",
						"provider": "wbpay",
						"amount": 1817,
						"payment_dt": 1637907727,
						"bank": "alpha",
						"delivery_cost": 1500,
						"goods_total": 317,
						"custom_fee": 0
					},
					"items": [
						{
						"chrt_id": 9934930,
						"track_number": "WBILMTESTTRACK",
						"price": 453,
						"rid": "ab4219087a764ae0btest",
						"name": "Mascaras",
						"sale": 30,
						"size": "0",
						"total_price": 317,
						"nm_id": 2389212,
						"brand": "Vivienne Sabo",
						"status": 202
						}
					],
					"locale": "en",
					"internal_signature": "",
					"customer_id": "test",
					"delivery_service": "meest",
					"shardkey": "9",
					"sm_id": 99,
					"date_created": "2021-11-26T06:22:19Z",
					"oof_shard": "1"
				}`,
			inputOrder: models.Order{
				OrderUID:    "",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
				Delivery: models.Delivery{
					Name:    "Test Testov",
					Phone:   "+9720000000",
					Zip:     "2639809",
					City:    "Kiryat Mozkin",
					Address: "Ploshad Mira 15",
					Region:  "Kraiot",
					Email:   "test@gmail.com",
				},
				Payment: models.Payment{
					Transaction:  "b563feb7b2b84b6test",
					RequestID:    "",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       1817,
					PaymentDT:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
					CustomFee:    0,
				},
				Items: []models.Item{
					{
						ChrtID:      9934930,
						TrackNumber: "WBILMTESTTRACK",
						Price:       453,
						Rid:         "ab4219087a764ae0btest",
						Name:        "Mascaras",
						Sale:        30,
						Size:        "0",
						TotalPrice:  317,
						NmID:        2389212,
						Brand:       "Vivienne Sabo",
						Status:      202,
					},
				},
				Locale:            "en",
				InternalSignature: "",
				CustomerID:        "test",
				DeliveryService:   "meest",
				ShardKey:          "9",
				SmID:              99,
				DateCreated:       "2021-11-26T06:22:19Z",
				OofShard:          "1",
			},
			mockBehavior: func(r *publ.MockInterfaceKafkaClient, ord models.Order) {
				//r.EXPECT().SendOrderToKafka("records", ord).Return(nil).Times(1)
			},
			expectedStatusCode:   400,
			expectedResponseBody: "",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			switch testCase.name {
			case "methodnotallowed":
				controller := gomock.NewController(t)
				defer controller.Finish()
				serv := publ.NewMockInterfaceKafkaClient(controller)
				// Создаем новый обработчик
				hb := HandlersBuilder{
					pub:  serv, // Подставьте свою реализацию или оставьте nil для теста
					rout: router.New(),
				}
				hb.rout.POST("/WB/set", hb.Set()) // Регистрируем обработчик для POST метода

				// Создаем запрос с методом GET
				req := fasthttp.AcquireRequest()
				req.Header.SetMethod("GET") // Устанавливаем метод GET
				req.SetRequestURI("/WB/set")
				// Создаем контекст для ответа
				resp := fasthttp.AcquireResponse()
				// Создаем контекст запроса
				ctx := &fasthttp.RequestCtx{}
				ctx.Request.SetRequestURI("/WB/set")
				ctx.Request.Header.SetMethod("GET")
				// Вызываем обработчик маршрутизатора
				hb.rout.Handler(ctx)
				// Проверяем статус ответа
				assert.Equal(t, fasthttp.StatusMethodNotAllowed, ctx.Response.StatusCode(), "Expected status should be 405 Method Not Allowed")
				// Освобождаем ресурсы
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
			case "Ok":
				controller := gomock.NewController(t)
				defer controller.Finish()

				serv := publ.NewMockInterfaceKafkaClient(controller)
				testCase.mockBehavior(serv, testCase.inputOrder)
				handler := HandlersBuilder{pub: serv, rout: router.New()}
				handler.rout.POST("/WB/set", handler.Set())

				// Create Request
				req := fasthttp.AcquireRequest()

				defer fasthttp.ReleaseRequest(req)

				req.Header.SetMethod("POST")

				req.SetRequestURI("/WB/set")
				req.SetBody([]byte(testCase.inputBody))

				// Create Response Context
				resp := fasthttp.AcquireResponse()
				defer fasthttp.ReleaseResponse(resp)
				ctx := &fasthttp.RequestCtx{}
				ctx.Request = *req
				ctx.Response = *resp

				// Make Request
				handler.rout.Handler(ctx) // Корректный вызов обработчика маршрутизатора

				// Assert

				assert.Equal(t, fasthttp.StatusOK, resp.StatusCode())

				assert.Equal(t, "", string(resp.Body())) // Проверка на пустое тело ответа

				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
			case "equeljson":

				controller := gomock.NewController(t)
				defer controller.Finish()

				// Создаем мок для Kafka клиента
				serv := publ.NewMockInterfaceKafkaClient(controller)

				// Здесь мы говорим, что SendOrderToKafka не должен быть вызван
				serv.EXPECT().SendOrderToKafka(gomock.Any(), gomock.Any()).Times(0)
				// Создаем новый обработчик
				hb := HandlersBuilder{
					pub:  serv,
					rout: router.New(),
				}

				// Регистрируем обработчик для POST метода
				hb.rout.POST("/WB/set", hb.Set())

				// Создаем контекст запроса
				ctx := &fasthttp.RequestCtx{}
				ctx.Request.SetRequestURI("/WB/set")
				ctx.Request.Header.SetMethod("POST")
				ctx.Request.SetBody([]byte(`{}`)) // Устанавливаем тело запроса с невалидным JSON

				// Вызываем обработчик маршрутизатора
				hb.rout.Handler(ctx)

				// Проверяем статус ответа
				assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode(), "Expected status should be 400 Bad Request")
			}

		})
	}
}

// func TestSetGenerateJson(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	testTable := []struct {
// 		name           string
// 		quantityS      string
// 		quantity       int
// 		expectedStatus int
// 		//mockBehavior   func(r *publ.MockInterfaceKafkaClient, g *gen.MockGenerateOrder, quantity int)
// 	}{
// 		{
// 			name:           "Ok",
// 			quantity:       3,
// 			quantityS:      "3",
// 			expectedStatus: fasthttp.StatusOK,
// 			// mockBehavior: func(r *publ.MockInterfaceKafkaClient, g *gen.MockGenerateOrder, quantity int) {
// 			// 	for i := 0; i < quantity; i++ {
// 			// 		orderJson := "orderJson" + strconv.Itoa(i) // Генерируем уникальный JSON для каждого заказа
// 			// 		g.EXPECT().GenerateOrder().Return(orderJson, nil).Times(1)
// 			// 		r.EXPECT().SendOrderToKafka("records", orderJson).Return(nil).Times(1)
// 			// 	}
// 			// },
// 		},
// 		{
// 			name:           "Invalid quantity", //
// 			quantityS:      "invalid",
// 			expectedStatus: fasthttp.StatusBadRequest,
// 			//mockBehavior:   func(r *publ.MockInterfaceKafkaClient, g *gen.MockGenerateOrder, quantity int) {},
// 		},
// 		{
// 			name:           "No quantity specified, default to 5", //
// 			quantityS:      "",
// 			quantity:       5,
// 			expectedStatus: fasthttp.StatusOK,
// 			// mockBehavior: func(r *publ.MockInterfaceKafkaClient, g *gen.MockGenerateOrder, quantity int) {
// 			// 	for i := 0; i < 5; i++ {
// 			// 		orderJson := "orderJson" + strconv.Itoa(i) // Генерируем уникальный JSON для каждого заказа
// 			// 		g.EXPECT().GenerateOrder().Return(orderJson, nil).Times(1)
// 			// 		r.EXPECT().SendOrderToKafka("records", orderJson).Return(nil).Times(1)
// 			// 	}
// 			// },
// 		},
// 		{
// 			name: "Error parsing order JSON",
// 			//quantityS:      "1",
// 			quantity:       1,
// 			expectedStatus: fasthttp.StatusBadRequest,
// 			// mockBehavior: func(r *publ.MockInterfaceKafkaClient, g *gen.MockGenerateOrder, quantity int) {
// 			// 	g.EXPECT().GenerateOrder().Return("invalidOrderJson", nil).Times(1)
// 			// },
// 		},
// 		{
// 			name: "SendOrderToKafka error",
// 			//quantityS:      "1",
// 			quantity:       1,
// 			expectedStatus: fasthttp.StatusOK,
// 			// mockBehavior: func(r *publ.MockInterfaceKafkaClient, g *gen.MockGenerateOrder, quantity int) {
// 			// 	orderJson := "orderJson0"
// 			// 	g.EXPECT().GenerateOrder().Return(orderJson, nil).Times(1)
// 			// 	r.EXPECT().SendOrderToKafka("records", orderJson).Return(errors.New("kafka error")).Times(1)
// 			// },
// 		},
// 		{
// 			name:           "Method not allowed",
// 			quantityS:      "",
// 			expectedStatus: fasthttp.StatusMethodNotAllowed,
// 			//mockBehavior:   func(r *publ.MockInterfaceKafkaClient, g *gen.MockGenerateOrder, quantity int) {},
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			switch testCase.name {
// 			case "methodnotallowed":
// 				controller := gomock.NewController(t)
// 				defer controller.Finish()
// 				serv := publ.NewMockInterfaceKafkaClient(controller)
// 				// Создаем новый обработчик
// 				hb := HandlersBuilder{
// 					pub:  serv, // Подставьте свою реализацию или оставьте nil для теста
// 					rout: router.New(),
// 				}
// 				hb.rout.POST("/WB/set/generate", hb.SetGenerateJson()) // Регистрируем обработчик для POST метода

// 				// Создаем запрос с методом GET
// 				req := fasthttp.AcquireRequest()
// 				req.Header.SetMethod("GET") // Устанавливаем метод GET
// 				req.SetRequestURI("/WB/set/generate")

// 				// Создаем контекст для ответа
// 				resp := fasthttp.AcquireResponse()

// 				// Создаем контекст запроса
// 				ctx := &fasthttp.RequestCtx{}
// 				ctx.Request.SetRequestURI("/WB/set/generate")
// 				ctx.Request.Header.SetMethod("GET")

// 				// Вызываем обработчик маршрутизатора
// 				hb.rout.Handler(ctx)

// 				// Проверяем статус ответа
// 				assert.Equal(t, fasthttp.StatusMethodNotAllowed, ctx.Response.StatusCode(), "Expected status should be 405 Method Not Allowed")

// 				// Освобождаем ресурсы
// 				fasthttp.ReleaseRequest(req)
// 				fasthttp.ReleaseResponse(resp)
// 			case "Ok":
// 				controllerG := gomock.NewController(t)
// 				controllerP := gomock.NewController(t)
// 				defer controllerG.Finish()
// 				defer controllerP.Finish()

// 				serv := publ.NewMockInterfaceKafkaClient(controllerP)
// 				gen := gen.NewMockGenerate(controllerG)
// 				orderJson := "orderJson0"
// 				for i := 0; i < testCase.quantity; i++ {
// 					gen.EXPECT().GenerateOrder().Return(orderJson, nil).Times(1)
// 					serv.EXPECT().SendOrderToKafka("records", orderJson).Return(nil).Times(1)
// 				}
// 				//testCase.mockBehavior(serv, gen, testCase.quantity)
// 				handler := HandlersBuilder{pub: serv, rout: router.New()}
// 				handler.rout.POST("/WB/set/generate", handler.SetGenerateJson())

// 				// Create Request
// 				req := fasthttp.AcquireRequest()

// 				defer fasthttp.ReleaseRequest(req)

// 				req.Header.SetMethod("POST")

// 				req.SetRequestURI("/WB/set/generate")

// 				// Create Response Context
// 				resp := fasthttp.AcquireResponse()
// 				defer fasthttp.ReleaseResponse(resp)
// 				ctx := &fasthttp.RequestCtx{}
// 				ctx.Request = *req
// 				ctx.Response = *resp
// 				ctx.QueryArgs().Add("quantity=", testCase.quantityS) // Задаем количество
// 				// Make Request
// 				handler.rout.Handler(ctx) // Корректный вызов обработчика маршрутизатора

// 				// Assert

// 				assert.Equal(t, fasthttp.StatusOK, resp.StatusCode())

// 				//assert.Equal(t, "", string(resp.Body())) // Проверка на пустое тело ответа

// 				fasthttp.ReleaseRequest(req)
// 				fasthttp.ReleaseResponse(resp)

// 				// case "equeljson":

// 				// 	controller := gomock.NewController(t)
// 				// 	defer controller.Finish()

// 				// 	// Создаем мок для Kafka клиента
// 				// 	serv := publ.NewMockInterfaceKafkaClient(controller)

// 				// 	// Здесь мы говорим, что SendOrderToKafka не должен быть вызван
// 				// 	serv.EXPECT().SendOrderToKafka(gomock.Any(), gomock.Any()).Times(0)
// 				// 	// Создаем новый обработчик
// 				// 	hb := HandlersBuilder{
// 				// 		pub:  serv,
// 				// 		rout: router.New(),
// 				// 	}

// 				// 	// Регистрируем обработчик для POST метода
// 				// 	hb.rout.POST("/WB/set", hb.Set())

// 				// 	// Создаем контекст запроса
// 				// 	ctx := &fasthttp.RequestCtx{}
// 				// 	ctx.Request.SetRequestURI("/WB/set")
// 				// 	ctx.Request.Header.SetMethod("POST")
// 				// 	ctx.Request.SetBody([]byte(`{}`)) // Устанавливаем тело запроса с невалидным JSON

// 				// 	// Вызываем обработчик маршрутизатора
// 				// 	hb.rout.Handler(ctx)

// 				// 	// Проверяем статус ответа
// 				// 	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.StatusCode(), "Expected status should be 400 Bad Request")
// 				// }
// 			}
// 		})
// 	}
// }
