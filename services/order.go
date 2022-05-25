package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
)

func (s *Service) GetOrder(c *gin.Context) {
	s.getOrder(c)
}

func (s *Service) CreateOrder(c *gin.Context) {
	s.createOrder(c)
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
