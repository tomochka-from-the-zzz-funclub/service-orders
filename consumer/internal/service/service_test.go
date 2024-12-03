package service

import (
	"testing"

	cmock "consumer/internal/cache/mocks"
	dbmock "consumer/internal/database/mocks"
	myErrors "consumer/internal/errors"
	"consumer/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetOrderSrv_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := cmock.NewMockCache(ctrl)
	mockDB := dbmock.NewMockInterfacePostgresDB(ctrl)

	orderUUID := "b563feb7b2b84b6test"
	expectedOrder := models.Order{OrderUID: orderUUID}

	mockCache.EXPECT().Get(orderUUID).Return(expectedOrder, nil)

	srv := Srv{
		cache: mockCache,
		db:    mockDB,
	}

	order, err := srv.GetOrderSrv(orderUUID)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestGetOrderSrv_CacheMiss_DBHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := cmock.NewMockCache(ctrl)
	mockDB := dbmock.NewMockInterfacePostgresDB(ctrl)

	orderUUID := "b563feb7b2b84b6test"
	expectedOrder := models.Order{OrderUID: orderUUID}

	mockCache.EXPECT().Get(orderUUID).Return(models.Order{}, myErrors.ErrNotFoundOrderCache)

	mockDB.EXPECT().GetOrder(orderUUID).Return(expectedOrder, nil)

	mockCache.EXPECT().Add(expectedOrder).Return(nil)

	srv := Srv{
		cache: mockCache,
		db:    mockDB,
	}

	order, err := srv.GetOrderSrv(orderUUID)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestSetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := cmock.NewMockCache(ctrl)
	mockDB := dbmock.NewMockInterfacePostgresDB(ctrl)

	order := models.Order{OrderUID: "b563feb7b2b84b6test"}

	mockDB.EXPECT().AddOrderStruct(order).Return(order.OrderUID, nil)
	mockCache.EXPECT().Add(order).Return(nil)

	srv := Srv{
		cache: mockCache,
		db:    mockDB,
	}

	id, err := srv.SetOrder(order)

	assert.NoError(t, err)
	assert.Equal(t, order.OrderUID, id)
}
