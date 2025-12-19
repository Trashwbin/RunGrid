//go:build windows

package scanner

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func latestFileTimestamp(path string) (time.Time, bool) {
	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return time.Time{}, false
	}

	var data windows.Win32FileAttributeData
	err = windows.GetFileAttributesEx(ptr, windows.GetFileExInfoStandard, (*byte)(unsafe.Pointer(&data)))
	if err != nil {
		return time.Time{}, false
	}

	created := time.Unix(0, data.CreationTime.Nanoseconds())
	updated := time.Unix(0, data.LastWriteTime.Nanoseconds())
	if updated.After(created) {
		return updated, true
	}
	return created, true
}
