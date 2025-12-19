//go:build windows

package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"rungrid/backend/domain"
)

func ClassifyPath(path string) domain.ItemType {
	clean := strings.TrimSpace(path)
	if clean == "" {
		return domain.ItemTypeApp
	}

	lower := strings.ToLower(clean)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return domain.ItemTypeURL
	}
	if strings.HasPrefix(lower, "shell:appsfolder") || strings.HasPrefix(lower, "ms-settings:") {
		return domain.ItemTypeSystem
	}

	ext := strings.ToLower(filepath.Ext(clean))
	if _, ok := webExtensions[ext]; ok {
		return domain.ItemTypeURL
	}
	if _, ok := docExtensions[ext]; ok {
		return domain.ItemTypeDoc
	}

	if info, err := os.Stat(clean); err == nil && info.IsDir() {
		return domain.ItemTypeFolder
	}

	if isSystemBinaryPath(clean) {
		return domain.ItemTypeSystem
	}

	switch ext {
	case ".lnk", ".exe":
		return domain.ItemTypeApp
	default:
		return domain.ItemTypeApp
	}
}
