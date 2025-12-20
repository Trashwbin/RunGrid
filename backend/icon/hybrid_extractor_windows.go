//go:build windows

package icon

func NewHybridExtractor() Extractor {
	return HybridExtractor{
		primary:  NewNativeExtractor(),
		fallback: NewDefaultExtractor(),
	}
}
