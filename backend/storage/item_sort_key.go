package storage

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"

	"rungrid/backend/domain"
)

var pinyinArgs = func() pinyin.Args {
	args := pinyin.NewArgs()
	args.Style = pinyin.FirstLetter
	return args
}()

func sortKeyForItem(item domain.Item) string {
	name := strings.TrimSpace(item.Name)
	if name == "" {
		name = trimExt(filepath.Base(item.Path))
	}

	key := buildNameSortKey(name)
	if key == "" {
		key = fallbackExeKey(item.Path)
	}
	if key == "" {
		key = strings.ToLower(strings.TrimSpace(name))
	}
	if key == "" {
		key = strings.ToLower(strings.TrimSpace(filepath.Base(item.Path)))
	}
	return key
}

func buildNameSortKey(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	if isASCIIOnly(trimmed) {
		return strings.ToLower(trimmed)
	}
	key := pinyinSortKey(trimmed)
	if key != "" {
		return key
	}
	return ""
}

func pinyinSortKey(name string) string {
	parts := pinyin.LazyPinyin(name, pinyinArgs)
	if len(parts) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		for _, r := range part {
			if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
				builder.WriteRune(unicode.ToLower(r))
				break
			}
		}
	}

	return builder.String()
}

func isASCIIOnly(value string) bool {
	for _, r := range value {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func fallbackExeKey(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".exe" {
		return ""
	}
	base := trimExt(filepath.Base(path))
	if base == "" {
		return ""
	}
	return strings.ToLower(base)
}

func trimExt(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}
	return strings.TrimSuffix(name, ext)
}
