package transport

import (
	"encoding/json"
	"testing"

	"github.com/fasthttp/router"
	"github.com/golang/mock/gomock"
	"github.com/valyala/fasthttp"

	"github.com/stretchr/testify/assert"

	myErrors "consumer/internal/errors"
	"consumer/internal/models"
	"consumer/internal/service/mocks"
)

func TestSetHandler_Get(t *testing.T) {
	type mockBehavior func(r *mocks.MockInterfaceService, ord string)

	testTable := []struct {
		name                 string // имя теста
		inputURLarg          string //тело запроса
		inputOrderid         string //структура заказа, которую мы пеердаем в метод сервиса
		mockBehavior         mockBehavior
		expectedStatusCode   int          // ожидаемый статус код
		expectedResponseBody models.Order //тело ответа
	}{
		// 	{ // тест на успешный вызов SendOrderToKafka
		// 		name:         "Ok",
		// 		inputURLarg:  "order_uid=b563feb7b2b84b6test",
		// 		inputOrderid: "b563feb7b2b84b6test",
		// 		mockBehavior: func(r *mocks.MockInterfaceService, id string) {
		// 			r.EXPECT().GetOrderSrv(id).Return(nil)
		// 		},
		// 		expectedStatusCode: 200,
		// 		expectedResponseBody: models.Order{
		// 			OrderUID:    "b563feb7b2b84b6test",
		// 			TrackNumber: "WBILMTESTTRACK",
		// 			Entry:       "WBIL",
		// 			Delivery: models.Delivery{
		// 				Name:    "Test Testov",
		// 				Phone:   "+9720000000",
		// 				Zip:     "2639809",
		// 				City:    "Kiryat Mozkin",
		// 				Address: "Ploshad Mira 15",
		// 				Region:  "Kraiot",
		// 				Email:   "test@gmail.com",
		// 			},
		// 			Payment: models.Payment{
		// 				Transaction:  "b563feb7b2b84b6test",
		// 				RequestID:    "",
		// 				Currency:     "USD",
		// 				Provider:     "wbpay",
		// 				Amount:       1817,
		// 				PaymentDT:    1637907727,
		// 				Bank:         "alpha",
		// 				DeliveryCost: 1500,
		// 				GoodsTotal:   317,
		// 				CustomFee:    0,
		// 			},
		// 			Items: []models.Item{
		// 				{
		// 					ChrtID:      9934930,
		// 					TrackNumber: "WBILMTESTTRACK",
		// 					Price:       453,
		// 					Rid:         "ab4219087a764ae0btest",
		// 					Name:        "Mascaras",
		// 					Sale:        30,
		// 					Size:        "0",
		// 					TotalPrice:  317,
		// 					NmID:        2389212,
		// 					Brand:       "Vivienne Sabo",
		// 					Status:      202,
		// 				},
		// 			},
		// 			Locale:            "en",
		// 			InternalSignature: "",
		// 			CustomerID:        "test",
		// 			DeliveryService:   "meest",
		// 			ShardKey:          "9",
		// 			SmID:              99,
		// 			DateCreated:       "2021-11-26T06:22:19Z",
		// 			OofShard:          "1",
		// 		},
		// 	},

		// 	{
		// 		name:         "methodnotallowed",
		// 		inputURLarg:  "order_uid=b563feb7b2b84b6test",
		// 		inputOrderid: "b563feb7b2b84b6test",

		// 		mockBehavior: func(r *mocks.MockInterfaceService, id string) {
		// 			//r.EXPECT().SendOrderToKafka("records", ord).Return(nil).Times(1)
		// 		},
		// 		expectedStatusCode:   405,
		// 		expectedResponseBody: models.Order{},
		// 	},
		// }

		// for _, testCase := range testTable {
		// 	t.Run(testCase.name, func(t *testing.T) {
		// 			switch testCase.name {
		// 			case "methodnotallowed":
		// 				controller := gomock.NewController(t)
		// 				defer controller.Finish()
		// 				serv := mocks.NewMockInterfaceService(controller)
		// 				// Создаем новый обработчик
		// 				hb := HandlersBuilder{
		// 					srv:  serv, // Подставьте свою реализацию или оставьте nil для теста
		// 					rout: router.New(),
		// 				}
		// 				hb.rout.GET("/WB/get", hb.Get()) // Регистрируем обработчик для POST метода

		// 				// Создаем запрос с методом GET
		// 				req := fasthttp.AcquireRequest()
		// 				req.Header.SetMethod("POST") // Устанавливаем метод GET
		// 				req.SetRequestURI("/WB/get?order_uid=b563feb7b2b84b6test")

		// 				// Создаем контекст для ответа
		// 				resp := fasthttp.AcquireResponse()

		// 				// Создаем контекст запроса
		// 				ctx := &fasthttp.RequestCtx{}
		// 				ctx.Request.SetRequestURI("/WB/get?order_uid=b563feb7b2b84b6test")
		// 				ctx.Request.Header.SetMethod("POST")

		// 				// Вызываем обработчик маршрутизатора
		// 				hb.rout.Handler(ctx)

		// 				// Проверяем статус ответа
		// 				assert.Equal(t, fasthttp.StatusMethodNotAllowed, ctx.Response.StatusCode(), "Expected status should be 405 Method Not Allowed")

		// 				// Освобождаем ресурсы
		// 				fasthttp.ReleaseRequest(req)
		// 				fasthttp.ReleaseResponse(resp)
		// 			case "Ok":
		// 				controller := gomock.NewController(t)
		// 				defer controller.Finish()

		// 				serv := mocks.NewMockInterfaceKafkaClient(controller)
		// 				testCase.mockBehavior(serv, testCase.inputOrder)
		// 				handler := HandlersBuilder{pub: serv, rout: router.New()}
		// 				handler.rout.POST("/WB/set", handler.Set())

		// 				// Create Request
		// 				req := fasthttp.AcquireRequest()

		// 				defer fasthttp.ReleaseRequest(req)

		// 				req.Header.SetMethod("POST")

		// 				req.SetRequestURI("/WB/set")
		// 				req.SetBody([]byte(testCase.inputBody))

		// 				// Create Response Context
		// 				resp := fasthttp.AcquireResponse()
		// 				defer fasthttp.ReleaseResponse(resp)
		// 				ctx := &fasthttp.RequestCtx{}
		// 				ctx.Request = *req
		// 				ctx.Response = *resp

		// 				// Make Request
		// 				handler.rout.Handler(ctx) // Корректный вызов обработчика маршрутизатора

		// 				// Assert

		// 				assert.Equal(t, fasthttp.StatusOK, resp.StatusCode())

		// 				assert.Equal(t, "", string(resp.Body())) // Проверка на пустое тело ответа

		// 				fasthttp.ReleaseRequest(req)
		// 				fasthttp.ReleaseResponse(resp)
		// 			}
		// 		})
		// 	}
		// }

		{ // тест на успешный вызов
			name:         "Ok",
			inputURLarg:  "order_uid=b563feb7b2b84b6test",
			inputOrderid: "b563feb7b2b84b6test",
			mockBehavior: func(r *mocks.MockInterfaceService, id string) {
				r.EXPECT().GetOrderSrv(id).Return(models.Order{
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
				}, nil)
			},
			expectedStatusCode: fasthttp.StatusOK,
			expectedResponseBody: models.Order{
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
		},
		{
			name:                 "methodnotallowed",
			inputURLarg:          "order_uid=b563feb7b2b84b6test",
			inputOrderid:         "b563feb7b2b84b6test",
			mockBehavior:         func(r *mocks.MockInterfaceService, id string) {},
			expectedStatusCode:   fasthttp.StatusMethodNotAllowed,
			expectedResponseBody: models.Order{},
		},
		{
			name:         "notFoundOrder",
			inputURLarg:  "order_uid=b563feb7b2b84b6test",
			inputOrderid: "b563feb7b2b84b6test",
			mockBehavior: func(r *mocks.MockInterfaceService, id string) {
				r.EXPECT().GetOrderSrv(id).Return(models.Order{}, myErrors.NewError(fasthttp.StatusNotFound, "not_found"))
			},
			expectedStatusCode:   fasthttp.StatusNotFound,
			expectedResponseBody: models.Order{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			serv := mocks.NewMockInterfaceService(controller)

			// Настройка поведения мока
			testCase.mockBehavior(serv, testCase.inputOrderid)

			// Создаем новый обработчик
			hb := HandlersBuilder{
				srv:  serv,
				rout: router.New(),
			}
			hb.rout.GET("/WB/get", hb.Get()) // Регистрируем обработчик для GET метода

			// Создаем контекст запроса
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.SetRequestURI(string("/WB/get?" + testCase.inputURLarg))

			if testCase.name == "methodnotallowed" {
				ctx.Request.Header.SetMethod("POST") // Устанавливаем метод POST для теста "methodnotallowed"
			} else {
				ctx.Request.Header.SetMethod("GET") // Устанавливаем метод GET для успешного теста
			}

			// Вызываем обработчик маршрутизатора
			hb.rout.Handler(ctx)

			// Проверяем статус ответа
			assert.Equal(t, testCase.expectedStatusCode, ctx.Response.StatusCode(), "Expected status code does not match")

			// Проверяем тело ответа (если требуется)
			if testCase.expectedStatusCode == fasthttp.StatusOK {
				var order models.Order
				err := json.Unmarshal(ctx.Response.Body(), &order)
				assert.NoError(t, err, "Error unmarshalling response body")
				assert.Equal(t, testCase.expectedResponseBody, order, "Response body does not match")
			}
		})
	}
}
