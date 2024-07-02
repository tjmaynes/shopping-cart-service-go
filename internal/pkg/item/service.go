package item

import (
	"context"
	"github.com/google/uuid"
)

// Service ..
type Service interface {
	GetItems(ctx context.Context, page int64, pageSize int64) ([]Item, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (Item, error)
	AddItem(
		ctx context.Context,
		item *ItemDTO,
	) (Item, error)
	UpdateItem(
		ctx context.Context,
		item *Item,
	) (Item, ServiceError)
	RemoveItem(ctx context.Context, id uuid.UUID) (uuid.UUID, ServiceError)
}

// NewService ..
func NewService(repository Repository) Service {
	return &service{
		Repository: repository,
	}
}

type service struct {
	Repository Repository
}

// GetItems ..
func (s *service) GetItems(ctx context.Context, page int64, pageSize int64) ([]Item, error) {
	return s.Repository.GetItems(ctx, page, pageSize)
}

// GetItemByID ..
func (s *service) GetItemByID(ctx context.Context, id uuid.UUID) (Item, error) {
	return s.Repository.GetItemByID(ctx, id)
}

// AddItem ..
func (s *service) AddItem(ctx context.Context, item *ItemDTO) (Item, error) {
	err := item.Validate()
	if err != nil {
		return Item{}, err
	}

	return s.Repository.AddItem(ctx, item.Name, item.Price, item.Manufacturer)
}

// UpdateItem ..
func (s *service) UpdateItem(ctx context.Context, item *Item) (Item, ServiceError) {
	err := item.Validate()
	if err != nil {
		return Item{}, CreateServiceError(err.Error(), InvalidItem)
	}

	result, err := s.Repository.UpdateItem(ctx, item)
	if err != nil {
		return Item{}, CreateServiceError(err.Error(), UnknownException)
	}

	return result, nil
}

// RemoveItem ..
func (s *service) RemoveItem(ctx context.Context, id uuid.UUID) (uuid.UUID, ServiceError) {
	result, err := s.Repository.RemoveItem(ctx, id)
	if err != nil {
		return id, CreateServiceError(err.Error(), UnknownException)
	}

	return result, nil
}
