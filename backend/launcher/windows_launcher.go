//go:build windows

package launcher

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

type WindowsLauncher struct{}

func NewDefaultLauncher() Launcher {
	return WindowsLauncher{}
}

func (WindowsLauncher) Open(ctx context.Context, target string) error {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" {
		return ErrUnsupported
	}

	// Try native ShellExecute first (fast and works for exe/lnk/url/dir)
	if err := shellExecute(trimmed); err == nil {
		return nil
	}

	// Fallback to cmd start with explicit quoting
	return startWithCmd(ctx, trimmed)
}

const (
	shellExecuteShow = 1 // SW_SHOWNORMAL
)

var (
	shell32             = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteW   = shell32.NewProc("ShellExecuteW")
	shellExecuteOpen, _ = syscall.UTF16PtrFromString("open")
)

func shellExecute(target string) error {
	ptr, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return err
	}

	// HWND=0, verb=open, parameters/cwd nil
	r, _, callErr := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(shellExecuteOpen)),
		uintptr(unsafe.Pointer(ptr)),
		0,
		0,
		uintptr(shellExecuteShow),
	)

	// ShellExecute returns >32 on success
	if r <= 32 {
		if callErr != nil && callErr != syscall.Errno(0) {
			return callErr
		}
		return errors.New("ShellExecute failed")
	}

	return nil
}

func startWithCmd(ctx context.Context, target string) error {
	quoted := quoteForCmd(target)
	cmd := exec.CommandContext(ctx, "cmd", "/C", "start", "", quoted)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}

func quoteForCmd(value string) string {
	trimmed := strings.Trim(value, "\"")
	return `"` + trimmed + `"`
}
