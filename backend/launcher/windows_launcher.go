//go:build windows

package launcher

import (
	"context"
	"os/exec"
	"syscall"
)

type WindowsLauncher struct{}

func NewDefaultLauncher() Launcher {
	return WindowsLauncher{}
}

func (WindowsLauncher) Open(ctx context.Context, target string) error {
	cmd := exec.CommandContext(ctx, "cmd", "/C", "start", "", target)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}
