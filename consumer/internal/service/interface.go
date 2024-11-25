package service

import (
	"consumer/internal/config"
	"consumer/internal/models"
)

type InterfaceService interface {
	GetOrderSrv(orderUUID string) (models.Order, error)
	SetOrder(order models.Order) (string, error)
	Read(cfg config.Config)
}
