package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Clipboard UI Demo",
		Width:     240,
		Height:    460,
		MinWidth:  240,
		MinHeight: 300,
		Frameless: true,
		AlwaysOnTop: true,
		BackgroundColour: &options.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		},
		DisableResize:     true,
		StartHidden:       false,
		HideWindowOnClose: false,
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
		Assets: assets,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
