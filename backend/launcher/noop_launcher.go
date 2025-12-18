//go:build !windows

package launcher

import "context"

type noopLauncher struct{}

func NewDefaultLauncher() Launcher {
	return noopLauncher{}
}

func (noopLauncher) Open(_ context.Context, _ string) error {
	return ErrUnsupported
}
