//go:build windows

package storage

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shlwapi            = windows.NewLazySystemDLL("shlwapi.dll")
	procStrCmpLogicalW = shlwapi.NewProc("StrCmpLogicalW")
)

func compareItemName(a, b string) int {
	if a == b {
		return 0
	}

	aPtr, err := windows.UTF16PtrFromString(a)
	if err != nil {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	}
	bPtr, err := windows.UTF16PtrFromString(b)
	if err != nil {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	}

	ret, _, _ := procStrCmpLogicalW.Call(
		uintptr(unsafe.Pointer(aPtr)),
		uintptr(unsafe.Pointer(bPtr)),
	)

	return int(int32(ret))
}
