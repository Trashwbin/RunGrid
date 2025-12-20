package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"rungrid/backend/domain"
	"rungrid/backend/icon"
	"rungrid/backend/storage"
)

type IconService struct {
	cache  *icon.Cache
	items  *ItemService
	mu     sync.Mutex
	busy   bool
	notify func()
}

func NewIconService(cache *icon.Cache, items *ItemService) *IconService {
	return &IconService{cache: cache, items: items}
}

func (s *IconService) EnsureForItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	if s.cache == nil {
		return item, icon.ErrUnsupported
	}
	if item.IconPath != "" || item.Path == "" {
		return item, nil
	}
	if item.Type == domain.ItemTypeURL {
		return item, nil
	}
	if !filepath.IsAbs(item.Path) {
		return item, nil
	}
	if strings.HasPrefix(strings.ToLower(item.Path), "http") {
		return item, nil
	}

	iconPath, err := s.cache.Ensure(ctx, item.Path, false)
	if err != nil {
		return item, err
	}

	if iconPath == "" {
		return item, nil
	}

	if err := s.items.SetIconPath(ctx, item.ID, iconPath); err != nil {
		return item, err
	}

	item.IconPath = iconPath
	return item, nil
}

func (s *IconService) RefreshItem(ctx context.Context, id string) (domain.Item, error) {
	if s.cache == nil {
		return domain.Item{}, icon.ErrUnsupported
	}
	if strings.TrimSpace(id) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	item, err := s.items.Get(ctx, id)
	if err != nil {
		return domain.Item{}, err
	}

	if item.Path == "" || item.Type == domain.ItemTypeURL {
		return domain.Item{}, storage.ErrInvalidInput
	}
	if !filepath.IsAbs(item.Path) {
		return domain.Item{}, storage.ErrInvalidInput
	}
	if strings.HasPrefix(strings.ToLower(item.Path), "http") {
		return domain.Item{}, storage.ErrInvalidInput
	}

	iconPath, err := s.cache.Ensure(ctx, item.Path, true)
	if err != nil {
		return domain.Item{}, err
	}
	if iconPath == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	if err := s.items.SetIconPath(ctx, item.ID, iconPath); err != nil {
		return domain.Item{}, err
	}

	item.IconPath = iconPath
	return item, nil
}

func (s *IconService) UpdateFromSource(ctx context.Context, id string, source string) (domain.Item, error) {
	if s.cache == nil {
		return domain.Item{}, icon.ErrUnsupported
	}
	if strings.TrimSpace(id) == "" || strings.TrimSpace(source) == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	item, err := s.items.Get(ctx, id)
	if err != nil {
		return domain.Item{}, err
	}

	iconPath, err := s.cache.Ensure(ctx, source, true)
	if err != nil {
		return domain.Item{}, err
	}
	if iconPath == "" {
		return domain.Item{}, storage.ErrInvalidInput
	}

	if err := s.items.SetIconPath(ctx, item.ID, iconPath); err != nil {
		return domain.Item{}, err
	}

	item.IconPath = iconPath
	return item, nil
}

func (s *IconService) PreviewFromSource(ctx context.Context, source string) (string, error) {
	if s.cache == nil {
		return "", icon.ErrUnsupported
	}
	if strings.TrimSpace(source) == "" {
		return "", storage.ErrInvalidInput
	}
	return s.cache.Ensure(ctx, source, true)
}

func (s *IconService) SyncMissing(ctx context.Context) (int, error) {
	return s.sync(ctx, false)
}

func (s *IconService) SyncMissingAsync(onDone func()) {
	if s.cache == nil {
		return
	}
	s.mu.Lock()
	if onDone != nil {
		s.notify = onDone
	}
	if s.busy {
		s.mu.Unlock()
		return
	}
	s.busy = true
	s.mu.Unlock()

	go func() {
		_, _ = s.sync(context.Background(), false)
		s.mu.Lock()
		notify := s.notify
		s.notify = nil
		s.busy = false
		s.mu.Unlock()
		if notify != nil {
			notify()
		}
	}()
}

func (s *IconService) RefreshAll(ctx context.Context) (int, error) {
	return s.sync(ctx, true)
}

func (s *IconService) sync(ctx context.Context, force bool) (int, error) {
	if s.cache == nil {
		return 0, icon.ErrUnsupported
	}

	items, err := s.items.List(ctx, storage.ItemFilter{})
	if err != nil {
		return 0, err
	}

	type iconTask struct {
		source string
		ids    []string
	}

	taskIndex := make(map[string]int)
	tasks := make([]iconTask, 0)
	for _, item := range items {
		if !force && item.IconPath != "" {
			continue
		}

		if item.Path == "" {
			continue
		}

		if item.Type == domain.ItemTypeURL {
			continue
		}
		if !filepath.IsAbs(item.Path) {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(item.Path))
		if key == "" {
			continue
		}

		index, ok := taskIndex[key]
		if !ok {
			index = len(tasks)
			taskIndex[key] = index
			tasks = append(tasks, iconTask{source: item.Path})
		}
		tasks[index].ids = append(tasks[index].ids, item.ID)
	}

	if len(tasks) == 0 {
		return 0, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var updated int64
	var once sync.Once
	var firstErr error
	setErr := func(err error) {
		if err == nil {
			return
		}
		once.Do(func() {
			firstErr = err
			cancel()
		})
	}

	const workerCount = 5
	taskCh := make(chan iconTask)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskCh:
					if !ok {
						return
					}
					iconPath, err := s.cache.Ensure(ctx, task.source, force)
					if err != nil {
						if errors.Is(err, icon.ErrUnsupported) {
							setErr(err)
							return
						}
						continue
					}
					if iconPath == "" {
						continue
					}

					for _, id := range task.ids {
						if ctx.Err() != nil {
							return
						}
						if err := s.items.SetIconPath(ctx, id, iconPath); err != nil {
							if errors.Is(err, storage.ErrNotFound) {
								continue
							}
							setErr(err)
							return
						}
						atomic.AddInt64(&updated, 1)
					}
				}
			}
		}()
	}

	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case taskCh <- task:
			}
		}
	}()

	wg.Wait()
	if firstErr != nil {
		return int(updated), firstErr
	}

	return int(updated), nil
}
