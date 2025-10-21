package rest

import (
	"net/http"
	"strconv"

	"github.com/bharathbbg/order-service/internal/model"
	"github.com/bharathbbg/order-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	orderService *service.OrderService
}

func NewRouter(orderService *service.OrderService) *gin.Engine {
	router := gin.Default()
	r := &Router{orderService: orderService}

	v1 := router.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", r.createOrder)
			orders.GET("", r.listOrders)
			orders.GET("/:id", r.getOrder)
			orders.PUT("/:id", r.updateOrder)
			orders.DELETE("/:id", r.deleteOrder)
		}
	}

	return router
}

func (r *Router) createOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := r.orderService.CreateOrder(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (r *Router) getOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := r.orderService.GetOrder(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (r *Router) listOrders(c *gin.Context) {
	customerID := c.Query("customer_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := r.orderService.ListOrders(c, customerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
	})
}

func (r *Router) updateOrder(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id
	order, err := r.orderService.UpdateOrder(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (r *Router) deleteOrder(c *gin.Context) {
	id := c.Param("id")
	success, err := r.orderService.DeleteOrder(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
