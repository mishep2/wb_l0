package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mishep2/wb_l0/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		orders := api.Group("/orders")
		{
			orders.GET("/:id", h.getOrderByUid)
		}
	}

	return router
}
