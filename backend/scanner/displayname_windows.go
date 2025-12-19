//go:build windows

package scanner

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

const shgfiDisplayName = 0x000000200

type shfileinfo struct {
	hIcon         windows.Handle
	iIcon         int32
	dwAttributes  uint32
	szDisplayName [windows.MAX_PATH]uint16
	szTypeName    [80]uint16
}

var (
	shell32Display     = windows.NewLazySystemDLL("shell32.dll")
	procSHGetFileInfoW = shell32Display.NewProc("SHGetFileInfoW")
)

func lookupDisplayName(path string) (string, bool) {
	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", false
	}

	var info shfileinfo
	ret, _, _ := procSHGetFileInfoW.Call(
		uintptr(unsafe.Pointer(ptr)),
		0,
		uintptr(unsafe.Pointer(&info)),
		uintptr(unsafe.Sizeof(info)),
		shgfiDisplayName,
	)
	if ret == 0 {
		return "", false
	}

	name := strings.TrimSpace(windows.UTF16ToString(info.szDisplayName[:]))
	if name == "" {
		return "", false
	}
	return name, true
}
