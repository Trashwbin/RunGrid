//go:build windows

package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type shortcutResolver struct {
	shell       *ole.IDispatch
	initialized bool
}

func newShortcutResolver() (*shortcutResolver, error) {
	const comSFalse = uintptr(0x00000001)
	initErr := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	if initErr != nil {
		if oleErr, ok := initErr.(*ole.OleError); !ok || oleErr.Code() != comSFalse {
			return nil, initErr
		}
	}

	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		ole.CoUninitialize()
		return nil, err
	}
	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	unknown.Release()
	if err != nil {
		ole.CoUninitialize()
		return nil, err
	}

	return &shortcutResolver{
		shell:       shell,
		initialized: true,
	}, nil
}

func (r *shortcutResolver) Close() {
	if r == nil {
		return
	}
	if r.shell != nil {
		r.shell.Release()
		r.shell = nil
	}
	if r.initialized {
		ole.CoUninitialize()
		r.initialized = false
	}
}

func (r *shortcutResolver) Resolve(path string) (string, string, error) {
	if r == nil || r.shell == nil {
		return "", "", nil
	}

	shortcutVariant, err := oleutil.CallMethod(r.shell, "CreateShortcut", path)
	if err != nil {
		return "", "", err
	}
	shortcut := shortcutVariant.ToIDispatch()
	if shortcut == nil {
		shortcutVariant.Clear()
		return "", "", nil
	}
	defer shortcutVariant.Clear()

	targetVar, err := oleutil.GetProperty(shortcut, "TargetPath")
	if err != nil {
		return "", "", err
	}
	argsVar, err := oleutil.GetProperty(shortcut, "Arguments")
	if err != nil {
		targetVar.Clear()
		return "", "", err
	}

	target := strings.TrimSpace(targetVar.ToString())
	args := strings.TrimSpace(argsVar.ToString())
	targetVar.Clear()
	argsVar.Clear()

	target = strings.Trim(target, "\"'")
	args = strings.Trim(args, "\"'")
	if target != "" {
		target = expandPercentEnv(target)
		target = strings.TrimSpace(target)
		if target != "" && !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(path), target)
		}
	}

	return target, args, nil
}

func shortcutDedupeKey(target, args string) string {
	target = normalizeShortcutPath(target)
	if target == "" {
		return ""
	}
	args = strings.ToLower(strings.TrimSpace(args))
	return target + "\x00" + args
}

func normalizeShortcutPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	return strings.ToLower(filepath.Clean(path))
}

func expandPercentEnv(value string) string {
	if !strings.Contains(value, "%") {
		return value
	}

	var builder strings.Builder
	builder.Grow(len(value))
	for i := 0; i < len(value); {
		if value[i] != '%' {
			builder.WriteByte(value[i])
			i++
			continue
		}

		end := strings.IndexByte(value[i+1:], '%')
		if end == -1 {
			builder.WriteString(value[i:])
			break
		}

		key := value[i+1 : i+1+end]
		if key == "" {
			builder.WriteByte('%')
			i += end + 2
			continue
		}

		if val, ok := os.LookupEnv(key); ok {
			builder.WriteString(val)
		} else {
			builder.WriteByte('%')
			builder.WriteString(key)
			builder.WriteByte('%')
		}
		i += end + 2
	}

	return builder.String()
}
