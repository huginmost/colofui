# Clipboard UI Demo

基于 Wails v2 + Vue 3 + TypeScript + Go 的实时透明列表 UI Demo。

一个半透明悬浮窗口，模拟游戏调试菜单 / 插件菜单 / 剪贴板历史面板，实时接收后端推送的文本数据并以彩色列表显示。

## 最终效果

```
┌──────────────────────┐
│ SYSTEM MENU      ●    │
│                      │
│  ● 聊天栏             │
│  ● 控制台             │
│  ● 调试模式           │
│  ● 快捷命令           │
│  ● 下载显示           │
│  ● 鼠标透视           │
│  ● 音乐信息           │
│  ● 实时新增文本 1      │
│  ● 实时新增文本 2      │
│                      │
│  12 items             │
└──────────────────────┘
```

- 半透明黑色毛玻璃背景
- 列表文字竖直排列，每行独立颜色
- 圆点颜色与文字颜色一致
- 文本带发光 (glow) 效果
- 鼠标悬停高亮 + 右移动画
- 新条目淡入动画
- 点击条目通知 Go 后端
- 无边框透明窗口，240x460

## 技术栈

| 层 | 技术 |
|---|------|
| 前端 | Vue 3 + Vite + TypeScript |
| 桌面框架 | Wails v2 |
| 后端 | Go |
| 样式 | 原生 CSS（毛玻璃 + glow） |

## 环境要求

- Go 1.21+
- Node.js 18+
- Wails v2 CLI

### 安装 Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

确保 `$GOPATH/bin` 在 PATH 中。

### 检查依赖

```bash
wails doctor
```

Windows 下需要 WebView2 运行时（Win11 已内置，Win10 可能需要安装）。

## 快速运行

```bash
# 1. 进入项目目录
cd clipboard-ui-demo

# 2. 安装前端依赖
cd frontend
npm install
cd ..

# 3. 开发模式运行
wails dev
```

首次运行 `wails dev` 会自动生成 `wailsjs/` 绑定目录（包含 `go/main/App.js` 和 `runtime/runtime.js`），前端通过这些绑定调用 Go 后端。

## 构建可执行文件

```bash
wails build
```

产物在 `build/bin/` 目录下。

## 项目结构

```
clipboard-ui-demo/
├── go.mod                  # Go 模块定义
├── main.go                 # 入口 + Wails 窗口配置
├── app.go                  # Go 后端逻辑
├── wails.json              # Wails 项目配置
├── README.md
├── wailsjs/                # 自动生成的绑定 (wails dev 后生成)
│   ├── go/main/App.js
│   └── runtime/runtime.js
└── frontend/
    ├── package.json
    ├── index.html
    ├── tsconfig.json
    ├── vite.config.ts
    └── src/
        ├── main.ts         # Vue 入口
        ├── App.vue         # 主组件
        ├── style.css       # 全局样式
        ├── types.ts        # TypeScript 类型
        └── vite-env.d.ts   # 类型声明
```

## Go 后端 API

### GetInitialItems() []ListItem

返回初始 8 条列表数据，每条颜色不同。

### StartMockStream()

启动后台 goroutine，每秒推送一条模拟数据，事件名 `ui:item:add`。

### SelectItem(id int, text string)

前端点击条目时调用，Go 后端打印选中信息。

### ListItem 结构

```go
type ListItem struct {
    ID    int    `json:"id"`
    Text  string `json:"text"`
    Color string `json:"color"`
}
```

## 未来扩展：接入真实项目

### 当前 Demo 数据流

```
Go 后端模拟实时数据
        ↓
Wails runtime.EventsEmit("ui:item:add")
        ↓
Vue UI 实时显示
```

### 方案 A：Go 直接获取剪贴板

```
Go 监听系统剪贴板 (github.com/atotto/clipboard 等)
        ↓
Go 生成 ListItem
        ↓
Wails Event → Vue UI 显示
```

无需修改前端代码。

### 方案 B：C++ 获取剪贴板历史 → Go 接收 → UI 显示

```
C++ 程序获取剪贴板历史
        ↓
WebSocket / Named Pipe / TCP
        ↓
Go 后端接收 C++ 推送的数据
        ↓
Wails runtime.EventsEmit("ui:item:add", item)
        ↓
Vue UI 实时显示
```

**实现要点：**

1. 在 `app.go` 中增加一个 goroutine，启动 TCP/WebSocket/Named Pipe 服务端
2. C++ 程序连接到该端口，发送 JSON 数据
3. Go 解析 JSON 后通过 `runtime.EventsEmit` 推送到前端
4. Vue UI 不需要任何修改 — 它只关心 `ui:item:add` 事件的 `ListItem` 数据

**示例：通过 TCP 接收 C++ 数据**

```go
// 在 app.go 中添加
func (a *App) StartTCPServer(port string) {
    go func() {
        ln, _ := net.Listen("tcp", ":"+port)
        for {
            conn, _ := ln.Accept()
            go func(c net.Conn) {
                defer c.Close()
                var item ListItem
                json.NewDecoder(c).Decode(&item)
                runtime.EventsEmit(a.ctx, "ui:item:add", item)
            }(conn)
        }
    }()
}
```

**C++ 端发送示例：**

```cpp
// 伪代码 - 通过 TCP 发送 JSON
std::string json = R"({"id":42,"text":"剪贴板文本","color":"#00ffff"})";
send(sock, json.c_str(), json.size(), 0);
```

**关键点：**
- Vue UI 完全不依赖数据来源 — 只要收到 `{id, text, color}` 格式的 JSON，就能正常显示
- Go 后端是数据中枢 — C++、Python、剪贴板监听等都可以通过它统一转发
- 事件名称 `ui:item:add` 是前端的唯一契约

## Wails 窗口配置说明

`main.go` 中的窗口配置：

```go
&options.App{
    Title:             "Clipboard UI Demo",
    Width:             240,
    Height:            460,
    Frameless:         true,           // 无边框
    AlwaysOnTop:       true,           // 窗口置顶
    DisableResize:     true,           // 禁止调整大小
    BackgroundColour:  &options.RGBA{  // 背景完全透明
        R: 0, G: 0, B: 0, A: 0,
    },
    Windows: &windows.Options{
        WebviewIsTransparent: true,    // WebView 透明
        WindowIsTranslucent:  true,    // 窗口半透明
    },
}
```

Linux 下透明窗口效果可能因桌面环境不同而有差异（X11/Wayland 兼容性取决于 Wails 和 WebKit 版本）。
