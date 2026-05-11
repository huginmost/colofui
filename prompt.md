# Vue + Wails + Go 透明文字列表 UI Demo Prompt

## 任务目标

创建一个完整可运行的桌面端 Demo，用于展示透明背景上的实时文字列表。它不需要传统窗口面板、边框或复杂玻璃卡片，只需要让文字像 Overlay 一样悬浮在桌面上。

技术栈:

- 前端：Vue 3 + Vite + TypeScript
- 桌面框架：Wails v2
- 后端：Go
- 样式：普通 CSS
- 目标平台：Windows 优先

## UI 要求

- 窗口透明、无边框、默认置顶
- 不做传统窗口外观，不使用大面板或卡片
- 以文字列表为主
- 每条文本有独立颜色
- 每条文本底部有半透明深色底条，增强可读性
- 文本可左键选择
- 文本可右键打开上下文菜单
- 上下文菜单至少包含：
  - 修改文本
  - 复制文本
  - 关闭此项
  - 只保留此项
- 列表最多保留 30 条
- 新增文本从顶部插入

## Go 后端要求

定义数据结构:

```go
type ListItem struct {
    ID    int    `json:"id"`
    Text  string `json:"text"`
    Color string `json:"color"`
}
```

实现方法:

- `GetInitialItems() []ListItem`
- `StartMockStream()`
- `SelectItem(id int, text string)`
- `UpdateItem(id int, text string) (ListItem, bool)`
- `CloseItem(id int) bool`
- `CloseOtherItems(id int) []ListItem`

`StartMockStream` 每秒生成一条模拟文本，并通过 Wails event 推送：

```text
ui:item:add
```

## Vue 前端要求

- `frontend/src/types.ts` 定义 `ListItem`
- 页面加载时调用 `GetInitialItems()`
- 页面加载后调用 `StartMockStream()`
- 监听 `ui:item:add`
- 收到新项后插入顶部
- 超过 30 条删除末尾
- 左键调用 `SelectItem(id, text)`
- 右键打开菜单并调用后端修改/关闭方法

## README 要求

README 需要说明:

- 如何运行 `wails dev`
- 如何构建 `wails build`
- 当前 Go 模拟数据到 Vue UI 的数据流
- 后续如何接入 C++ 项目，例如 C++ 通过 Named Pipe、TCP 或 WebSocket 把剪贴板历史传给 Go
