package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katakeda/boardhop-api-service-go/repositories"
	"github.com/katakeda/boardhop-api-service-go/utils"
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
		log.Printf("Failed to get post | %v", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong while getting post")
		return
	}

	c.JSON(http.StatusOK, post)
}

func (s *Service) CreatePost(c *gin.Context) {
	s.createPost(c)
}

func (s *Service) GetTags(c *gin.Context) {
	s.getTags(c)
}

func (s *Service) GetCategories(c *gin.Context) {
	s.getCategories(c)
}

func (s *Service) createPost(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to create post |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong while creating post")
		}
	}()

	payload := repositories.CreatePostPayload{}
	if err = c.ShouldBind(&payload); err != nil {
		return fmt.Errorf("failed to parse payload | %w", err)
	}

	ctx, err := s.repo.BeginTxn(c)
	if err != nil {
		return fmt.Errorf("failed to begin db txn | %w", err)
	}

	post, err := s.repo.CreatePost(ctx, payload.Data)
	if err != nil {
		s.repo.RollbackTxn(ctx)
		return fmt.Errorf("failed to insert post | %w", err)
	}

	tags := []repositories.CreatePostTag{}
	for idx := range payload.TagIds {
		tags = append(tags, repositories.CreatePostTag{
			PostId: post.Id,
			TagId:  payload.TagIds[idx],
		})
	}

	if len(tags) > 0 {
		if err = s.repo.CreatePostTags(ctx, tags); err != nil {
			s.repo.RollbackTxn(ctx)
			return fmt.Errorf("failed to create post tags | %w", err)
		}
	}

	categories := []repositories.CreatePostCategory{}
	for idx := range payload.CategoryIds {
		categories = append(categories, repositories.CreatePostCategory{
			PostId:     post.Id,
			CategoryId: payload.CategoryIds[idx],
		})
	}

	if len(categories) > 0 {
		if err = s.repo.CreatePostCategories(ctx, categories); err != nil {
			s.repo.RollbackTxn(ctx)
			return fmt.Errorf("failed to create post categories | %w", err)
		}
	}

	images := []repositories.CreatePostMedia{}
	for idx := range payload.Images {
		images = append(images, repositories.CreatePostMedia{
			PostId:   post.Id,
			MediaUrl: fmt.Sprintf("%s_%d.jpg", post.Id, idx),
			Type:     "image",
		})
	}

	if len(images) > 0 {
		if err = s.repo.CreatePostMedias(ctx, images); err != nil {
			s.repo.RollbackTxn(ctx)
			return fmt.Errorf("failed to create post medias | %w", err)
		}

		bucket, err := utils.GetDefaultBucket(c)
		if err != nil {
			return fmt.Errorf("failed to get default bucket | %w", err)
		}

		var uploadErr error
		defer func() {
			if uploadErr != nil {
				for idx := range images {
					utils.DeleteFile(c, bucket, images[idx].MediaUrl)
				}
			}
		}()

		for idx := range payload.Images {
			if err := utils.UploadFile(c, bucket, images[idx].MediaUrl, payload.Images[idx]); err != nil {
				uploadErr = err
				s.repo.RollbackTxn(ctx)
				return fmt.Errorf("failed to upload files | %w", err)
			}
		}
	}

	c.JSON(http.StatusOK, post)

	return s.repo.CommitTxn(ctx)
}

func (s *Service) getTags(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to get tags |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong while getting tags")
		}
	}()

	params := c.Request.URL.Query()

	tags, err := s.repo.GetTags(c, params)
	if err != nil {
		return fmt.Errorf("failed to get tags | %w", err)
	}

	if len(tags) <= 0 {
		c.JSON(http.StatusNotFound, "No tags found")
		return nil
	}

	c.JSON(http.StatusOK, tags)

	return nil
}

func (s *Service) getCategories(c *gin.Context) (err error) {
	defer func() {
		if err != nil {
			log.Println("Failed to get categories |", err)
			c.JSON(http.StatusInternalServerError, "Something went wrong while getting tags")
		}
	}()

	categories, err := s.repo.GetCategories(c)
	if err != nil {
		return fmt.Errorf("failed to get categories | %w", err)
	}

	if len(categories) <= 0 {
		c.JSON(http.StatusNotFound, "No categories found")
		return nil
	}

	c.JSON(http.StatusOK, categories)

	return nil
}
