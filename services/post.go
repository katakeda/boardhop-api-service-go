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
		c.IndentedJSON(http.StatusInternalServerError, "Something went wrong while getting posts")
	}

	if len(posts) <= 0 {
		log.Println("No posts found")
		c.IndentedJSON(http.StatusNotFound, "No posts found")
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func (s *Service) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := s.repo.GetPost(c, id)
	if err != nil {
		log.Println("Failed to get post", err)
		c.IndentedJSON(http.StatusInternalServerError, "Something went wrong while getting post")
	}

	c.IndentedJSON(http.StatusOK, post)
}

func (s *Service) CreatePost(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "OK")
}

func (s *Service) GetCategories(c *gin.Context) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		log.Println("Failed to get categories", err)
		c.IndentedJSON(http.StatusInternalServerError, "Something went wrong while getting categories")
	}

	if len(categories) <= 0 {
		log.Println("No categories found")
		c.IndentedJSON(http.StatusNotFound, "No categories found")
	}

	c.IndentedJSON(http.StatusOK, categories)
}
