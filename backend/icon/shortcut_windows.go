//go:build windows

package icon

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type shortcutIconInfo struct {
	source string
	index  int
}

func resolveShortcutIcon(path string) (shortcutIconInfo, error) {
	info := shortcutIconInfo{}
	if strings.TrimSpace(path) == "" {
		return info, nil
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	const comSFalse = uintptr(0x00000001)
	if initErr := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); initErr != nil {
		if oleErr, ok := initErr.(*ole.OleError); !ok || oleErr.Code() != comSFalse {
			return info, initErr
		}
	}
	defer ole.CoUninitialize()

	shellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return info, err
	}
	defer shellObject.Release()

	shellDispatch, err := shellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return info, err
	}
	defer shellDispatch.Release()

	shortcutVariant, err := oleutil.CallMethod(shellDispatch, "CreateShortcut", path)
	if err != nil {
		return info, err
	}
	shortcut := shortcutVariant.ToIDispatch()
	if shortcut == nil {
		shortcutVariant.Clear()
		return info, nil
	}
	defer shortcutVariant.Clear()

	iconLocation, _ := oleutil.GetProperty(shortcut, "IconLocation")
	targetPath, _ := oleutil.GetProperty(shortcut, "TargetPath")

	iconLocationValue := ""
	if iconLocation != nil {
		iconLocationValue = strings.TrimSpace(iconLocation.ToString())
		iconLocation.Clear()
	}

	targetValue := ""
	if targetPath != nil {
		targetValue = strings.TrimSpace(targetPath.ToString())
		targetPath.Clear()
	}

	iconCandidate, iconIndex := parseIconLocation(iconLocationValue)
	if strings.EqualFold(filepath.Ext(iconCandidate), ".lnk") {
		iconCandidate = ""
		iconIndex = 0
	}

	if iconCandidate == "" && targetValue != "" {
		iconCandidate = targetValue
		iconIndex = 0
	}

	resolved := resolveIconCandidate(path, iconCandidate)
	if resolved != "" {
		info.source = resolved
		info.index = iconIndex
	}

	return info, nil
}

func parseIconLocation(value string) (string, int) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", 0
	}
	path := value
	index := 0
	if comma := strings.LastIndex(value, ","); comma > -1 {
		path = strings.TrimSpace(value[:comma])
		idx := strings.TrimSpace(value[comma+1:])
		if parsed, err := strconv.Atoi(idx); err == nil {
			index = parsed
		}
	}
	path = strings.Trim(path, "\"'")
	return path, index
}

func resolveIconCandidate(shortcutPath, candidate string) string {
	if strings.TrimSpace(candidate) == "" {
		return ""
	}
	candidate = strings.Trim(candidate, "\"'")
	candidate = expandEnvVars(candidate)
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(filepath.Dir(shortcutPath), candidate)
	}
	candidate = filepath.Clean(candidate)

	if fileExists(candidate) {
		return candidate
	}

	lower := strings.ToLower(candidate)
	for _, ext := range []string{".exe", ".dll", ".cpl", ".msc", ".ico"} {
		if idx := strings.Index(lower, ext); idx >= 0 {
			fixed := candidate[:idx+len(ext)]
			if fileExists(fixed) {
				return fixed
			}
		}
	}

	return ""
}

func expandEnvVars(value string) string {
	result := value
	for {
		start := strings.Index(result, "%")
		if start == -1 {
			break
		}
		end := strings.Index(result[start+1:], "%")
		if end == -1 {
			break
		}
		end = start + 1 + end
		key := result[start+1 : end]
		repl := os.Getenv(key)
		result = result[:start] + repl + result[end+1:]
	}
	return result
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
