package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"rungrid/backend/domain"
	"rungrid/backend/storage"
)

type GroupRepository struct {
	mu     sync.RWMutex
	groups map[string]domain.Group
}

func NewGroupRepository() *GroupRepository {
	return &GroupRepository{groups: make(map[string]domain.Group)}
}

func (r *GroupRepository) List(_ context.Context) ([]domain.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	groups := make([]domain.Group, 0, len(r.groups))
	for _, group := range r.groups {
		groups = append(groups, group)
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Order == groups[j].Order {
			return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
		}
		return groups[i].Order < groups[j].Order
	})

	return groups, nil
}

func (r *GroupRepository) Get(_ context.Context, id string) (domain.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	group, ok := r.groups[id]
	if !ok {
		return domain.Group{}, storage.ErrNotFound
	}

	return group, nil
}

func (r *GroupRepository) Create(_ context.Context, group domain.Group) (domain.Group, error) {
	if strings.TrimSpace(group.ID) == "" || strings.TrimSpace(group.Name) == "" {
		return domain.Group{}, storage.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[group.ID]; exists {
		return domain.Group{}, storage.ErrInvalidInput
	}

	r.groups[group.ID] = group
	return group, nil
}

func (r *GroupRepository) Update(_ context.Context, group domain.Group) (domain.Group, error) {
	if strings.TrimSpace(group.ID) == "" || strings.TrimSpace(group.Name) == "" {
		return domain.Group{}, storage.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[group.ID]; !exists {
		return domain.Group{}, storage.ErrNotFound
	}

	r.groups[group.ID] = group
	return group, nil
}

func (r *GroupRepository) Delete(_ context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return storage.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[id]; !exists {
		return storage.ErrNotFound
	}

	delete(r.groups, id)
	return nil
}
