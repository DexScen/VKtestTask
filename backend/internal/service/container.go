package service

import (
	"context"

	"github.com/DexScen/VKtestTask/backend/internal/domain"
)

type ContainersRepository interface {
	GetContainers(ctx context.Context, list *domain.ListContainer) error
	PostContainers(ctx context.Context, list *domain.ListContainer) error
}

type Containers struct {
	repo ContainersRepository
}

func NewContainers(repo ContainersRepository) *Containers {
	return &Containers{
		repo: repo,
	}
}

func (c *Containers) GetContainers(ctx context.Context, list *domain.ListContainer) error {
	return c.repo.GetContainers(ctx, list)
}

func (c *Containers) PostContainers(ctx context.Context, list *domain.ListContainer) error {
	return c.repo.PostContainers(ctx, list)
}
