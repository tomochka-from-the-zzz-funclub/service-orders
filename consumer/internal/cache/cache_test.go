package store

import (
	"consumer/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	m := map[string]models.Order{}
	s := NewStore(m)
	order := models.Order{OrderUID: "test-uuid"}

	// Добавляем заказ
	err := s.Add(order)
	assert.NoError(t, err)

	// Пытаемся добавить тот же заказ снова
	err = s.Add(order)
	assert.Error(t, err)
	assert.Equal(t, "record with this OrderUID exists in database", err.Error())
}

func TestGet(t *testing.T) {
	m := map[string]models.Order{}
	s := NewStore(m)
	order := models.Order{OrderUID: "test-uuid"}

	// Пытаемся получить заказ, который отсутствует
	_, err := s.Get(order.OrderUID)
	assert.Error(t, err)
	assert.Equal(t, "there isn't record with such OrderUID", err.Error())

	// Добавляем заказ
	err = s.Add(order)
	assert.NoError(t, err)

	// Получаем заказ
	retrievedOrder, err := s.Get(order.OrderUID)
	assert.NoError(t, err)
	assert.Equal(t, order, retrievedOrder)
}

func TestDeleteExpiredKeys(t *testing.T) {
	m := map[string]models.Order{}
	s := NewStore(m)
	order := models.Order{OrderUID: "test-uuid"}

	// Добавляем заказ
	err := s.Add(order)
	assert.NoError(t, err)

	// Убеждаемся, что заказ существует
	retrievedOrder, err := s.Get(order.OrderUID)
	assert.NoError(t, err)
	assert.Equal(t, order, retrievedOrder)

	// Устанавливаем время создания заказа на 2 часа назад
	s.data[order.OrderUID] = OrderCache{
		Ord:        order,
		TimeCreate: time.Now().Add(-2 * time.Hour),
	}

	// Запускаем сборщик мусора
	s.DeleteExpiredKeys()

	// Проверяем, что заказ был удален
	_, err = s.Get(order.OrderUID)
	assert.Error(t, err)
	assert.Equal(t, "there isn't record with such OrderUID", err.Error())
}
