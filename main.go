package main

import (
	"embed"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

var ScreenW, ScreenH int

func getScreenSize() (int, int) {
	user32 := syscall.NewLazyDLL("user32.dll")
	getSystemMetrics := user32.NewProc("GetSystemMetrics")
	const smCxScreen = 0
	const smCyScreen = 1
	w, _, _ := getSystemMetrics.Call(uintptr(smCxScreen))
	h, _, _ := getSystemMetrics.Call(uintptr(smCyScreen))
	ScreenW, ScreenH = int(w), int(h)
	return ScreenW, ScreenH
}

func main() {
	app := NewApp()

	screenW, screenH := getScreenSize()

	err := wails.Run(&options.App{
		Title:       "Colorfui",
		Width:       screenW,
		Height:      screenH,
		Frameless:   true,
		AlwaysOnTop: true,
		BackgroundColour: &options.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		},
		DisableResize:     true,
		StartHidden:       true,
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
