package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
)

func (s *Service) CreateMessage(c *gin.Context) {
	s.createMessage(c)
}

func (s *Service) createMessage(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to create message |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong while creating message")
		}
	}()

	payload := repositories.CreateMessagePayload{}
	if err := c.BindJSON(&payload); err != nil {
		return fmt.Errorf("failed to parse payload | %w", err)
	}

	user, err := s.getUser(c)
	if err != nil || user == nil {
		return fmt.Errorf("failed to authorize user | %w", err)
	}

	payload.UserId = user.Id

	message, err := s.repo.CreateMessage(c, payload)
	if err != nil {
		return fmt.Errorf("failed to insert message | %w", err)
	}

	c.JSON(http.StatusOK, message)

	return nil
}
