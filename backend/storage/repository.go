package storage

import (
	"context"
	"time"

	"rungrid/backend/domain"
)

type ItemFilter struct {
	GroupID string
	Query   string
}

type ItemRepository interface {
	List(ctx context.Context, filter ItemFilter) ([]domain.Item, error)
	Get(ctx context.Context, id string) (domain.Item, error)
	GetByPath(ctx context.Context, path string) (domain.Item, error)
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	Update(ctx context.Context, item domain.Item) (domain.Item, error)
	Delete(ctx context.Context, id string) error
	IncrementLaunch(ctx context.Context, id string, usedAt time.Time) (domain.Item, error)
}

type GroupRepository interface {
	List(ctx context.Context) ([]domain.Group, error)
	Get(ctx context.Context, id string) (domain.Group, error)
	Create(ctx context.Context, group domain.Group) (domain.Group, error)
	Update(ctx context.Context, group domain.Group) (domain.Group, error)
	Delete(ctx context.Context, id string) error
}
