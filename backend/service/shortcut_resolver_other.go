//go:build !windows

package service

import "context"

func resolveShortcutTarget(_ context.Context, _ string) (string, error) {
	return "", nil
}
