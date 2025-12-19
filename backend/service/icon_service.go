package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"rungrid/backend/domain"
	"rungrid/backend/icon"
	"rungrid/backend/storage"
)

type IconService struct {
	cache *icon.Cache
	items *ItemService
}

func NewIconService(cache *icon.Cache, items *ItemService) *IconService {
	return &IconService{cache: cache, items: items}
}

func (s *IconService) EnsureForItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	if s.cache == nil {
		return item, icon.ErrUnsupported
	}
	if item.IconPath != "" || item.Path == "" {
		return item, nil
	}
	if item.Type == domain.ItemTypeURL {
		return item, nil
	}
	if !filepath.IsAbs(item.Path) {
		return item, nil
	}
	if strings.HasPrefix(strings.ToLower(item.Path), "http") {
		return item, nil
	}

	iconPath, err := s.cache.Ensure(ctx, item.Path, false)
	if err != nil {
		return item, err
	}

	if iconPath == "" {
		return item, nil
	}

	if err := s.items.SetIconPath(ctx, item.ID, iconPath); err != nil {
		return item, err
	}

	item.IconPath = iconPath
	return item, nil
}

func (s *IconService) RefreshItem(ctx context.Context, id string) (domain.Item, error) {
	if s.cache == nil {
		return domain.Item{}, icon.ErrUnsupported
	}
	if strings.TrimSpace(id) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	item, err := s.items.Get(ctx, id)
	if err != nil {
		return domain.Item{}, err
	}

	if item.Path == "" || item.Type == domain.ItemTypeURL {
		return domain.Item{}, storage.ErrInvalidInput
	}
	if !filepath.IsAbs(item.Path) {
		return domain.Item{}, storage.ErrInvalidInput
	}
	if strings.HasPrefix(strings.ToLower(item.Path), "http") {
		return domain.Item{}, storage.ErrInvalidInput
	}

	iconPath, err := s.cache.Ensure(ctx, item.Path, true)
	if err != nil {
		return domain.Item{}, err
	}
	if iconPath == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	if err := s.items.SetIconPath(ctx, item.ID, iconPath); err != nil {
		return domain.Item{}, err
	}

	item.IconPath = iconPath
	return item, nil
}

func (s *IconService) UpdateFromSource(ctx context.Context, id string, source string) (domain.Item, error) {
	if s.cache == nil {
		return domain.Item{}, icon.ErrUnsupported
	}
	if strings.TrimSpace(id) == "" || strings.TrimSpace(source) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	item, err := s.items.Get(ctx, id)
	if err != nil {
		return domain.Item{}, err
	}

	iconPath, err := s.cache.Ensure(ctx, source, true)
	if err != nil {
		return domain.Item{}, err
	}
	if iconPath == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	if err := s.items.SetIconPath(ctx, item.ID, iconPath); err != nil {
		return domain.Item{}, err
	}

	item.IconPath = iconPath
	return item, nil
}

func (s *IconService) PreviewFromSource(ctx context.Context, source string) (string, error) {
	if s.cache == nil {
		return "", icon.ErrUnsupported
	}
	if strings.TrimSpace(source) == "" {
		return "", storage.ErrInvalidInput
	}
	return s.cache.Ensure(ctx, source, true)
}

func (s *IconService) SyncMissing(ctx context.Context) (int, error) {
	return s.sync(ctx, false)
}

func (s *IconService) RefreshAll(ctx context.Context) (int, error) {
	return s.sync(ctx, true)
}

func (s *IconService) sync(ctx context.Context, force bool) (int, error) {
	if s.cache == nil {
		return 0, icon.ErrUnsupported
	}

	items, err := s.items.List(ctx, storage.ItemFilter{})
	if err != nil {
		return 0, err
	}

	updated := 0
	for _, item := range items {
		if !force && item.IconPath != "" {
			continue
		}

		if item.Path == "" {
			continue
		}

		if item.Type == domain.ItemTypeURL {
			continue
		}
		if !filepath.IsAbs(item.Path) {
			continue
		}

		iconPath, err := s.cache.Ensure(ctx, item.Path, force)
		if err != nil {
			if errors.Is(err, icon.ErrUnsupported) {
				return updated, err
			}
			continue
		}
		if iconPath == "" {
			continue
		}

		if err := s.items.SetIconPath(ctx, item.ID, iconPath); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				continue
			}
			return updated, err
		}

		item.IconPath = iconPath
		updated++
	}

	return updated, nil
}
