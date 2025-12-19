package storage

import (
	"sort"
	"time"

	"rungrid/backend/domain"
)

func SortItems(items []domain.Item) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Favorite != items[j].Favorite {
			return items[i].Favorite
		}
		iUsed := lastUsedAt(items[i])
		jUsed := lastUsedAt(items[j])
		if !iUsed.Equal(jUsed) {
			return iUsed.After(jUsed)
		}
		if items[i].LaunchCount != items[j].LaunchCount {
			return items[i].LaunchCount > items[j].LaunchCount
		}
		return compareItemName(items[i].Name, items[j].Name) < 0
	})
}

func lastUsedAt(item domain.Item) time.Time {
	if item.LastUsedAt != nil {
		return *item.LastUsedAt
	}
	return time.Time{}
}
