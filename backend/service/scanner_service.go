package service

import (
	"context"
	"errors"

	"rungrid/backend/domain"
	"rungrid/backend/scanner"
	"rungrid/backend/storage"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ScannerService struct {
	scanner scanner.Scanner
	items   *ItemService
	icons   *IconService
}

func NewScannerService(scanner scanner.Scanner, items *ItemService, icons *IconService) *ScannerService {
	return &ScannerService{scanner: scanner, items: items, icons: icons}
}

func (s *ScannerService) Scan(ctx context.Context) (domain.ScanResult, error) {
	return s.scan(ctx, nil)
}

func (s *ScannerService) ScanWithRoots(ctx context.Context, roots []string) (domain.ScanResult, error) {
	return s.scan(ctx, roots)
}

func (s *ScannerService) scan(ctx context.Context, roots []string) (domain.ScanResult, error) {
	if s.scanner == nil {
		return domain.ScanResult{}, scanner.ErrUnsupported
	}

	if setter, ok := s.scanner.(scanner.RootSetter); ok {
		if roots == nil {
			setter.SetRoots(scanner.DefaultRoots())
		} else {
			setter.SetRoots(roots)
		}
	}
	if reporter, ok := s.scanner.(scanner.ProgressReporter); ok {
		reporter.SetProgressReporter(func(progress scanner.ScanProgress) {
			runtime.EventsEmit(ctx, "scan:progress", progress)
		})
		defer reporter.SetProgressReporter(nil)
	}

	inputs, err := s.scanner.Scan(ctx)
	if err != nil {
		return domain.ScanResult{}, err
	}

	result := domain.ScanResult{Total: len(inputs)}
	for _, input := range inputs {
		if input.Path == "" || input.Name == "" {
			result.Skipped++
			continue
		}

		existing, err := s.items.GetByPath(ctx, input.Path)
		if err == nil {
			if existing.Type != input.Type && input.Type.IsValid() {
				_, updateErr := s.items.Update(ctx, domain.ItemUpdate{
					ID:       existing.ID,
					Type:     input.Type,
					Favorite: existing.Favorite,
					Hidden:   existing.Hidden,
				})
				if updateErr != nil && !errors.Is(updateErr, storage.ErrInvalidInput) {
					return result, updateErr
				}
			}
			result.Skipped++
			continue
		}
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return result, err
		}

		if _, err := s.items.Create(ctx, input); err != nil {
			if errors.Is(err, storage.ErrInvalidInput) {
				result.Skipped++
				continue
			}
			return result, err
		}

		result.Inserted++
	}

	if s.icons != nil {
		_, _ = s.icons.SyncMissing(ctx)
	}

	return result, nil
}
