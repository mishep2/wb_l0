package cache

import (
	"errors"
	"github.com/IsThatASkyline/wb_l0/models"
	"sync"
)

type Cache struct {
	sync.RWMutex
	orders map[string]models.Order
}

// "NewCache нечитать" Этот код определяет пользовательский тип Cache, который имеет следующие поля:
//1. sync.RWMutex - это представляет собой мьютекс для обеспечения безопасности при работе с данными в кэше. Он предоставляет механизм блокировок для чтения и записи.
//2. orders map[string]models.Order - это поле типа map, которое хранит пары ключ-значение. Ключами являются строки (string), а значениями являются объекты типа models.Order. Это предполагает, что в кэше будут храниться заказы (orders) по ключу, который является строкой.
//Таким образом, данный пользовательский тип Cache представляет некий кэш, который использует мьютекс для обеспечения безопасного доступа к данным, и содержит список заказов, доступных по строковому ключу.

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]models.Order), //Эта функция NewCache() создает новый объект типа Cache, инициализируя его пустым списком заказов (orders). Функция возвращает указатель на созданный объект типа Cache.
		//Таким образом, при вызове cache.NewCache() будет создан и возвращен новый объект кэша Cache с пустым списком заказов.
	}
}

func (c *Cache) Set(uid string, order models.Order) {
	c.Lock()
	defer c.Unlock()
	c.orders[uid] = order
}

func (c *Cache) GetOrderByUid(uid string) (models.Order, error) {
	c.RLock()
	defer c.RUnlock()
	order, ok := c.orders[uid]
	if !ok {
		return order, errors.New("error cache for GetOrderByUid")
	}
	return order, nil
}
