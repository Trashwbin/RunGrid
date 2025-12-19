//go:build windows

package scanner

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	type dedupeCandidate struct {
		item      domain.ItemInput
		timestamp time.Time
	}

	candidates := map[string]dedupeCandidate{}
	keys := []string{}
	var resolver *shortcutResolver
	var resolverErr error
	defer func() {
		if resolver != nil {
			resolver.Close()
		}
	}()

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

			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			name = strings.TrimSpace(name)
			if ext == ".lnk" {
				if displayName, ok := lookupDisplayName(path); ok {
					name = displayName
				}
			}
			if name == "" {
				name = entry.Name()
			}

			timestamp := time.Time{}
			if latest, ok := latestFileTimestamp(path); ok {
				timestamp = latest
			} else if info, infoErr := entry.Info(); infoErr == nil {
				timestamp = info.ModTime()
			}

			dedupeKey := strings.ToLower(path)
			if ext == ".lnk" {
				if resolver == nil && resolverErr == nil {
					resolver, resolverErr = newShortcutResolver()
				}
				if resolver != nil {
					target, args, err := resolver.Resolve(path)
					if err == nil {
						itemType = classifyShortcutTarget(path, target, args, itemType)
						if shortcutKey := shortcutDedupeKey(target, args); shortcutKey != "" {
							dedupeKey = shortcutKey
						}
					}
				}
			} else if ext == ".exe" {
				if isSystemBinaryPath(path) {
					itemType = domain.ItemTypeSystem
				}
			}

			item := domain.ItemInput{
				Name:     name,
				Path:     path,
				Type:     itemType,
				IconPath: "",
				GroupID:  "",
				Tags:     nil,
				Favorite: false,
				Hidden:   false,
			}

			if existing, ok := candidates[dedupeKey]; ok {
				if timestamp.After(existing.timestamp) {
					candidates[dedupeKey] = dedupeCandidate{item: item, timestamp: timestamp}
				}
				return nil
			}

			candidates[dedupeKey] = dedupeCandidate{item: item, timestamp: timestamp}
			keys = append(keys, dedupeKey)

			return nil
		})
		if walkErr != nil {
			return nil, walkErr
		}
	}

	items := make([]domain.ItemInput, 0, len(keys))
	for _, key := range keys {
		candidate, ok := candidates[key]
		if !ok {
			continue
		}
		items = append(items, candidate.item)
	}

	return items, nil
}

func mapExtensionType(ext string) (domain.ItemType, bool) {
	switch ext {
	case ".lnk", ".exe":
		return domain.ItemTypeApp, true
	case ".url", ".htm", ".html", ".mht", ".mhtml":
		return domain.ItemTypeURL, true
	default:
		return "", false
	}
}
