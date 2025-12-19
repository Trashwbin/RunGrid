//go:build windows

package scanner

import (
	"path/filepath"
	"strings"
)

func isUninstallerEntry(name, sourcePath, targetPath, args string) bool {
	if hasUninstallKeyword(name) {
		return true
	}
	if hasUninstallKeyword(sourcePath) {
		return true
	}
	if hasUninstallKeyword(targetPath) {
		return true
	}
	if hasUninstallKeyword(args) {
		return true
	}

	targetBase := strings.ToLower(filepath.Base(strings.Trim(targetPath, "\"'")))
	if targetBase == "uninstall.exe" || targetBase == "uninstaller.exe" {
		return true
	}
	if strings.HasPrefix(targetBase, "unins") && strings.HasSuffix(targetBase, ".exe") {
		return true
	}

	return false
}

func hasUninstallKeyword(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "uninstall") ||
		strings.Contains(lower, "uninstaller") ||
		strings.Contains(lower, "unins000") ||
		strings.Contains(lower, "unins") && strings.Contains(lower, ".exe") ||
		strings.Contains(value, "卸载")
}
