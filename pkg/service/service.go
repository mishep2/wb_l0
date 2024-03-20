package service

import (
	"github.com/mishep2/wb_l0/models"
	"github.com/mishep2/wb_l0/pkg/repository"
)

type Order interface {
	GetOrderByUid(uid string) (models.Order, error)
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.Order),
	}
}
