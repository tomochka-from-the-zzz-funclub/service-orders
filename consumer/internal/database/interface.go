package database

import "consumer/internal/models"

type InterfacePostgresDB interface {
	AddOrderStruct(order models.Order) (string, error)
	AddOrder(order models.Order, items []string, delivery string, payment string) error
	AddItems(items []models.Item) ([]string, error)
	AddPayment(payment models.Payment) (string, error)
	AddDelivery(delivery models.Delivery) (string, error)
	GetOrder(order_uuid string) (models.Order, error)
}
