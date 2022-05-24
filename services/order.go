package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
)

func (s *Service) CreateOrder(c *gin.Context) {
	s.createOrder(c)
}

func (s *Service) createOrder(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to create order |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong while creating order")
		}
	}()

	payload := repositories.CreateOrderPayload{}
	if err := c.BindJSON(&payload); err != nil {
		return fmt.Errorf("failed to parse payload | %w", err)
	}

	order, err := s.repo.CreateOrder(c, payload)
	if err != nil {
		return fmt.Errorf("failed to insert order | %w", err)
	}

	c.JSON(http.StatusOK, order)

	return nil
}
