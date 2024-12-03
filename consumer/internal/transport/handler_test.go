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
		name                 string
		inputURLarg          string
		inputOrderid         string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody models.Order
	}{
		{
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
				r.EXPECT().GetOrderSrv(id).Return(models.Order{}, myErrors.ErrNotFoundOrder)
			},
			expectedStatusCode:   fasthttp.StatusNotFound,
			expectedResponseBody: models.Order{},
		},
		{
			name:         "servererror",
			inputURLarg:  "order_uid=b563feb7b2b84b6test",
			inputOrderid: "b563feb7b2b84b6test",
			mockBehavior: func(r *mocks.MockInterfaceService, id string) {
				r.EXPECT().GetOrderSrv(id).Return(models.Order{}, myErrors.NewError(fasthttp.StatusInternalServerError, "server err"))
			},
			expectedStatusCode:   fasthttp.StatusInternalServerError,
			expectedResponseBody: models.Order{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			serv := mocks.NewMockInterfaceService(controller)

			testCase.mockBehavior(serv, testCase.inputOrderid)

			hb := HandlersBuilder{
				srv:  serv,
				rout: router.New(),
			}
			hb.rout.GET("/WB/get", hb.Get())

			ctx := &fasthttp.RequestCtx{}
			ctx.Request.SetRequestURI(string("/WB/get?" + testCase.inputURLarg))

			if testCase.name == "methodnotallowed" {
				ctx.Request.Header.SetMethod("POST")
			} else {
				ctx.Request.Header.SetMethod("GET")
			}

			hb.rout.Handler(ctx)

			assert.Equal(t, testCase.expectedStatusCode, ctx.Response.StatusCode(), "Expected status code does not match")

			if testCase.expectedStatusCode == fasthttp.StatusOK {
				var order models.Order
				err := json.Unmarshal(ctx.Response.Body(), &order)
				assert.NoError(t, err, "Error unmarshalling response body")
				assert.Equal(t, testCase.expectedResponseBody, order, "Response body does not match")
			}
		})
	}
}
