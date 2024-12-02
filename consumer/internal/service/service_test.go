package service

import (
	"testing"

	cmock "consumer/internal/cache/mocks"     // Путь к вашим мокам
	dbmock "consumer/internal/database/mocks" // Путь к вашим мокам
	myErrors "consumer/internal/errors"
	"consumer/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetOrderSrv_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := cmock.NewMockCache(ctrl)             // Создаем мок для кеша
	mockDB := dbmock.NewMockInterfacePadtgresDB(ctrl) // Создаем мок для базы данных

	orderUUID := "b563feb7b2b84b6test"
	expectedOrder := models.Order{OrderUID: orderUUID}

	// Настраиваем мок кеша
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

	mockCache := cmock.NewMockCache(ctrl)             // Создаем мок для кеша
	mockDB := dbmock.NewMockInterfacePadtgresDB(ctrl) // Создаем мок для базы данных

	orderUUID := "b563feb7b2b84b6test"
	expectedOrder := models.Order{OrderUID: orderUUID} // Ожидаемый объект заказа

	// Настраиваем мок кеша для возврата пустого заказа при кеш-промахе
	mockCache.EXPECT().Get(orderUUID).Return(models.Order{}, myErrors.ErrNotFoundOrderCache)

	// Настраиваем мок базы данных для возврата ожидаемого заказа
	mockDB.EXPECT().GetOrder(orderUUID).Return(expectedOrder, nil)

	// Ожидаем вызов метода Add с ожидаемым объектом
	mockCache.EXPECT().Add(expectedOrder).Return(nil) // Добавьте эту строку, если метод Add должен быть вызван

	// Настраиваем мок базы данных для возврата ожидаемого заказа

	srv := Srv{
		cache: mockCache,
		db:    mockDB,
	}

	order, err := srv.GetOrderSrv(orderUUID)

	assert.NoError(t, err)                // Проверяем, что ошибка отсутствует
	assert.Equal(t, expectedOrder, order) // Проверяем, что возвращаемый заказ соответствует ожидаемому
}

func TestSetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := cmock.NewMockCache(ctrl)             // Создаем мок для кеша
	mockDB := dbmock.NewMockInterfacePadtgresDB(ctrl) // Создаем мок для базы данных

	order := models.Order{OrderUID: "b563feb7b2b84b6test"}

	// Настраиваем мок базы данных
	mockDB.EXPECT().AddOrder(order).Return(order.OrderUID, nil)
	// Настраиваем мок кеша (предполагается, что при добавлении заказа в БД мы также добавляем его в кеш)
	mockCache.EXPECT().Add(order).Return(nil)

	srv := Srv{
		cache: mockCache,
		db:    mockDB,
	}

	id, err := srv.SetOrder(order)

	assert.NoError(t, err)
	assert.Equal(t, order.OrderUID, id)
}
