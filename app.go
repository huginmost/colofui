package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

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
