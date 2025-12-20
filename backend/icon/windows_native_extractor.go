//go:build windows

package icon

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type NativeExtractor struct{}

func NewNativeExtractor() Extractor {
	return NativeExtractor{}
}

func (NativeExtractor) Extract(_ context.Context, source string, dest string) error {
	if err := ValidateSource(source); err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(source))
	if ext == ".png" {
		return CopyFile(source, dest)
	}

	iconSource := source
	iconIndex := 0
	if ext == ".lnk" {
		info, err := resolveShortcutIcon(source)
		if err != nil {
			return err
		}
		if info.source != "" {
			iconSource = info.source
			iconIndex = info.index
		}
	}

	sizes := []int{256, 128, 96, 64, 48, 32}
	for _, size := range sizes {
		handle, err := extractIconHandle(iconSource, iconIndex, size)
		if err != nil || handle == 0 {
			continue
		}
		if err := saveIconPNG(handle, size, dest); err != nil {
			destroyIcon(handle)
			continue
		}
		destroyIcon(handle)
		return nil
	}

	return fmt.Errorf("icon not found")
}

func extractIconHandle(path string, index int, size int) (windows.Handle, error) {
	if strings.TrimSpace(path) == "" {
		return 0, fmt.Errorf("empty source")
	}

	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	var large windows.Handle
	var small windows.Handle
	_, _, _ = procSHDefExtractIconW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(index),
		0,
		uintptr(unsafe.Pointer(&large)),
		uintptr(unsafe.Pointer(&small)),
		uintptr(size),
	)
	if large != 0 {
		if small != 0 {
			destroyIcon(small)
		}
		return large, nil
	}
	if large != 0 {
		destroyIcon(large)
	}
	if small != 0 {
		destroyIcon(small)
	}

	var icon windows.Handle
	var id uint32
	count, _, _ := procPrivateExtractIconsW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(index),
		uintptr(size),
		uintptr(size),
		uintptr(unsafe.Pointer(&icon)),
		uintptr(unsafe.Pointer(&id)),
		1,
		0,
	)
	if count > 0 && icon != 0 {
		return icon, nil
	}

	if icon != 0 {
		destroyIcon(icon)
	}

	return 0, fmt.Errorf("no icon extracted")
}

func saveIconPNG(icon windows.Handle, size int, dest string) error {
	hdc, _, err := procCreateCompatibleDC.Call(0)
	if hdc == 0 {
		return err
	}
	defer procDeleteDC.Call(hdc)

	info := bitmapInfo{
		Header: bitmapInfoHeader{
			Size:        uint32(unsafe.Sizeof(bitmapInfoHeader{})),
			Width:       int32(size),
			Height:      -int32(size),
			Planes:      1,
			BitCount:    32,
			Compression: biRGB,
		},
	}

	var bits unsafe.Pointer
	hbm, _, err := procCreateDIBSection.Call(
		hdc,
		uintptr(unsafe.Pointer(&info)),
		0,
		uintptr(unsafe.Pointer(&bits)),
		0,
		0,
	)
	if hbm == 0 || bits == nil {
		return err
	}
	defer procDeleteObject.Call(hbm)

	prev, _, _ := procSelectObject.Call(hdc, hbm)
	if prev != 0 {
		defer procSelectObject.Call(hdc, prev)
	}

	ret, _, err := procDrawIconEx.Call(
		hdc,
		0,
		0,
		uintptr(icon),
		uintptr(size),
		uintptr(size),
		0,
		0,
		diNormal,
	)
	if ret == 0 {
		return err
	}

	pixels := unsafe.Slice((*byte)(bits), size*size*4)
	data := make([]byte, len(pixels))
	copy(data, pixels)

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	hasAlpha := false
	for i := 0; i+3 < len(data); i += 4 {
		b := data[i]
		g := data[i+1]
		r := data[i+2]
		a := data[i+3]
		if a != 0 {
			hasAlpha = true
		}
		img.Pix[i] = r
		img.Pix[i+1] = g
		img.Pix[i+2] = b
		img.Pix[i+3] = a
	}
	if !hasAlpha {
		if err := applyIconMaskAlpha(icon, size, img); err != nil {
			for i := 3; i < len(img.Pix); i += 4 {
				img.Pix[i] = 0xFF
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return err
	}

	return out.Sync()
}

func destroyIcon(icon windows.Handle) {
	if icon == 0 {
		return
	}
	procDestroyIcon.Call(uintptr(icon))
}

const (
	biRGB        = 0
	diNormal     = 0x0003
	dibRGBColors = 0
)

type bitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type bitmapInfo struct {
	Header bitmapInfoHeader
	Colors [1]uint32
}

type bitmapInfo1bpp struct {
	Header bitmapInfoHeader
	Colors [2]uint32
}

type iconInfo struct {
	IsIcon   int32
	HotspotX uint32
	HotspotY uint32
	Mask     windows.Handle
	Color    windows.Handle
}

func applyIconMaskAlpha(icon windows.Handle, size int, img *image.RGBA) error {
	var info iconInfo
	ret, _, err := procGetIconInfo.Call(
		uintptr(icon),
		uintptr(unsafe.Pointer(&info)),
	)
	if ret == 0 {
		return err
	}
	if info.Mask != 0 {
		defer procDeleteObject.Call(uintptr(info.Mask))
	}
	if info.Color != 0 {
		defer procDeleteObject.Call(uintptr(info.Color))
	}
	if info.Mask == 0 {
		return fmt.Errorf("icon mask unavailable")
	}

	hdc, _, err := procCreateCompatibleDC.Call(0)
	if hdc == 0 {
		return err
	}
	defer procDeleteDC.Call(hdc)

	maskInfo := bitmapInfo1bpp{
		Header: bitmapInfoHeader{
			Size:        uint32(unsafe.Sizeof(bitmapInfoHeader{})),
			Width:       int32(size),
			Height:      -int32(size),
			Planes:      1,
			BitCount:    1,
			Compression: biRGB,
		},
	}
	stride := ((size + 31) / 32) * 4
	buffer := make([]byte, stride*size)
	if len(buffer) == 0 {
		return fmt.Errorf("empty mask buffer")
	}
	ret, _, err = procGetDIBits.Call(
		hdc,
		uintptr(info.Mask),
		0,
		uintptr(size),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&maskInfo)),
		dibRGBColors,
	)
	if ret == 0 {
		return err
	}

	for y := 0; y < size; y++ {
		row := y * stride
		for x := 0; x < size; x++ {
			byteIndex := row + (x / 8)
			bit := 7 - (x % 8)
			maskBit := (buffer[byteIndex] >> bit) & 1
			idx := y*img.Stride + x*4 + 3
			if maskBit == 1 {
				img.Pix[idx] = 0
			} else {
				img.Pix[idx] = 0xFF
			}
		}
	}

	return nil
}

var (
	modShell32               = windows.NewLazySystemDLL("shell32.dll")
	procSHDefExtractIconW    = modShell32.NewProc("SHDefExtractIconW")
	modUser32                = windows.NewLazySystemDLL("user32.dll")
	procPrivateExtractIconsW = modUser32.NewProc("PrivateExtractIconsW")
	procDestroyIcon          = modUser32.NewProc("DestroyIcon")
	procGetIconInfo          = modUser32.NewProc("GetIconInfo")
	procDrawIconEx           = modUser32.NewProc("DrawIconEx")
	modGdi32                 = windows.NewLazySystemDLL("gdi32.dll")
	procCreateCompatibleDC   = modGdi32.NewProc("CreateCompatibleDC")
	procCreateDIBSection     = modGdi32.NewProc("CreateDIBSection")
	procGetDIBits            = modGdi32.NewProc("GetDIBits")
	procSelectObject         = modGdi32.NewProc("SelectObject")
	procDeleteObject         = modGdi32.NewProc("DeleteObject")
	procDeleteDC             = modGdi32.NewProc("DeleteDC")
)
