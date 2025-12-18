package service

import (
	"context"
	"errors"

	"rungrid/backend/domain"
	"rungrid/backend/scanner"
	"rungrid/backend/storage"
)

type ScannerService struct {
	scanner scanner.Scanner
	items   *ItemService
}

func NewScannerService(scanner scanner.Scanner, items *ItemService) *ScannerService {
	return &ScannerService{scanner: scanner, items: items}
}

func (s *ScannerService) Scan(ctx context.Context) (domain.ScanResult, error) {
	if s.scanner == nil {
		return domain.ScanResult{}, scanner.ErrUnsupported
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

		_, err := s.items.GetByPath(ctx, input.Path)
		if err == nil {
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

	return result, nil
}
