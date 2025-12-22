//go:build windows

package main

import (
	"errors"
	"syscall"
	"unsafe"

	"rungrid/backend/domain"
)

type winPoint struct {
	X int32
	Y int32
}

type winRect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type winMonitorInfo struct {
	CbSize    uint32
	RcMonitor winRect
	RcWork    winRect
	DwFlags   uint32
}

const monitorDefaultToNearest = 0x00000002

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos     = user32.NewProc("GetCursorPos")
	procMonitorFromPoint = user32.NewProc("MonitorFromPoint")
	procGetMonitorInfoW  = user32.NewProc("GetMonitorInfoW")
)

func getCursorPos(point *winPoint) error {
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(point)))
	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return errors.New("GetCursorPos failed")
	}
	return nil
}

func monitorFromPoint(point winPoint, flags uint32) uintptr {
	returnUint := uintptr(uint32(point.X)) | (uintptr(uint32(point.Y)) << 32)
	monitor, _, _ := procMonitorFromPoint.Call(returnUint, uintptr(flags))
	return monitor
}

func getMonitorInfo(monitor uintptr, info *winMonitorInfo) error {
	ret, _, err := procGetMonitorInfoW.Call(monitor, uintptr(unsafe.Pointer(info)))
	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return errors.New("GetMonitorInfo failed")
	}
	return nil
}

func (a *App) GetCursorAnchorPosition(width, height int) (domain.Point, error) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

	var pt winPoint
	if err := getCursorPos(&pt); err != nil {
		return domain.Point{}, err
	}

	x := int(pt.X) - width/2
	y := int(pt.Y) - height/2

	monitor := monitorFromPoint(pt, monitorDefaultToNearest)
	if monitor != 0 {
		var info winMonitorInfo
		info.CbSize = uint32(unsafe.Sizeof(info))
		if err := getMonitorInfo(monitor, &info); err == nil {
			left := int(info.RcWork.Left)
			top := int(info.RcWork.Top)
			right := int(info.RcWork.Right)
			bottom := int(info.RcWork.Bottom)

			if width > 0 {
				if right-left < width {
					x = left
				} else {
					if x < left {
						x = left
					}
					if x+width > right {
						x = right - width
					}
				}
			}

			if height > 0 {
				if bottom-top < height {
					y = top
				} else {
					if y < top {
						y = top
					}
					if y+height > bottom {
						y = bottom - height
					}
				}
			}
		}
	}

	return domain.Point{X: x, Y: y}, nil
}
