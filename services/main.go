package services

import (
	"fmt"

	"github.com/katakeda/boardhop-api-service-go/repositories"
)

type Service struct {
	repo *repositories.Repository
}

func NewService(repo *repositories.Repository) (*Service, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository is required to start a new service")
	}

	return &Service{
		repo: repo,
	}, nil
}
