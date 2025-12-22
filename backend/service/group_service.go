package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"rungrid/backend/domain"
	"rungrid/backend/storage"
)

type GroupService struct {
	repo storage.GroupRepository
}

func NewGroupService(repo storage.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) List(ctx context.Context) ([]domain.Group, error) {
	return s.repo.List(ctx)
}

func (s *GroupService) Create(ctx context.Context, input domain.GroupInput) (domain.Group, error) {
	if strings.TrimSpace(input.Name) == "" {
		return domain.Group{}, storage.ErrInvalidInput
	}

	category, err := normalizeGroupCategory(input.Category)
	if err != nil {
		return domain.Group{}, err
	}

	group := domain.Group{
		ID:       uuid.NewString(),
		Name:     strings.TrimSpace(input.Name),
		Order:    input.Order,
		Color:    strings.TrimSpace(input.Color),
		Category: category,
	}

	return s.repo.Create(ctx, group)
}

func (s *GroupService) Update(ctx context.Context, group domain.Group) (domain.Group, error) {
	if strings.TrimSpace(group.ID) == "" || strings.TrimSpace(group.Name) == "" {
		return domain.Group{}, storage.ErrInvalidInput
	}
	group.Name = strings.TrimSpace(group.Name)
	group.Color = strings.TrimSpace(group.Color)
	if strings.TrimSpace(group.Category) == "" {
		existing, err := s.repo.Get(ctx, group.ID)
		if err != nil {
			return domain.Group{}, err
		}
		group.Category = existing.Category
	} else {
		category, err := normalizeGroupCategory(group.Category)
		if err != nil {
			return domain.Group{}, err
		}
		group.Category = category
	}
	return s.repo.Update(ctx, group)
}

func (s *GroupService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func normalizeGroupCategory(category string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(category))
	if normalized == "" {
		return "app", nil
	}

	switch normalized {
	case "app", "system", "doc", "folder", "url":
		return normalized, nil
	default:
		return "", storage.ErrInvalidInput
	}
}
