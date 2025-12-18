package service

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"rungrid/backend/domain"
	"rungrid/backend/launcher"
	"rungrid/backend/storage"
)

type LauncherService struct {
	launcher launcher.Launcher
	items    *ItemService
}

func NewLauncherService(launcher launcher.Launcher, items *ItemService) *LauncherService {
	return &LauncherService{launcher: launcher, items: items}
}

func (s *LauncherService) LaunchItem(ctx context.Context, id string) (domain.Item, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	if s.launcher == nil {
		return domain.Item{}, launcher.ErrUnsupported
	}

	item, err := s.items.Get(ctx, id)
	if err != nil {
		return domain.Item{}, err
	}

	if err := validateLaunchTarget(item); err != nil {
		return domain.Item{}, err
	}

	if err := s.launcher.Open(ctx, item.Path); err != nil {
		return domain.Item{}, err
	}

	return s.items.RecordLaunch(ctx, id)
}

func validateLaunchTarget(item domain.Item) error {
	target := strings.TrimSpace(item.Path)
	if target == "" {
		return storage.ErrInvalidInput
	}

	if isWebURL(target) {
		if !isAllowedScheme(target) {
			return storage.ErrInvalidInput
		}
		return nil
	}

	if isUNCPath(target) {
		return storage.ErrInvalidInput
	}

	if !filepath.IsAbs(target) {
		return storage.ErrInvalidInput
	}

	info, err := os.Stat(target)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return storage.ErrInvalidInput
		}
		return err
	}

	if info.IsDir() {
		return nil
	}

	if info.Mode().IsRegular() {
		return nil
	}

	return storage.ErrInvalidInput
}

func isWebURL(target string) bool {
	lower := strings.ToLower(strings.TrimSpace(target))
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func isAllowedScheme(target string) bool {
	parsed, err := url.Parse(target)
	if err != nil {
		return false
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
		return true
	default:
		return false
	}
}

func isUNCPath(target string) bool {
	return strings.HasPrefix(target, `\\`) || strings.HasPrefix(target, "//")
}
