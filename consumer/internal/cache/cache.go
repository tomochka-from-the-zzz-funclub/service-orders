package store

import (
	myLog "consumer/internal/logger"
	"consumer/internal/models"
	"errors"

	"sync"
	"time"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]OrderCache
}

type OrderCache struct {
	Ord        models.Order
	TimeCreate time.Time
}

func NewOrderCache(ord models.Order) OrderCache {

	return OrderCache{
		Ord:        ord,
		TimeCreate: time.Now(),
	}
}

func NewStore(c map[string]models.Order) *Store {
	result := make(map[string]OrderCache)
	for id, str := range c {
		result[id] = NewOrderCache(str)
	}
	return &Store{data: result}
}

func (s *Store) Add(order models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	myLog.Log.Debugf("Start fun Add in cache")
	_, ok := s.data[order.OrderUID]
	if ok {
		return errors.New("record with this OrderUID exists in database")
	}
	s.data[order.OrderUID] = NewOrderCache(order)
	myLog.Log.Debugf("Succes add order in cache")
	return nil
}

func (s *Store) Get(OrderUID string) (models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[OrderUID]
	if ok {
		return val.Ord, nil
	}
	return models.Order{}, errors.New("there isn't record with such OrderUID")
}

func (c *Store) StartGC() {
	go c.GC()
}

func (c *Store) GC() {
	for {
		time.Sleep(10 * time.Minute)
		c.DeleteExpiredKeys()
	}

}

func (c *Store) DeleteExpiredKeys() {

	c.mu.Lock()

	defer c.mu.Unlock()

	for k, i := range c.data {
		if time.Hour < time.Since(i.TimeCreate) {
			delete(c.data, k)
		}
	}
}
