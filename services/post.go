package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
)

func (s *Service) GetPosts(c *gin.Context) {
	params := c.Request.URL.Query()

	posts, err := s.repo.GetPosts(c, params)
	if err != nil {
		log.Println("Failed to get posts", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting posts")
		return
	}

	if len(posts) <= 0 {
		log.Println("No posts found")
		c.JSON(http.StatusNotFound, "No posts found")
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (s *Service) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := s.repo.GetPost(c, id)
	if err != nil {
		log.Println("Failed to get post", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting post")
		return
	}

	c.JSON(http.StatusOK, post)
}

func (s *Service) CreatePost(c *gin.Context) {
	payload := repositories.CreatePostPayload{}
	if err := c.BindJSON(&payload); err != nil {
		log.Println("Failed to parse payload", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		return
	}

	post, err := s.repo.CreatePost(c, payload)
	if err != nil {
		log.Println("Failed to create post", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		return
	}

	c.JSON(http.StatusOK, post)
}

func (s *Service) GetCategories(c *gin.Context) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		log.Println("Failed to get categories", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting categories")
		return
	}

	if len(categories) <= 0 {
		log.Println("No categories found")
		c.JSON(http.StatusNotFound, "No categories found")
		return
	}

	c.JSON(http.StatusOK, categories)
}
