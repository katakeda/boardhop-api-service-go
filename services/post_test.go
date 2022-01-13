package services

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/mocks"
	"github.com/katakeda/boardhop-api-service-go/repositories"
	"github.com/stretchr/testify/mock"
)

func TestGetPosts(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(mocks.IRepository)

	mockRepo.
		On("GetPosts", ctx, mock.AnythingOfType("url.Values{}")).
		Return([]repositories.Post{}, nil)

	svc, _ := NewService(mockRepo)

	router := gin.Default()
	router.GET("/posts", svc.GetPosts)
}
