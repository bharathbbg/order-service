package service

import (
	"context"
	"errors"
	"time"

	"github.com/bharathbbg/order-service/internal/model"
	"github.com/bharathbbg/order-service/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	repo  *repository.PostgresRepository
	cache *repository.RedisCache
}

func NewOrderService(repo *repository.PostgresRepository, cache *repository.RedisCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.Order, error) {
	// Validate request
	if req.CustomerID == "" {
		return nil, errors.New("customer_id is required")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// Create order
	order := &model.Order{
		ID:              uuid.New().String(),
		CustomerID:      req.CustomerID,
		Items:           req.Items,
		ShippingAddress: req.ShippingAddress,
		Status:          "PENDING",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Calculate total price
	var totalAmount int64
	for _, item := range req.Items {
		totalAmount += item.UnitPrice.Amount * int64(item.Quantity)
	}

	// Assuming all items have same currency for simplicity
	order.TotalPrice = model.Money{
		Currency: req.Items[0].UnitPrice.Currency,
		Amount:   totalAmount,
	}

	// Save to database
	savedOrder, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := s.cache.CacheOrder(ctx, savedOrder); err != nil {
		// Just log error, don't fail the request
		// log.Printf("Failed to cache order: %v", err)
	}

	return savedOrder, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	// Try to get from cache first
	cachedOrder, err := s.cache.GetCachedOrder(ctx, id)
	if err == nil && cachedOrder != nil {
		return cachedOrder, nil
	}

	// If not in cache, get from database
	res, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	// Ensure the returned value is the expected type
	order, ok := res.(*model.Order)
	if !ok {
		return nil, errors.New("repository returned unexpected type for order")
	}

	// Cache the result for future requests
	if order != nil {
		if err := s.cache.CacheOrder(ctx, order); err != nil {
			// Just log error, don't fail the request
			// log.Printf("Failed to cache order: %v", err)
		}
	}

	return order, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, req *model.UpdateOrderRequest) (*model.Order, error) {
	// Implementation would update the order in database and invalidate cache
	return nil, errors.New("not implemented")
}

func (s *OrderService) ListOrders(ctx context.Context, customerID string, page, pageSize int) ([]*model.Order, int, error) {
	// Implementation would list orders from database
	return nil, 0, errors.New("not implemented")
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) (bool, error) {
	// Implementation would delete order from database and invalidate cache
	return false, errors.New("not implemented")
}
