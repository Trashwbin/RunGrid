package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"rungrid/backend/domain"
	"rungrid/backend/hotkey"
	"rungrid/backend/icon"
	"rungrid/backend/launcher"
	"rungrid/backend/scanner"
	"rungrid/backend/service"
	"rungrid/backend/storage"
	"rungrid/backend/storage/sqlite"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	items    *service.ItemService
	groups   *service.GroupService
	icons    *service.IconService
	scanner  *service.ScannerService
	launcher *service.LauncherService
	hotkeys  *hotkey.Manager
	closeFn  func() error
}

// NewApp creates a new App application struct
func NewApp() (*App, error) {
	dbPath, err := appDataPath("rungrid")
	if err != nil {
		return nil, err
	}

	db, err := sqlite.Open(dbPath)
	if err != nil {
		return nil, err
	}

	if err := sqlite.EnsureSchema(context.Background(), db); err != nil {
		_ = db.Close()
		return nil, err
	}

	itemRepo := sqlite.NewItemRepository(db)
	groupRepo := sqlite.NewGroupRepository(db)

	itemService := service.NewItemService(itemRepo)
	groupService := service.NewGroupService(groupRepo)

	iconRoot := filepath.Join(filepath.Dir(dbPath), "icons")
	iconCache := icon.NewCache(iconRoot, icon.NewHybridExtractor())
	iconService := service.NewIconService(iconCache, itemService)
	hotkeyManager := hotkey.NewManager()
	app := &App{
		items:    itemService,
		groups:   groupService,
		icons:    iconService,
		scanner:  service.NewScannerService(scanner.NewDefaultScanner(), itemService, iconService),
		launcher: service.NewLauncherService(launcher.NewDefaultLauncher(), itemService),
		hotkeys:  hotkeyManager,
		closeFn:  db.Close,
	}

	globalTray.setApp(app)

	if err := service.EnsureDefaultGroups(context.Background(), app.groups); err != nil {
		_ = db.Close()
		return nil, err
	}

	return app, nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	globalTray.start(ctx)
	if a.hotkeys != nil {
		a.hotkeys.Start(ctx)
	}
}

// shutdown is called when the app is terminating.
func (a *App) shutdown(ctx context.Context) {
	globalTray.stop()
	if a.hotkeys != nil {
		a.hotkeys.Stop()
	}
	if a.closeFn != nil {
		_ = a.closeFn()
	}
}

func (a *App) ListItems(groupID string, query string) ([]domain.Item, error) {
	return a.items.List(a.context(), storage.ItemFilter{GroupID: groupID, Query: query})
}

func (a *App) CreateItem(input domain.ItemInput) (domain.Item, error) {
	return a.items.Create(a.context(), input)
}

func (a *App) UpdateItem(input domain.ItemUpdate) (domain.Item, error) {
	if strings.TrimSpace(input.Path) != "" && input.Type == "" {
		input.Type = scanner.ClassifyPath(input.Path)
	}
	return a.items.Update(a.context(), input)
}

func (a *App) DeleteItem(id string) error {
	return a.items.Delete(a.context(), id)
}

func (a *App) ClearItems() (int, error) {
	return a.items.Clear(a.context())
}

func (a *App) RecordLaunch(id string) (domain.Item, error) {
	return a.items.RecordLaunch(a.context(), id)
}

func (a *App) ScanShortcuts(roots []string) (domain.ScanResult, error) {
	if a.scanner == nil {
		return domain.ScanResult{}, scanner.ErrUnsupported
	}
	return a.scanner.ScanWithRoots(a.context(), roots)
}

func (a *App) ListScanRoots() ([]string, error) {
	return scanner.NormalizeRoots(scanner.DefaultRoots()), nil
}

func (a *App) PickScanRoot() (string, error) {
	return runtime.OpenDirectoryDialog(a.context(), runtime.OpenDialogOptions{
		Title: "选择扫描目录",
	})
}

func (a *App) PickIconSource() (string, error) {
	return runtime.OpenFileDialog(a.context(), runtime.OpenDialogOptions{
		Title: "选择图标来源",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "图标/程序文件 (*.png;*.ico;*.exe;*.lnk)",
				Pattern:     "*.png;*.ico;*.exe;*.lnk",
			},
		},
	})
}

func (a *App) PreviewIconFromSource(source string) (string, error) {
	if a.icons == nil {
		return "", icon.ErrUnsupported
	}
	return a.icons.PreviewFromSource(a.context(), source)
}

func (a *App) PickTargetPath() (string, error) {
	return runtime.OpenFileDialog(a.context(), runtime.OpenDialogOptions{
		Title: "选择目标文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "应用/快捷方式 (*.exe;*.lnk)",
				Pattern:     "*.exe;*.lnk",
			},
			{
				DisplayName: "网页/链接 (*.url;*.htm;*.html;*.mht;*.mhtml)",
				Pattern:     "*.url;*.htm;*.html;*.mht;*.mhtml",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	})
}

func (a *App) PickTargetFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.context(), runtime.OpenDialogOptions{
		Title: "选择目标文件夹",
	})
}

func (a *App) PickRuleFile() (string, error) {
	return runtime.OpenFileDialog(a.context(), runtime.OpenDialogOptions{
		Title: "选择分组规则",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "规则文件 (*.json)",
				Pattern:     "*.json",
			},
		},
	})
}

func (a *App) ImportGroupRules(path string) (domain.RuleImportResult, error) {
	if strings.TrimSpace(path) == "" {
		return domain.RuleImportResult{}, storage.ErrInvalidInput
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.RuleImportResult{}, err
	}
	return service.ImportGroupRules(a.context(), data, a.groups, a.items)
}

func (a *App) UpdateItemIconFromSource(id string, source string) (domain.Item, error) {
	if a.icons == nil {
		return domain.Item{}, icon.ErrUnsupported
	}
	return a.icons.UpdateFromSource(a.context(), id, source)
}

func (a *App) SyncIcons() (int, error) {
	if a.icons == nil {
		return 0, icon.ErrUnsupported
	}
	return a.icons.RefreshAll(a.context())
}

func (a *App) RefreshItemIcon(id string) (domain.Item, error) {
	if a.icons == nil {
		return domain.Item{}, icon.ErrUnsupported
	}
	return a.icons.RefreshItem(a.context(), id)
}

func (a *App) LaunchItem(id string) (domain.Item, error) {
	if a.launcher == nil {
		return domain.Item{}, launcher.ErrUnsupported
	}
	return a.launcher.LaunchItem(a.context(), id)
}

func (a *App) OpenItemLocation(id string) error {
	if a.launcher == nil {
		return launcher.ErrUnsupported
	}
	return a.launcher.OpenItemLocation(a.context(), id)
}

func (a *App) SetFavorite(id string, favorite bool) (domain.Item, error) {
	return a.items.SetFavorite(a.context(), id, favorite)
}

func (a *App) ApplyHotkeys(bindings []domain.HotkeyBinding) (domain.HotkeyApplyResult, error) {
	if a.hotkeys == nil {
		return domain.HotkeyApplyResult{}, hotkey.ErrUnsupported
	}

	issues, err := a.hotkeys.Apply(bindings)
	if err != nil {
		return domain.HotkeyApplyResult{}, err
	}

	return domain.HotkeyApplyResult{Issues: issues}, nil
}

func (a *App) ListGroups() ([]domain.Group, error) {
	return a.groups.List(a.context())
}

func (a *App) CreateGroup(input domain.GroupInput) (domain.Group, error) {
	return a.groups.Create(a.context(), input)
}

func (a *App) UpdateGroup(input domain.Group) (domain.Group, error) {
	return a.groups.Update(a.context(), input)
}

func (a *App) DeleteGroup(id string) error {
	return a.groups.Delete(a.context(), id)
}

func (a *App) context() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

func appDataPath(appName string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		configDir = "."
	}

	root := filepath.Join(configDir, appName)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(root, "rungrid.db"), nil
}
