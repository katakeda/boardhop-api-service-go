package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) GetPosts(c *gin.Context) {
	params := c.Request.URL.Query()

	posts, err := s.repo.GetPosts(c, params)
	if err != nil {
		log.Println("Failed to get posts", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting posts")
	}

	if len(posts) <= 0 {
		log.Println("No posts found")
		c.JSON(http.StatusNotFound, "No posts found")
	}

	c.JSON(http.StatusOK, posts)
}

func (s *Service) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := s.repo.GetPost(c, id)
	if err != nil {
		log.Println("Failed to get post", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting post")
	}

	c.JSON(http.StatusOK, post)
}

func (s *Service) CreatePost(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func (s *Service) GetCategories(c *gin.Context) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		log.Println("Failed to get categories", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting categories")
	}

	if len(categories) <= 0 {
		log.Println("No categories found")
		c.JSON(http.StatusNotFound, "No categories found")
	}

	c.JSON(http.StatusOK, categories)
}
