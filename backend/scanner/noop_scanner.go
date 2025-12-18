//go:build !windows

package scanner

import (
	"context"

	"rungrid/backend/domain"
)

type noopScanner struct{}

func NewDefaultScanner() Scanner {
	return noopScanner{}
}

func (noopScanner) Scan(_ context.Context) ([]domain.ItemInput, error) {
	return nil, ErrUnsupported
}
