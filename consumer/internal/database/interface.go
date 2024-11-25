package database

import "consumer/internal/models"

type InterfacePadtgresDB interface {
	AddOrder(order models.Order) (string, error)
	GetOrder(order_uuid string) (models.Order, error)
}
