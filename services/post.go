package services

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
	"github.com/katakeda/boardhop-api-service-go/utils"
)

type CreatePostPayload struct {
	Data   repositories.CreatePost `json:"data" form:"data" binding:"required"`
	Images []*multipart.FileHeader `json:"images" form:"images"`
}

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
		log.Printf("Failed to get post | %v", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting post")
		return
	}

	c.JSON(http.StatusOK, post)
}

func (s *Service) CreatePost(c *gin.Context) {
	payload := CreatePostPayload{}
	if err := c.ShouldBind(&payload); err != nil {
		log.Println("Failed to parse payload |", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		return
	}

	ctx, err := s.repo.BeginTxn(c)
	if err != nil {
		log.Println("Failed to begin db txn |", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		return
	}

	post, err := s.repo.CreatePost(ctx, payload.Data)
	if err != nil {
		log.Println("Failed to create post |", err)
		s.repo.RollbackTxn(ctx)
		c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		return
	}

	images := []repositories.CreatePostMedia{}
	for idx := range payload.Images {
		images = append(images, repositories.CreatePostMedia{
			PostId:   post.Id,
			MediaUrl: fmt.Sprintf("%s_%d.jpg", post.Id, idx),
			Type:     "image",
		})
	}

	if err := s.repo.CreatePostMedias(ctx, images); err != nil {
		log.Println("Failed to create post medias |", err)
		s.repo.RollbackTxn(ctx)
		c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		return
	}

	bucket, err := utils.GetDefaultBucket(c)
	if err != nil {
		log.Println("Failed to get default bucket |", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		return
	}

	for idx := range payload.Images {
		if err := utils.UploadFile(c, bucket, images[idx].MediaUrl, payload.Images[idx]); err != nil {
			// Revert all uploaded
			log.Println("Failed to upload files |", err)
			s.repo.RollbackTxn(ctx)
			c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
			return
		}
	}

	s.repo.CommitTxn(ctx)

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
