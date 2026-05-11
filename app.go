package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ListItem struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

type App struct {
	ctx    context.Context
	nextID atomic.Int64
}

func NewApp() *App {
	a := &App{}
	a.nextID.Store(100)
	return a
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("=== [DIAG] Startup called ===")
	go a.applyWindowStyles()
}

const (
	gwlStyle       = ^uintptr(15) // -16
	gwlExStyle     = ^uintptr(19) // -20
	wsSysMenu      = 0x00080000
	wsExNoActivate = 0x08000000
	wsExToolWindow = 0x00000080
	swpNoActivate  = 0x0010
	swpNoMove      = 0x0002
	swpNoSize      = 0x0001
	swpFrameChanged = 0x0020
	swpShowWindow  = 0x0040
	hwndTopmost    = ^uintptr(0) // -1
)

var (
	hCurrentWnd uintptr
)

func findWindowByProcess() uintptr {
	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	enumWindows := user32.NewProc("EnumWindows")
	getWindowThreadProcessId := user32.NewProc("GetWindowThreadProcessId")
	getCurrentProcessId := kernel32.NewProc("GetCurrentProcessId")
	getClassNameW := user32.NewProc("GetClassNameW")
	getWindowTextW := user32.NewProc("GetWindowTextW")

	pid, _, _ := getCurrentProcessId.Call()
	fmt.Printf("=== [DIAG] Current PID: %d ===\n", pid)

	var found uintptr
	var foundClass string
	cb := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		var wndPid uint32
		getWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&wndPid)))
		if uintptr(wndPid) != pid {
			return 1
		}
		buf := make([]uint16, 256)
		n, _, _ := getClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), 256)
		if n == 0 {
			return 1
		}
		className := syscall.UTF16ToString(buf)
		// Get window title too
		titleBuf := make([]uint16, 256)
		getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&titleBuf[0])), 256)
		title := syscall.UTF16ToString(titleBuf)
		fmt.Printf("=== [DIAG] EnumWindow found: hwnd=0x%x class=%q title=%q ===\n", hwnd, className, title)
		found = hwnd
		foundClass = className
		return 0
	})

	enumWindows.Call(cb, 0)
	fmt.Printf("=== [DIAG] Selected hwnd=0x%x class=%q ===\n", found, foundClass)
	return found
}

func (a *App) applyWindowStyles() {
	user32 := syscall.NewLazyDLL("user32.dll")
	getWindowLongW := user32.NewProc("GetWindowLongW")
	setWindowLongW := user32.NewProc("SetWindowLongW")
	setWindowPos := user32.NewProc("SetWindowPos")
	showWindow := user32.NewProc("ShowWindow")

	dwmapi := syscall.NewLazyDLL("dwmapi.dll")
	dwmSetWindowAttribute := dwmapi.NewProc("DwmSetWindowAttribute")

	fmt.Println("=== [DIAG] applyWindowStyles() started ===")

	// Retry until window handle is found
	var hwnd uintptr
	for i := 0; i < 30; i++ {
		fmt.Printf("=== [DIAG] Attempt %d to find window... ===\n", i+1)
		hwnd = findWindowByProcess()
		if hwnd != 0 {
			fmt.Printf("=== [DIAG] Window found on attempt %d: 0x%x ===\n", i+1, hwnd)
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if hwnd == 0 {
		fmt.Println("=== [DIAG] FAILED to find window handle! ===")
		return
	}

	hCurrentWnd = hwnd

	// Read current regular style
	ws, _, _ := getWindowLongW.Call(hwnd, gwlStyle)
	fmt.Printf("=== [DIAG] Original GWL_STYLE: 0x%x (has WS_SYSMENU=%v) ===\n",
		ws, (ws&uintptr(wsSysMenu)) != 0)

	// Remove WS_SYSMENU
	setWindowLongW.Call(hwnd, gwlStyle, ws & ^uintptr(wsSysMenu))

	// Verify
	ws2, _, _ := getWindowLongW.Call(hwnd, gwlStyle)
	fmt.Printf("=== [DIAG] New GWL_STYLE:      0x%x (has WS_SYSMENU=%v) ===\n",
		ws2, (ws2&uintptr(wsSysMenu)) != 0)

	// Read current extended style
	ex, _, _ := getWindowLongW.Call(hwnd, gwlExStyle)
	fmt.Printf("=== [DIAG] Original GWL_EXSTYLE: 0x%x (NOACTIVATE=%v TOOLWINDOW=%v) ===\n",
		ex, (ex&wsExNoActivate) != 0, (ex&wsExToolWindow) != 0)

	// Add WS_EX_NOACTIVATE | WS_EX_TOOLWINDOW
	newEx := ex | wsExNoActivate | wsExToolWindow
	setWindowLongW.Call(hwnd, gwlExStyle, newEx)

	// Verify
	ex2, _, _ := getWindowLongW.Call(hwnd, gwlExStyle)
	fmt.Printf("=== [DIAG] New GWL_EXSTYLE:       0x%x (NOACTIVATE=%v TOOLWINDOW=%v) ===\n",
		ex2, (ex2&wsExNoActivate) != 0, (ex2&wsExToolWindow) != 0)

	// Disable DWM transitions to prevent focus-based rendering changes
	const dwmwaTransitionsForceDisabled = 3
	dwmSetWindowAttribute.Call(hwnd, uintptr(dwmwaTransitionsForceDisabled),
		uintptr(unsafe.Pointer(&[]bool{true}[0])), 4)
	fmt.Println("=== [DIAG] DWM transitions disabled ===")

	// Disable DWM non-client rendering to prevent focus repaint
	const dwmwaNcRenderingPolicy = 2
	const dwmncrpDisabled = 1
	dwmSetWindowAttribute.Call(hwnd, uintptr(dwmwaNcRenderingPolicy),
		uintptr(unsafe.Pointer(&[]uint32{dwmncrpDisabled}[0])), 4)
	fmt.Println("=== [DIAG] DWM NC rendering disabled ===")

	// Show window WITHOUT activation (SW_SHOWNOACTIVATE = 4)
	const swShowNoActivate = 4
	showWindow.Call(hwnd, uintptr(swShowNoActivate))
	fmt.Println("=== [DIAG] ShowWindow(SW_SHOWNOACTIVATE) called ===")

	// Update position to topmost (no need for SWP_SHOWWINDOW since already shown)
	setWindowPos.Call(hwnd, hwndTopmost, 0, 0, 0, 0,
		swpNoActivate|swpNoMove|swpNoSize|swpFrameChanged)

	// Force transparent background via runtime
	fmt.Println("=== [DIAG] Calling WindowSetBackgroundColour(0,0,0,0) ===")
	runtime.WindowSetBackgroundColour(a.ctx, 0, 0, 0, 0)
	fmt.Println("=== [DIAG] applyWindowStyles() finished ===")
}

func (a *App) GetInitialItems() []ListItem {
	return []ListItem{
		{ID: 1, Text: "聊天栏", Color: "#00ffff"},
		{ID: 2, Text: "控制台", Color: "#00ff88"},
		{ID: 3, Text: "调试模式", Color: "#ff66ff"},
		{ID: 4, Text: "快捷命令", Color: "#ffaa00"},
		{ID: 5, Text: "下载显示", Color: "#66aaff"},
		{ID: 6, Text: "鼠标透视", Color: "#ff5555"},
		{ID: 7, Text: "音乐信息", Color: "#aaff00"},
		{ID: 8, Text: "系统状态", Color: "#ffffff"},
	}
}

func (a *App) StartMockStream() {
	colors := []string{
		"#00ffff", "#00ff88", "#ff66ff", "#ffaa00",
		"#66aaff", "#ff5555", "#aaff00", "#ffffff",
	}

	sampleTexts := []string{
		"实时新增文本",
		"剪贴板历史",
		"调试日志条目",
		"系统通知消息",
		"网络状态更新",
	}

	go func() {
		textIdx := 0
		for {
			time.Sleep(1 * time.Second)

			id := int(a.nextID.Add(1))
			text := fmt.Sprintf("%s %d", sampleTexts[textIdx%len(sampleTexts)], id)
			color := colors[id%len(colors)]

			item := ListItem{
				ID:    id,
				Text:  text,
				Color: color,
			}

			runtime.EventsEmit(a.ctx, "ui:item:add", item)

			textIdx++
		}
	}()
}

func (a *App) SelectItem(id int, text string) {
	fmt.Printf("Selected item: id=%d, text=%s\n", id, text)
}
