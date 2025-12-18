package scanner

import (
	"context"
	"errors"

	"rungrid/backend/domain"
)

var ErrUnsupported = errors.New("scanner not supported")

type Scanner interface {
	Scan(ctx context.Context) ([]domain.ItemInput, error)
}
