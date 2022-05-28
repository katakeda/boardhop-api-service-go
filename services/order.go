package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
)

func (s *Service) GetOrders(c *gin.Context) {
	s.getOrders(c)
}

func (s *Service) GetOrder(c *gin.Context) {
	s.getOrder(c)
}

func (s *Service) CreateOrder(c *gin.Context) {
	s.createOrder(c)
}

func (s *Service) getOrders(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to get orders |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong while getting orders")
		}
	}()

	user, err := s.getUser(c)
	if err != nil {
		return fmt.Errorf("failed to authorize user | %w", err)
	}

	orders, err := s.repo.GetOrders(c, repositories.GetOrdersFilter{UserId: &user.Id})
	if err != nil {
		return fmt.Errorf("failed to fetch orders | %w", err)
	}

	c.JSON(http.StatusOK, orders)

	return nil
}

func (s *Service) getOrder(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to get order |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong while getting order")
		}
	}()

	id := c.Param("id")
	order, err := s.repo.GetOrder(c, id)
	if err != nil {
		return fmt.Errorf("failed to get order | %w", err)
	}

	c.JSON(http.StatusOK, order)

	return nil
}

func (s *Service) createOrder(c *gin.Context) (err error) {
	ctx, _ := s.repo.BeginTxn(c)

	defer func() {
		if err != nil {
			log.Println("Failed to create order |", err)
			s.repo.RollbackTxn(ctx)
			c.JSON(http.StatusInternalServerError, "Something went wrong while creating order")
		}
	}()

	payload := repositories.CreateOrderPayload{}
	if err := c.BindJSON(&payload); err != nil {
		return fmt.Errorf("failed to parse payload | %w", err)
	}

	order, err := s.repo.CreateOrder(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to insert order | %w", err)
	}

	if payload.Message != nil {
		message, err := s.repo.CreateMessage(ctx, repositories.CreateMessagePayload{OrderId: &order.Id, Message: payload.Message})
		if err != nil {
			return fmt.Errorf("failed to insert order message | %w", err)
		}
		order.Messages = []repositories.Message{*message}
	}

	c.JSON(http.StatusOK, order)

	return s.repo.CommitTxn(ctx)
}
