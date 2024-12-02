package service

import (
	"consumer/internal/config"
	"consumer/internal/models"
)

//mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mocks

type InterfaceService interface {
	GetOrderSrv(orderUUID string) (models.Order, error)
	SetOrder(order models.Order) (string, error)
	Read(cfg config.Config)
}
