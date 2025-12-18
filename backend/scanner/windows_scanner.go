//go:build windows

package scanner

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"rungrid/backend/domain"
)

type WindowsScanner struct {
	Roots []string
}

func NewDefaultScanner() Scanner {
	return &WindowsScanner{Roots: DefaultRoots()}
}

func DefaultRoots() []string {
	roots := []string{}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		roots = append(roots, filepath.Join(home, "Desktop"))
	}

	if appData := os.Getenv("APPDATA"); appData != "" {
		roots = append(roots, filepath.Join(appData, "Microsoft", "Windows", "Start Menu"))
	}

	if programData := os.Getenv("PROGRAMDATA"); programData != "" {
		roots = append(roots, filepath.Join(programData, "Microsoft", "Windows", "Start Menu"))
	}

	return roots
}

func (s *WindowsScanner) Scan(ctx context.Context) ([]domain.ItemInput, error) {
	items := []domain.ItemInput{}
	seen := map[string]struct{}{}

	for _, root := range s.Roots {
		if root == "" {
			continue
		}
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			continue
		}

		walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if entry.IsDir() {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(entry.Name()))
			itemType, ok := mapExtensionType(ext)
			if !ok {
				return nil
			}

			if _, exists := seen[path]; exists {
				return nil
			}
			seen[path] = struct{}{}

			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			name = strings.TrimSpace(name)
			if name == "" {
				name = entry.Name()
			}

			items = append(items, domain.ItemInput{
				Name:     name,
				Path:     path,
				Type:     itemType,
				IconPath: "",
				GroupID:  "",
				Tags:     nil,
				Favorite: false,
				Hidden:   false,
			})

			return nil
		})
		if walkErr != nil {
			return nil, walkErr
		}
	}

	return items, nil
}

func mapExtensionType(ext string) (domain.ItemType, bool) {
	switch ext {
	case ".lnk", ".exe":
		return domain.ItemTypeApp, true
	case ".url":
		return domain.ItemTypeURL, true
	default:
		return "", false
	}
}
