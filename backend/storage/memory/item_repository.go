package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"rungrid/backend/domain"
	"rungrid/backend/storage"
)

type ItemRepository struct {
	mu    sync.RWMutex
	items map[string]domain.Item
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{items: make(map[string]domain.Item)}
}

func (r *ItemRepository) List(_ context.Context, filter storage.ItemFilter) ([]domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.Item, 0, len(r.items))
	query := strings.ToLower(strings.TrimSpace(filter.Query))

	for _, item := range r.items {
		if filter.GroupID != "" && filter.GroupID != "all" && item.GroupID != filter.GroupID {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(item.Name), query) {
			continue
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

func (r *ItemRepository) Get(_ context.Context, id string) (domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return domain.Item{}, storage.ErrNotFound
	}

	return item, nil
}

func (r *ItemRepository) GetByPath(_ context.Context, path string) (domain.Item, error) {
	if strings.TrimSpace(path) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if strings.EqualFold(item.Path, path) {
			return item, nil
		}
	}

	return domain.Item{}, storage.ErrNotFound
}

func (r *ItemRepository) Create(_ context.Context, item domain.Item) (domain.Item, error) {
	if strings.TrimSpace(item.ID) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; exists {
		return domain.Item{}, storage.ErrInvalidInput
	}

	r.items[item.ID] = item
	return item, nil
}

func (r *ItemRepository) Update(_ context.Context, item domain.Item) (domain.Item, error) {
	if strings.TrimSpace(item.ID) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return domain.Item{}, storage.ErrNotFound
	}

	r.items[item.ID] = item
	return item, nil
}

func (r *ItemRepository) Delete(_ context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return storage.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return storage.ErrNotFound
	}

	delete(r.items, id)
	return nil
}

func (r *ItemRepository) IncrementLaunch(_ context.Context, id string, usedAt time.Time) (domain.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.items[id]
	if !exists {
		return domain.Item{}, storage.ErrNotFound
	}

	item.LaunchCount++
	item.LastUsedAt = &usedAt

	r.items[id] = item
	return item, nil
}
