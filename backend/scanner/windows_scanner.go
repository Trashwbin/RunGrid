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
	Roots    []string
	progress ProgressFunc
}

func NewDefaultScanner() Scanner {
	return &WindowsScanner{Roots: NormalizeRoots(DefaultRoots())}
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

func (s *WindowsScanner) SetRoots(roots []string) {
	s.Roots = NormalizeRoots(roots)
}

func (s *WindowsScanner) SetProgressReporter(fn ProgressFunc) {
	s.progress = fn
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

	roots := NormalizeRoots(s.Roots)
	if len(roots) == 0 {
		roots = NormalizeRoots(DefaultRoots())
	}
	totalRoots := len(roots)
	scanned := 0
	lastEmit := time.Now().Add(-time.Second)

	for index, root := range roots {
		if root == "" {
			continue
		}
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			continue
		}

		s.emitProgress(ScanProgress{
			Root:      root,
			Path:      root,
			RootIndex: index + 1,
			RootTotal: totalRoots,
			Scanned:   scanned,
			Percent:   rootPercent(index, totalRoots),
		})

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

			scanned++
			if time.Since(lastEmit) >= 200*time.Millisecond {
				lastEmit = time.Now()
				s.emitProgress(ScanProgress{
					Root:      root,
					Path:      path,
					RootIndex: index + 1,
					RootTotal: totalRoots,
					Scanned:   scanned,
					Percent:   -1,
				})
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
						if isUninstallerEntry(name, path, target, args) {
							return nil
						}
						itemType = classifyShortcutTarget(path, target, args, itemType)
						if shortcutKey := shortcutDedupeKey(target, args); shortcutKey != "" {
							dedupeKey = shortcutKey
						}
					}
				}
			} else if ext == ".exe" {
				if isUninstallerEntry(name, path, path, "") {
					return nil
				}
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

		s.emitProgress(ScanProgress{
			Root:      root,
			Path:      root,
			RootIndex: index + 1,
			RootTotal: totalRoots,
			Scanned:   scanned,
			Percent:   rootPercent(index+1, totalRoots),
		})
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

func (s *WindowsScanner) emitProgress(progress ScanProgress) {
	if s.progress == nil {
		return
	}
	s.progress(progress)
}

func rootPercent(index, total int) int {
	if total <= 0 {
		return 0
	}
	if index < 0 {
		index = 0
	}
	if index > total {
		index = total
	}
	return index * 100 / total
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
