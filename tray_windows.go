//go:build windows

package main

import (
	"context"
	_ "embed"
	"errors"
	"os"
	"runtime"
	"sync/atomic"
	"unsafe"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
)

//go:embed assets/icons/tray.ico
var trayIconData []byte

type trayController struct {
	ctx     context.Context
	app     *App
	hWnd    windows.Handle
	hIcon   windows.Handle
	running atomic.Bool
}

var globalTray trayController

func (t *trayController) setApp(app *App) {
	t.app = app
}

func (t *trayController) start(ctx context.Context) {
	if !t.running.CompareAndSwap(false, true) {
		return
	}
	t.ctx = ctx
	go t.loop()
}

func (t *trayController) stop() {
	if !t.running.Load() {
		return
	}
	if t.hWnd != 0 {
		procPostMessageW.Call(uintptr(t.hWnd), wmClose, 0, 0)
	}
}

func (t *trayController) loop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	inst := getModuleHandle()
	className := windows.StringToUTF16Ptr("RunGridTrayWindow")
	wndProc := windows.NewCallback(trayWndProc)
	wc := wndClassEx{
		cbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		lpfnWndProc:   wndProc,
		hInstance:     inst,
		lpszClassName: className,
	}

	if atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); atom == 0 {
		logError(t.ctx, "tray: RegisterClassExW failed", err)
		t.running.Store(false)
		_ = err
		return
	}

	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(className)),
		0,
		0, 0, 0, 0,
		0,
		0,
		uintptr(inst),
		0,
	)
	if hwnd == 0 {
		logError(t.ctx, "tray: CreateWindowExW failed", err)
		t.running.Store(false)
		_ = err
		return
	}
	t.hWnd = windows.Handle(hwnd)

	icon, iconErr := loadIconFromBytes(trayIconData)
	if iconErr == nil && icon != 0 {
		t.hIcon = icon
		if err := t.addIcon("RunGrid"); err != nil {
			logError(t.ctx, "tray: Shell_NotifyIconW add failed", err)
		}
	} else {
		logError(t.ctx, "tray: create icon failed", iconErr)
		t.running.Store(false)
		return
	}

	t.runMessageLoop()

	t.removeIcon()
	if t.hIcon != 0 {
		_ = destroyIcon(t.hIcon)
	}
	t.hWnd = 0
	t.running.Store(false)
}

func (t *trayController) runMessageLoop() {
	var msg msg
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		switch int32(ret) {
		case -1:
			return
		case 0:
			return
		default:
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}
}

func trayWndProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case trayCallbackMessage:
		switch lParam {
		case wmLButtonUp, wmLButtonDblClk:
			globalTray.showMainWindow()
		case wmRButtonUp:
			globalTray.showContextMenu()
		}
		return 0
	case wmCommand:
		globalTray.handleCommand(uint16(wParam & 0xffff))
		return 0
	case wmDestroy:
		globalTray.removeIcon()
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	return ret
}

func (t *trayController) showMainWindow() {
	if t.ctx == nil {
		return
	}
	wailsRuntime.WindowShow(t.ctx)
	wailsRuntime.WindowUnminimise(t.ctx)
	wailsRuntime.EventsEmit(t.ctx, "window:show")
}

func (t *trayController) handleCommand(id uint16) {
	switch id {
	case cmdShow:
		t.showMainWindow()
	case cmdScan:
		if t.app != nil {
			go t.app.ScanShortcuts(nil)
		}
	case cmdRefreshIcons:
		if t.app != nil {
			go t.app.SyncIcons()
		}
	case cmdQuit:
		if t.ctx != nil {
			wailsRuntime.Quit(t.ctx)
		}
	}
}

func (t *trayController) showContextMenu() {
	menu, _, err := procCreatePopupMenu.Call()
	if menu == 0 {
		_ = err
		return
	}
	defer procDestroyMenu.Call(menu)

	appendMenu(menu, cmdShow, "显示窗口")
	appendMenu(menu, cmdScan, "重新扫描")
	appendMenu(menu, cmdRefreshIcons, "刷新图标缓存")
	appendSeparator(menu)
	appendMenu(menu, cmdQuit, "退出")

	var pos point
	if ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos))); ret == 0 {
		return
	}

	procSetForegroundWindow.Call(uintptr(t.hWnd))
	cmd, _, _ := procTrackPopupMenu.Call(
		menu,
		tpmRightButton|tpmReturnCmd,
		uintptr(pos.X),
		uintptr(pos.Y),
		0,
		uintptr(t.hWnd),
		0,
	)
	if cmd != 0 {
		t.handleCommand(uint16(cmd))
	}
}

func appendMenu(menu uintptr, id uint16, label string) {
	text := windows.StringToUTF16Ptr(label)
	procAppendMenuW.Call(menu, mfString, uintptr(id), uintptr(unsafe.Pointer(text)))
}

func appendSeparator(menu uintptr) {
	procAppendMenuW.Call(menu, mfSeparator, 0, 0)
}

func (t *trayController) addIcon(tooltip string) error {
	var nid notifyIconData
	nid.cbSize = uint32(unsafe.Sizeof(nid))
	nid.hWnd = t.hWnd
	nid.uID = 1
	nid.uFlags = nifMessage | nifIcon | nifTip
	nid.uCallbackMessage = trayCallbackMessage
	nid.hIcon = t.hIcon
	setUTF16(nid.szTip[:], tooltip)
	r, _, err := procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&nid)))
	if r == 0 {
		return err
	}
	procShellNotifyIconW.Call(nimSetVersion, uintptr(unsafe.Pointer(&nid)))
	return nil
}

func (t *trayController) removeIcon() {
	if t.hWnd == 0 {
		return
	}
	var nid notifyIconData
	nid.cbSize = uint32(unsafe.Sizeof(nid))
	nid.hWnd = t.hWnd
	nid.uID = 1
	procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&nid)))
}

func createIconFromResource(data []byte) (windows.Handle, error) {
	return loadIconFromBytes(data)
}

func setUTF16(dst []uint16, value string) {
	encoded, _ := windows.UTF16FromString(value)
	max := len(dst)
	for i := 0; i < len(encoded) && i < max; i++ {
		dst[i] = encoded[i]
	}
	if max > 0 {
		dst[max-1] = 0
	}
}

type wndClassEx struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     windows.Handle
	hIcon         windows.Handle
	hCursor       windows.Handle
	hbrBackground windows.Handle
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       windows.Handle
}

type msg struct {
	hWnd    windows.Handle
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
}

type point struct {
	X int32
	Y int32
}

type notifyIconData struct {
	cbSize           uint32
	hWnd             windows.Handle
	uID              uint32
	uFlags           uint32
	uCallbackMessage uint32
	hIcon            windows.Handle
	szTip            [128]uint16
	dwState          uint32
	dwStateMask      uint32
	szInfo           [256]uint16
	uTimeout         uint32
	szInfoTitle      [64]uint16
	dwInfoFlags      uint32
	guidItem         windows.GUID
	hBalloonIcon     windows.Handle
}

const (
	nimAdd        = 0x00000000
	nimDelete     = 0x00000002
	nimSetVersion = 0x00000004

	nifMessage = 0x00000001
	nifIcon    = 0x00000002
	nifTip     = 0x00000004

	wmCommand           = 0x0111
	wmClose             = 0x0010
	wmDestroy           = 0x0002
	wmLButtonUp         = 0x0202
	wmLButtonDblClk     = 0x0203
	wmRButtonUp         = 0x0205
	wmUser              = 0x0400
	trayCallbackMessage = wmUser + 1

	tpmRightButton = 0x0002
	tpmReturnCmd   = 0x0100

	mfString    = 0x0000
	mfSeparator = 0x0800

	cmdShow         = 1001
	cmdScan         = 1002
	cmdRefreshIcons = 1003
	cmdQuit         = 1004

	lrDefaultColor = 0x00000000
	imageIcon      = 1
	lrLoadFromFile = 0x00000010
)

var (
	modUser32                    = windows.NewLazySystemDLL("user32.dll")
	modShell32                   = windows.NewLazySystemDLL("shell32.dll")
	modKernel32                  = windows.NewLazySystemDLL("kernel32.dll")
	procRegisterClassExW         = modUser32.NewProc("RegisterClassExW")
	procCreateWindowExW          = modUser32.NewProc("CreateWindowExW")
	procDefWindowProcW           = modUser32.NewProc("DefWindowProcW")
	procPostQuitMessage          = modUser32.NewProc("PostQuitMessage")
	procGetMessageW              = modUser32.NewProc("GetMessageW")
	procTranslateMessage         = modUser32.NewProc("TranslateMessage")
	procDispatchMessageW         = modUser32.NewProc("DispatchMessageW")
	procPostMessageW             = modUser32.NewProc("PostMessageW")
	procCreatePopupMenu          = modUser32.NewProc("CreatePopupMenu")
	procAppendMenuW              = modUser32.NewProc("AppendMenuW")
	procTrackPopupMenu           = modUser32.NewProc("TrackPopupMenu")
	procGetCursorPos             = modUser32.NewProc("GetCursorPos")
	procSetForegroundWindow      = modUser32.NewProc("SetForegroundWindow")
	procDestroyMenu              = modUser32.NewProc("DestroyMenu")
	procShellNotifyIconW         = modShell32.NewProc("Shell_NotifyIconW")
	procCreateIconFromResourceEx = modUser32.NewProc("CreateIconFromResourceEx")
	procGetModuleHandleW         = modKernel32.NewProc("GetModuleHandleW")
	procDestroyIcon              = modUser32.NewProc("DestroyIcon")
	procLoadImageW               = modUser32.NewProc("LoadImageW")
)

func getModuleHandle() windows.Handle {
	h, _, _ := procGetModuleHandleW.Call(0)
	return windows.Handle(h)
}

func destroyIcon(icon windows.Handle) error {
	r, _, err := procDestroyIcon.Call(uintptr(icon))
	if r == 0 {
		return err
	}
	return nil
}

func loadIconFromBytes(data []byte) (windows.Handle, error) {
	if len(data) == 0 {
		return 0, errors.New("empty icon data")
	}

	tmp, err := os.CreateTemp("", "rungrid_tray_*.ico")
	if err != nil {
		return 0, err
	}
	path := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(path)
		return 0, err
	}
	_ = tmp.Close()
	defer os.Remove(path)

	ptr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	h, _, loadErr := procLoadImageW.Call(
		0,
		uintptr(unsafe.Pointer(ptr)),
		imageIcon,
		0,
		0,
		lrLoadFromFile,
	)
	if h == 0 {
		return 0, loadErr
	}
	return windows.Handle(h), nil
}

func logError(ctx context.Context, msg string, err error) {
	if err != nil {
		msg = msg + ": " + err.Error()
	}
	if ctx != nil {
		wailsRuntime.LogError(ctx, msg)
		return
	}
	println(msg)
}
