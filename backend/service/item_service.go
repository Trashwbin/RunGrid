package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"rungrid/backend/domain"
	"rungrid/backend/storage"
)

type ItemService struct {
	repo storage.ItemRepository
}

func NewItemService(repo storage.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) List(ctx context.Context, filter storage.ItemFilter) ([]domain.Item, error) {
	return s.repo.List(ctx, filter)
}

func (s *ItemService) GetByPath(ctx context.Context, path string) (domain.Item, error) {
	clean := strings.TrimSpace(path)
	if clean == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}
	return s.repo.GetByPath(ctx, clean)
}

func (s *ItemService) Create(ctx context.Context, input domain.ItemInput) (domain.Item, error) {
	if err := validateItemInput(input); err != nil {
		return domain.Item{}, err
	}

	item := domain.Item{
		ID:       uuid.NewString(),
		Name:     strings.TrimSpace(input.Name),
		Path:     strings.TrimSpace(input.Path),
		Type:     input.Type,
		IconPath: strings.TrimSpace(input.IconPath),
		GroupID:  strings.TrimSpace(input.GroupID),
		Tags:     dedupeTags(input.Tags),
		Favorite: input.Favorite,
		Hidden:   input.Hidden,
	}

	return s.repo.Create(ctx, item)
}

func (s *ItemService) Update(ctx context.Context, input domain.ItemUpdate) (domain.Item, error) {
	if strings.TrimSpace(input.ID) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	if input.Type != "" && !input.Type.IsValid() {
		return domain.Item{}, storage.ErrInvalidInput
	}

	current, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return domain.Item{}, err
	}

	updated := current
	if strings.TrimSpace(input.Name) != "" {
		updated.Name = strings.TrimSpace(input.Name)
	}
	if strings.TrimSpace(input.Path) != "" {
		updated.Path = strings.TrimSpace(input.Path)
	}
	if input.Type != "" {
		updated.Type = input.Type
	}
	if strings.TrimSpace(input.IconPath) != "" {
		updated.IconPath = strings.TrimSpace(input.IconPath)
	}
	if strings.TrimSpace(input.GroupID) != "" {
		updated.GroupID = strings.TrimSpace(input.GroupID)
	}
	if input.Tags != nil {
		updated.Tags = dedupeTags(input.Tags)
	}
	updated.Favorite = input.Favorite
	updated.Hidden = input.Hidden

	return s.repo.Update(ctx, updated)
}

func (s *ItemService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ItemService) RecordLaunch(ctx context.Context, id string) (domain.Item, error) {
	return s.repo.IncrementLaunch(ctx, id, time.Now())
}

func validateItemInput(input domain.ItemInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return storage.ErrInvalidInput
	}
	if strings.TrimSpace(input.Path) == "" {
		return storage.ErrInvalidInput
	}
	if !input.Type.IsValid() {
		return storage.ErrInvalidInput
	}
	return nil
}

func dedupeTags(tags []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		clean := strings.TrimSpace(tag)
		if clean == "" {
			continue
		}
		key := strings.ToLower(clean)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, clean)
	}
	return result
}
