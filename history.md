# Changelog

## v1.1 (2026-05-12)

**latest** — Replace dot with `|`, item background fits text width
- 行首圆点替换为竖线 `|`，保留 glow 效果
- 文本栏底部阴影框宽度改为跟随文本长度（`inline-flex` + `width: fit-content`）

**063888e** — Window-wide drag
- `--wails-draggable: drag` 加到 `scroll-area`，窗口内任意位置均可拖动

**b5e8549** — Hide scrollbar, lighter item background, larger border-radius
- 隐藏滚动条（保留滚轮滚动）
- 每行文本背景透明度降低：`rgba(0,0,0,0.45)` → `rgba(0,0,0,0.25)`
- hover 背景同步变淡
- 圆角加大：`4px` → `10px`

**b4f0e94** — Full transparency, fix focus opacity bug, remove glassmorphism
- 窗口查找方式从 `FindWindowW`（按标题）改为 `EnumWindows`（按进程 ID），更可靠
- 窗口显示后才调用 `runtime.WindowSetBackgroundColour` 强制透明
- 去掉每行文本的 `backdrop-filter` 毛玻璃效果，改为纯深色半透明背景
- `scroll-area` 显式设 `background: transparent`

**89b7ba2** — Fix: start hidden, remove X button, fix initial transparency
- 窗口启动时先隐藏，样式修改完成后再显示（解决启动时短暂不透明的 bug）
- 移除 `WS_SYSMENU` 样式去掉右上角 X 按钮
- `applyWindowStyles()` 加入轮询等待窗口句柄就绪

**472f6b8** — No-focus click-through window, per-item glassmorphism, right-click quit
- 通过 `SetWindowLongW` 添加 `WS_EX_NOACTIVATE` + `WS_EX_TOOLWINDOW`，窗口不获取焦点、不显示在任务栏
- 移除左键点击处理
- 右键任意位置调用 `Quit()` 关闭程序
- 每行文本底部添加独立毛玻璃效果：`background: rgba(0,0,0,0.35)` + `backdrop-filter: blur(6px)`

**83232ac** — Remove window chrome: title bar, footer, panel background
- 去掉标题栏（"SYSTEM MENU" + 状态灯）、底部计数栏、面板毛玻璃外壳
- 只保留纯文本列表，`scroll-area` 直接撑满窗口

## v1.0 (2026-05-11)

**e90a831** — Initial project: Wails v2 + Vue 3 transparent overlay UI
- Go 后端：`ListItem` 结构体、`GetInitialItems()`（8 条不同颜色的初始数据）、`StartMockStream()`（每秒推送一条模拟数据）、`SelectItem()`（打印选中信息）
- Vue 3 前端：透明悬浮窗口 240x460，彩色文本列表带圆点、glow 文字阴影、hover 高亮、淡入动画
- 窗口配置：`Frameless` + `AlwaysOnTop` + `BackgroundColour: rgba(0,0,0,0)` + `WebviewIsTransparent` + `WindowIsTranslucent`
- Wails 事件系统：Go 端 `runtime.EventsEmit("ui:item:add")` → Vue 端 `EventsOn` 接收
- 列表最大 30 条，新条目插入顶部
