package storage

import (
	"sort"
	"time"

	"rungrid/backend/domain"
)

func SortItems(items []domain.Item) {
	type itemSortEntry struct {
		item domain.Item
		key  string
	}

	entries := make([]itemSortEntry, len(items))
	for i, item := range items {
		entries[i] = itemSortEntry{
			item: item,
			key:  sortKeyForItem(item),
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].item.Favorite != entries[j].item.Favorite {
			return entries[i].item.Favorite
		}
		iUsed := lastUsedAt(entries[i].item)
		jUsed := lastUsedAt(entries[j].item)
		if !iUsed.Equal(jUsed) {
			return iUsed.After(jUsed)
		}
		if entries[i].item.LaunchCount != entries[j].item.LaunchCount {
			return entries[i].item.LaunchCount > entries[j].item.LaunchCount
		}
		if entries[i].key != entries[j].key {
			return compareItemName(entries[i].key, entries[j].key) < 0
		}
		return compareItemName(entries[i].item.Name, entries[j].item.Name) < 0
	})

	for i := range items {
		items[i] = entries[i].item
	}
}

func lastUsedAt(item domain.Item) time.Time {
	if item.LastUsedAt != nil {
		return *item.LastUsedAt
	}
	return time.Time{}
}
