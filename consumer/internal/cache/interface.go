package store

import "consumer/internal/models"

type Cache interface {
	Add(order models.Order) error
	Get(OrderUID string) (models.Order, error)
	StartGC()
	GC()
	DeleteExpiredKeys()
}
