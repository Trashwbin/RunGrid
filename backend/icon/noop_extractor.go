//go:build !windows

package icon

import "context"

type noopExtractor struct{}

func NewDefaultExtractor() Extractor {
	return noopExtractor{}
}

func (noopExtractor) Extract(_ context.Context, _ string, _ string) error {
	return ErrUnsupported
}
