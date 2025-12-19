//go:build windows

package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"rungrid/backend/domain"
)

var docExtensions = map[string]struct{}{
	".txt":      {},
	".md":       {},
	".markdown": {},
	".pdf":      {},
	".rtf":      {},
	".doc":      {},
	".docx":     {},
	".xls":      {},
	".xlsx":     {},
	".ppt":      {},
	".pptx":     {},
	".csv":      {},
	".log":      {},
}

var webExtensions = map[string]struct{}{
	".url":   {},
	".htm":   {},
	".html":  {},
	".mht":   {},
	".mhtml": {},
}

func classifyShortcutTarget(source, target, args string, fallback domain.ItemType) domain.ItemType {
	source = normalizePath(source)
	target = strings.TrimSpace(target)
	args = strings.TrimSpace(args)

	if hasWebURL(target) || hasWebURL(args) {
		return domain.ItemTypeURL
	}

	ext := strings.ToLower(filepath.Ext(target))
	if _, ok := webExtensions[ext]; ok {
		return domain.ItemTypeURL
	}
	if _, ok := docExtensions[ext]; ok {
		return domain.ItemTypeDoc
	}

	if isSystemShortcutSource(source) || isSystemTarget(target, args) {
		return domain.ItemTypeSystem
	}

	if target != "" {
		if info, err := os.Stat(target); err == nil && info.IsDir() {
			return domain.ItemTypeFolder
		}
	}

	if target != "" {
		return domain.ItemTypeApp
	}

	return fallback
}

func hasWebURL(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "http://") || strings.Contains(lower, "https://")
}

func isSystemTarget(target, args string) bool {
	target = normalizePath(target)
	argsLower := strings.ToLower(args)

	if strings.Contains(argsLower, "shell:appsfolder") || strings.Contains(argsLower, "ms-settings:") {
		return true
	}

	if target == "" {
		return false
	}

	return isSystemBinaryPath(target)
}

func isSystemBinaryPath(path string) bool {
	clean := normalizePath(path)
	if clean == "" {
		return false
	}

	for _, root := range systemRoots() {
		if hasPathPrefix(clean, root) {
			return true
		}
	}
	return false
}

func isSystemShortcutSource(path string) bool {
	if path == "" {
		return false
	}

	for _, root := range startMenuRoots() {
		if !hasPathPrefix(path, root) {
			continue
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			continue
		}
		lowerRel := normalizePath(rel)
		for _, folder := range systemShortcutFolders() {
			if strings.HasPrefix(lowerRel, folder+string(os.PathSeparator)) || lowerRel == folder {
				return true
			}
		}
	}

	return false
}

func systemRoots() []string {
	roots := []string{}
	if windir := os.Getenv("WINDIR"); windir != "" {
		roots = append(roots,
			filepath.Join(windir, "System32"),
			filepath.Join(windir, "SysWOW64"),
			filepath.Join(windir, "SystemApps"),
			filepath.Join(windir, "Explorer.exe"),
		)
	}
	if systemRoot := os.Getenv("SystemRoot"); systemRoot != "" {
		roots = append(roots,
			filepath.Join(systemRoot, "System32"),
			filepath.Join(systemRoot, "SysWOW64"),
			filepath.Join(systemRoot, "SystemApps"),
			filepath.Join(systemRoot, "Explorer.exe"),
		)
	}

	normalized := make([]string, 0, len(roots))
	seen := make(map[string]struct{})
	for _, root := range roots {
		root = normalizePath(root)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		normalized = append(normalized, root)
	}
	return normalized
}

func startMenuRoots() []string {
	roots := []string{}
	if appData := os.Getenv("APPDATA"); appData != "" {
		roots = append(roots, filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs"))
	}
	if programData := os.Getenv("PROGRAMDATA"); programData != "" {
		roots = append(roots, filepath.Join(programData, "Microsoft", "Windows", "Start Menu", "Programs"))
	}

	normalized := make([]string, 0, len(roots))
	seen := make(map[string]struct{})
	for _, root := range roots {
		root = normalizePath(root)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		normalized = append(normalized, root)
	}
	return normalized
}

func systemShortcutFolders() []string {
	return []string{
		normalizePath("Administrative Tools"),
		normalizePath("System Tools"),
		normalizePath("Windows Tools"),
	}
}

func normalizePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = filepath.Clean(value)
	value = strings.ReplaceAll(value, "/", string(os.PathSeparator))
	return strings.ToLower(value)
}

func hasPathPrefix(path, prefix string) bool {
	if path == prefix {
		return true
	}
	prefixWithSep := prefix + string(os.PathSeparator)
	return strings.HasPrefix(path, prefixWithSep)
}
