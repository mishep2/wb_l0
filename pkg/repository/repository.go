package repository

import (
	"github.com/mishep2/wb_l0/models"
	"github.com/mishep2/wb_l0/pkg/repository/cache"
)

type Order interface {
	GetOrderByUid(uid string) (models.Order, error)
}

type Repository struct {
	Order
}

func NewRepository(cache *cache.Cache) *Repository {
	return &Repository{
		Order: NewOrderCache(cache),
	}
}
