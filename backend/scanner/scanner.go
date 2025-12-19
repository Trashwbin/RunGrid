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

type RootSetter interface {
	SetRoots(roots []string)
}

type ProgressFunc func(ScanProgress)

type ProgressReporter interface {
	SetProgressReporter(fn ProgressFunc)
}

type ScanProgress struct {
	Root      string `json:"root"`
	Path      string `json:"path"`
	RootIndex int    `json:"rootIndex"`
	RootTotal int    `json:"rootTotal"`
	Scanned   int    `json:"scanned"`
	Percent   int    `json:"percent"`
}
