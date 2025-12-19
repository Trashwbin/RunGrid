//go:build !windows

package scanner

import (
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

	switch strings.ToLower(filepath.Ext(clean)) {
	case ".url", ".htm", ".html", ".mht", ".mhtml":
		return domain.ItemTypeURL
	case ".txt", ".md", ".markdown", ".pdf", ".rtf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".csv", ".log", ".chm":
		return domain.ItemTypeDoc
	default:
		return domain.ItemTypeApp
	}
}
