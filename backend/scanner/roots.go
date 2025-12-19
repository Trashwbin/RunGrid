package scanner

import (
	"path/filepath"
	"strings"
)

func NormalizeRoots(roots []string) []string {
	if len(roots) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(roots))
	result := make([]string, 0, len(roots))

	for _, root := range roots {
		trimmed := strings.TrimSpace(root)
		if trimmed == "" {
			continue
		}
		cleaned := filepath.Clean(trimmed)
		if cleaned == "." {
			continue
		}
		key := strings.ToLower(cleaned)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, cleaned)
	}

	return result
}
