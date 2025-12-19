//go:build !windows

package storage

import "strings"

func compareItemName(a, b string) int {
	return strings.Compare(strings.ToLower(a), strings.ToLower(b))
}
