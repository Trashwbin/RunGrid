//go:build windows

package hotkey

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"rungrid/backend/domain"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
)

type Manager struct {
	mu     sync.Mutex
	ctx    context.Context
	worker *hotkeyWorker
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Start(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ctx = ctx
	if m.worker == nil {
		m.worker = newHotkeyWorker(ctx)
	}
}

func (m *Manager) Stop() {
	m.mu.Lock()
	worker := m.worker
	m.worker = nil
	m.mu.Unlock()
	if worker != nil {
		worker.stop()
	}
}

func (m *Manager) Apply(bindings []domain.HotkeyBinding) ([]domain.HotkeyIssue, error) {
	m.mu.Lock()
	if m.worker == nil {
		if m.ctx == nil {
			m.mu.Unlock()
			return nil, ErrUnsupported
		}
		m.worker = newHotkeyWorker(m.ctx)
	}
	worker := m.worker
	m.mu.Unlock()
	return worker.apply(bindings)
}

type hotkeyWorker struct {
	ctx        context.Context
	commands   chan applyCommand
	registered map[int]registeredHotkey
	ready      chan struct{}
	done       chan struct{}
	threadID   uint32
}

type registeredHotkey struct {
	id     int
	action string
	keys   string
	mod    uint
	vk     uint
}

type applyCommand struct {
	bindings []domain.HotkeyBinding
	reply    chan applyResult
}

type applyResult struct {
	issues []domain.HotkeyIssue
}

func newHotkeyWorker(ctx context.Context) *hotkeyWorker {
	worker := &hotkeyWorker{
		ctx:        ctx,
		commands:   make(chan applyCommand, 1),
		registered: make(map[int]registeredHotkey),
		ready:      make(chan struct{}),
		done:       make(chan struct{}),
	}
	go worker.loop()
	return worker
}

func (w *hotkeyWorker) apply(bindings []domain.HotkeyBinding) ([]domain.HotkeyIssue, error) {
	<-w.ready
	reply := make(chan applyResult, 1)
	w.commands <- applyCommand{bindings: bindings, reply: reply}
	w.wake()
	result := <-reply
	return result.issues, nil
}

func (w *hotkeyWorker) stop() {
	<-w.ready
	postThreadMessage(w.threadID, wmQuit, 0, 0)
	<-w.done
}

func (w *hotkeyWorker) loop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var msg msg
	peekMessage(&msg, 0, 0, 0, pmNoremove)
	w.threadID = windows.GetCurrentThreadId()
	close(w.ready)

	for {
		result, _ := getMessage(&msg, 0, 0, 0)
		if result == 0 {
			break
		}
		if result == -1 {
			continue
		}

		switch msg.message {
		case wmHotkey:
			w.handleHotkey(int(msg.wParam))
		case wmHotkeyCommand:
			w.drainCommands()
		default:
			translateMessage(&msg)
			dispatchMessage(&msg)
		}
	}

	w.unregisterAll()
	close(w.done)
}

func (w *hotkeyWorker) handleHotkey(id int) {
	if w.ctx == nil {
		return
	}
	hk, ok := w.registered[id]
	if !ok {
		return
	}
	wailsRuntime.EventsEmit(w.ctx, "hotkey:trigger", hk.action)
}

func (w *hotkeyWorker) drainCommands() {
	for {
		select {
		case cmd := <-w.commands:
			issues := w.applyBindings(cmd.bindings)
			cmd.reply <- applyResult{issues: issues}
		default:
			return
		}
	}
}

func (w *hotkeyWorker) applyBindings(bindings []domain.HotkeyBinding) []domain.HotkeyIssue {
	w.unregisterAll()
	w.registered = make(map[int]registeredHotkey)

	var issues []domain.HotkeyIssue
	seen := make(map[string]string)
	nextID := 1

	for _, binding := range bindings {
		if strings.TrimSpace(binding.ID) == "" {
			issues = append(issues, domain.HotkeyIssue{
				ID:     binding.ID,
				Keys:   binding.Keys,
				Reason: "无效的动作编号",
			})
			continue
		}

		keys := strings.TrimSpace(binding.Keys)
		if keys == "" {
			continue
		}

		mod, vk, err := parseHotkey(keys)
		if err != nil {
			issues = append(issues, domain.HotkeyIssue{
				ID:     binding.ID,
				Keys:   binding.Keys,
				Reason: err.Error(),
			})
			continue
		}

		signature := fmt.Sprintf("%d:%d", mod, vk)
		if _, exists := seen[signature]; exists {
			issues = append(issues, domain.HotkeyIssue{
				ID:     binding.ID,
				Keys:   binding.Keys,
				Reason: "快捷键重复",
			})
			continue
		}
		seen[signature] = binding.ID

		if err := registerHotKey(nextID, mod, vk); err != nil {
			issues = append(issues, domain.HotkeyIssue{
				ID:     binding.ID,
				Keys:   binding.Keys,
				Reason: "快捷键冲突或已被占用",
			})
			continue
		}

		w.registered[nextID] = registeredHotkey{
			id:     nextID,
			action: binding.ID,
			keys:   keys,
			mod:    mod,
			vk:     vk,
		}
		nextID++
	}

	return issues
}

func (w *hotkeyWorker) unregisterAll() {
	for id := range w.registered {
		_ = unregisterHotKey(id)
	}
	w.registered = make(map[int]registeredHotkey)
}

func (w *hotkeyWorker) wake() {
	if w.threadID == 0 {
		return
	}
	postThreadMessage(w.threadID, wmHotkeyCommand, 0, 0)
}

func parseHotkey(value string) (uint, uint, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, 0, fmt.Errorf("快捷键为空")
	}

	var mod uint
	var keyToken string

	parts := strings.Split(value, "+")
	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token == "" {
			continue
		}
		upper := strings.ToUpper(token)
		switch upper {
		case "CTRL", "CONTROL":
			mod |= modControl
		case "ALT":
			mod |= modAlt
		case "SHIFT":
			mod |= modShift
		case "WIN", "META", "CMD", "COMMAND":
			mod |= modWin
		default:
			if keyToken != "" {
				return 0, 0, fmt.Errorf("快捷键格式不支持")
			}
			keyToken = token
		}
	}

	if keyToken == "" {
		return 0, 0, fmt.Errorf("快捷键缺少主键")
	}

	vk, ok := keyToVK(keyToken)
	if !ok {
		return 0, 0, fmt.Errorf("不支持的按键")
	}
	if mod == 0 {
		if vk < vkF1 || vk > vkF1+23 {
			return 0, 0, fmt.Errorf("快捷键需要包含修饰键")
		}
	}
	return mod, vk, nil
}

func keyToVK(token string) (uint, bool) {
	upper := strings.ToUpper(strings.TrimSpace(token))
	if upper == "" {
		return 0, false
	}

	if len(upper) == 1 {
		ch := upper[0]
		if ch >= 'A' && ch <= 'Z' {
			return uint(ch), true
		}
		if ch >= '0' && ch <= '9' {
			return uint(ch), true
		}
	}

	if strings.HasPrefix(upper, "F") && len(upper) <= 3 {
		number := strings.TrimPrefix(upper, "F")
		if n := parseNumber(number); n >= 1 && n <= 24 {
			return uint(vkF1 + n - 1), true
		}
	}

	if vk, ok := keyMap[upper]; ok {
		return vk, true
	}
	return 0, false
}

func parseNumber(value string) int {
	total := 0
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return -1
		}
		total = total*10 + int(ch-'0')
	}
	return total
}

func registerHotKey(id int, mod uint, vk uint) error {
	ret, _, err := procRegisterHotKey.Call(
		0,
		uintptr(id),
		uintptr(mod),
		uintptr(vk),
	)
	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return fmt.Errorf("register hotkey failed")
	}
	return nil
}

func unregisterHotKey(id int) error {
	ret, _, err := procUnregisterHotKey.Call(0, uintptr(id))
	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return fmt.Errorf("unregister hotkey failed")
	}
	return nil
}

func postThreadMessage(threadID uint32, msg uint32, wParam uintptr, lParam uintptr) {
	if threadID == 0 {
		return
	}
	procPostThreadMessageW.Call(uintptr(threadID), uintptr(msg), wParam, lParam)
}

func getMessage(msg *msg, hwnd uintptr, minMsg uint32, maxMsg uint32) (int32, error) {
	ret, _, err := procGetMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		hwnd,
		uintptr(minMsg),
		uintptr(maxMsg),
	)
	return int32(ret), err
}

func translateMessage(msg *msg) {
	procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}

func dispatchMessage(msg *msg) {
	procDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
}

func peekMessage(msg *msg, hwnd uintptr, minMsg uint32, maxMsg uint32, removeMsg uint32) {
	procPeekMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		hwnd,
		uintptr(minMsg),
		uintptr(maxMsg),
		uintptr(removeMsg),
	)
}

const (
	modAlt     = 0x0001
	modControl = 0x0002
	modShift   = 0x0004
	modWin     = 0x0008

	wmHotkey        = 0x0312
	wmQuit          = 0x0012
	wmHotkeyCommand = 0x8001

	pmNoremove = 0x0000

	vkSpace = 0x20
	vkTab   = 0x09
	vkEnter = 0x0D
	vkEsc   = 0x1B
	vkLeft  = 0x25
	vkUp    = 0x26
	vkRight = 0x27
	vkDown  = 0x28

	vkOemComma     = 0xBC
	vkOemPeriod    = 0xBE
	vkOemMinus     = 0xBD
	vkOemPlus      = 0xBB
	vkOemSemicolon = 0xBA
	vkOemSlash     = 0xBF
	vkOemBackslash = 0xDC
	vkOemLBracket  = 0xDB
	vkOemRBracket  = 0xDD
	vkOemQuote     = 0xDE
	vkOemTilde     = 0xC0

	vkF1 = 0x70
)

var keyMap = map[string]uint{
	"SPACE":        vkSpace,
	"TAB":          vkTab,
	"ENTER":        vkEnter,
	"RETURN":       vkEnter,
	"ESC":          vkEsc,
	"ESCAPE":       vkEsc,
	"LEFT":         vkLeft,
	"UP":           vkUp,
	"RIGHT":        vkRight,
	"DOWN":         vkDown,
	",":            vkOemComma,
	"COMMA":        vkOemComma,
	".":            vkOemPeriod,
	"PERIOD":       vkOemPeriod,
	"DOT":          vkOemPeriod,
	"-":            vkOemMinus,
	"MINUS":        vkOemMinus,
	"=":            vkOemPlus,
	"PLUS":         vkOemPlus,
	"EQUALS":       vkOemPlus,
	";":            vkOemSemicolon,
	"SEMICOLON":    vkOemSemicolon,
	"/":            vkOemSlash,
	"SLASH":        vkOemSlash,
	"\\":           vkOemBackslash,
	"BACKSLASH":    vkOemBackslash,
	"[":            vkOemLBracket,
	"BRACKETLEFT":  vkOemLBracket,
	"LBRACKET":     vkOemLBracket,
	"]":            vkOemRBracket,
	"BRACKETRIGHT": vkOemRBracket,
	"RBRACKET":     vkOemRBracket,
	"'":            vkOemQuote,
	"QUOTE":        vkOemQuote,
	"`":            vkOemTilde,
	"BACKQUOTE":    vkOemTilde,
	"TILDE":        vkOemTilde,
}

type point struct {
	X int32
	Y int32
}

type msg struct {
	hwnd     windows.Handle
	message  uint32
	wParam   uintptr
	lParam   uintptr
	time     uint32
	pt       point
	lPrivate uint32
}

var (
	modUser32              = windows.NewLazySystemDLL("user32.dll")
	procRegisterHotKey     = modUser32.NewProc("RegisterHotKey")
	procUnregisterHotKey   = modUser32.NewProc("UnregisterHotKey")
	procPostThreadMessageW = modUser32.NewProc("PostThreadMessageW")
	procGetMessageW        = modUser32.NewProc("GetMessageW")
	procTranslateMessage   = modUser32.NewProc("TranslateMessage")
	procDispatchMessageW   = modUser32.NewProc("DispatchMessageW")
	procPeekMessageW       = modUser32.NewProc("PeekMessageW")
)
