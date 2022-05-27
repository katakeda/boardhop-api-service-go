package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		app, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Println("Error initializing app |", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		idToken, err := parseIdToken(c)
		if err != nil {
			log.Println("Error parsing ID token |", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		client, err := app.Auth(c)
		if err != nil {
			log.Println("Error getting Auth client |", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		token, err := client.VerifyIDToken(c, *idToken)
		if err != nil {
			log.Println("Error verifying ID token |", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("googleAuthId", token.UID)

		c.Next()
	}
}

func parseIdToken(c *gin.Context) (*string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	authArr := strings.Split(authHeader, "Bearer ")

	if len(authArr) > 1 {
		return &authArr[1], nil
	}

	return nil, fmt.Errorf("failed to parse Authorization header")
}
