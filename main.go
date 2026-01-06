package main

import (
	"embed"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app, err := NewApp()
	if err != nil {
		println("Error:", err.Error())
		return
	}

	iconRoot, err := iconRootPath()
	if err != nil {
		println("Error:", err.Error())
		return
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "RunGrid",
		Width:     1024,
		Height:    768,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: iconFileHandler(iconRoot),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "rungrid",
			OnSecondInstanceLaunch: func(_ options.SecondInstanceData) {
				ctx := app.context()
				runtime.WindowShow(ctx)
				runtime.WindowUnminimise(ctx)
				runtime.EventsEmit(ctx, "window:show")
			},
		},
		Windows: &windows.Options{
			DisableWindowIcon: false,
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func iconRootPath() (string, error) {
	dbPath, err := appDataPath("rungrid")
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(dbPath), "icons"), nil
}
