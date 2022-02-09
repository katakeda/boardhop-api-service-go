package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
)

func (s *Service) UserSignup(c *gin.Context) {
	payload := repositories.UserSignupPayload{}
	if err := c.BindJSON(&payload); err != nil {
		log.Println("Failed to parse payload", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong during signup")
		return
	}

	user, err := s.repo.UserSignup(c, payload)
	if err != nil {
		log.Println("Failed to signup user", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong during signup")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Service) GetUser(c *gin.Context) {
	googleAuthId, ok := c.Get("googleAuthId")
	if !ok {
		c.JSON(http.StatusUnauthorized, "Failed to authorize user")
		return
	}

	user, err := s.repo.GetUserByGoogleAuthId(c, googleAuthId)
	if err != nil {
		log.Println("Failed to get user", err)
		c.JSON(http.StatusNotFound, "Failed to get user")
		return
	}

	c.JSON(http.StatusOK, user)
}
