//go:build !windows

package icon

func NewHybridExtractor() Extractor {
	return NewDefaultExtractor()
}
