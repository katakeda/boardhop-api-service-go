package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
)

func (s *Service) UserSignup(c *gin.Context) {
	s.userSignup(c)
}

func (s *Service) UserLogin(c *gin.Context) {
	payload := repositories.UserLoginPayload{}
	if err := c.BindJSON(&payload); err != nil {
		log.Println("Failed to parse payload", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong during login")
		return
	}

	user, err := s.repo.GetUserByGoogleAuthId(c, payload.GoogleAuthId)
	if err != nil {
		log.Println("Failed to get user", err)
		c.JSON(http.StatusNotFound, "Failed to get user")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Service) GetUser(c *gin.Context) {
	user, err := s.getUser(c)
	if err != nil {
		log.Println("Failed to get user |", err)
		c.JSON(http.StatusUnauthorized, "Failed to authorize user")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Service) getUser(c *gin.Context) (user *repositories.User, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to get user | %w", err)
		}
	}()

	googleAuthId, ok := c.Get("googleAuthId")
	if !ok {
		return nil, fmt.Errorf("failed to authorize user")
	}

	user, err = s.repo.GetUserByGoogleAuthId(c, googleAuthId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user | %w", err)
	}

	return user, nil
}

func (s *Service) userSignup(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to signup user |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong during signup")
		}
	}()

	payload := repositories.UserSignupPayload{}
	if err := c.BindJSON(&payload); err != nil {
		return fmt.Errorf("failed to parse payload | %w", err)
	}

	user, err := s.repo.UserSignup(c, payload)
	if err != nil {
		return fmt.Errorf("failed to insert user | %w", err)
	}

	c.JSON(http.StatusOK, user)

	return nil
}
