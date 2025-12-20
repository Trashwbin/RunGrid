package icon

import "context"

type HybridExtractor struct {
	primary  Extractor
	fallback Extractor
}

func (h HybridExtractor) Extract(ctx context.Context, source string, dest string) error {
	var firstErr error
	if h.primary != nil {
		if err := h.primary.Extract(ctx, source, dest); err == nil {
			return nil
		} else {
			firstErr = err
		}
	}
	if h.fallback != nil {
		if err := h.fallback.Extract(ctx, source, dest); err == nil {
			return nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return firstErr
	}
	return ErrUnsupported
}
