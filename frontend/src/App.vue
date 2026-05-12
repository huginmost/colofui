<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { main } from '../wailsjs/go/models'
import { GetInitialItems, StartMockStream, SetPanelBounds, SetContextMenuVisible } from '../wailsjs/go/main/App'
import { EventsOn, Quit } from '../wailsjs/runtime/runtime'

type ListItem = main.ListItem

const items = ref<ListItem[]>([])

// Panel drag state
const panelLeft = ref(0)
const panelTop = ref(0)
const isDragging = ref(false)
let dragStartX = 0
let dragStartY = 0
let dragPanelLeft = 0
let dragPanelTop = 0

function onPanelMouseDown(e: MouseEvent) {
  if (e.button !== 0) return
  isDragging.value = true
  dragStartX = e.clientX
  dragStartY = e.clientY
  dragPanelLeft = panelLeft.value
  dragPanelTop = panelTop.value
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
}

function onMouseMove(e: MouseEvent) {
  if (!isDragging.value) return
  panelLeft.value = dragPanelLeft + (e.clientX - dragStartX)
  panelTop.value = dragPanelTop + (e.clientY - dragStartY)
}

function onMouseUp() {
  if (!isDragging.value) return
  isDragging.value = false
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)
  SetPanelBounds(panelLeft.value, panelTop.value, 240, 460)
}

// Context menu state
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  item: null as ListItem | null,
})

function showContextMenu(e: MouseEvent, item: ListItem) {
  e.preventDefault()
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    item,
  }
  SetContextMenuVisible(true)
}

function hideContextMenu() {
  contextMenu.value.visible = false
  SetContextMenuVisible(false)
}

function editItem() {
  const target = contextMenu.value.item
  if (!target) return
  const item = items.value.find(i => i.id === target.id)
  if (item) {
    item.text = '已修改: ' + item.text
  }
  hideContextMenu()
}

function deleteItem() {
  const target = contextMenu.value.item
  if (!target) return
  items.value = items.value.filter(i => i.id !== target.id)
  hideContextMenu()
}

function onGlobalClick() {
  if (contextMenu.value.visible) {
    hideContextMenu()
  }
}

onMounted(async () => {
  items.value = await GetInitialItems()

  EventsOn('ui:item:add', (item: ListItem) => {
    items.value.unshift(item)

    if (items.value.length > 30) {
      items.value.pop()
    }
  })

  await StartMockStream()

  document.addEventListener('click', onGlobalClick)

  // Report initial panel bounds to Go backend for hit-test
  setTimeout(() => {
    SetPanelBounds(0, 0, 240, 460)
  }, 100)
})

onUnmounted(() => {
  document.removeEventListener('click', onGlobalClick)
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)
})
</script>

<template>
  <div
    class="app-panel"
    :style="{ left: panelLeft + 'px', top: panelTop + 'px' }"
    @mousedown="onPanelMouseDown"
  >
    <div class="scroll-area" @contextmenu.prevent>
      <div
        v-for="item in items"
        :key="item.id"
        class="menu-item item-enter"
        :style="{ color: item.color }"
        @contextmenu.prevent="showContextMenu($event, item)"
      >
        <span class="bar">|</span>
        <span class="item-text">{{ item.text }}</span>
      </div>
    </div>
  </div>

  <!-- Context menu backdrop -->
  <div v-if="contextMenu.visible" class="context-backdrop" @click="hideContextMenu" @contextmenu.prevent="hideContextMenu"></div>

  <!-- Context menu -->
  <div
    v-if="contextMenu.visible"
    class="context-menu"
    :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
  >
    <div class="context-menu-item" @click="editItem">
      <span class="context-icon">&#9998;</span>
      <span>修改</span>
    </div>
    <div class="context-menu-separator"></div>
    <div class="context-menu-item" @click="deleteItem">
      <span class="context-icon">&#10006;</span>
      <span>删除</span>
    </div>
    <div class="context-menu-separator"></div>
    <div class="context-menu-item context-menu-item--danger" @click="Quit()">
      <span class="context-icon">&#10005;</span>
      <span>关闭</span>
    </div>
  </div>
</template>

<style scoped>
.app-panel {
  position: fixed;
  width: 240px;
  height: 460px;
  pointer-events: auto;
}

.scroll-area {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 4px 6px;
  box-sizing: border-box;
  background: transparent;
}

.scroll-area::-webkit-scrollbar {
  display: none;
}

.menu-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 22px;
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 14px;
  line-height: 22px;
  text-shadow: 0 0 8px currentColor;
  background: rgba(0, 0, 0, 0.25);
  width: fit-content;
  max-width: 100%;
}

.menu-item:hover {
  background: rgba(20, 20, 20, 0.35);
}

.bar {
  flex-shrink: 0;
  text-shadow: 0 0 6px currentColor;
}

.item-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Fade-in animation for new items */
.item-enter {
  animation: fadeSlideIn 0.35s ease-out;
}

@keyframes fadeSlideIn {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Context menu backdrop */
.context-backdrop {
  position: fixed;
  inset: 0;
  z-index: 999;
  pointer-events: auto;
}

/* Context menu */
.context-menu {
  position: fixed;
  z-index: 1000;
  min-width: 140px;
  background: rgba(30, 30, 30, 0.88);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  padding: 4px 0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
  pointer-events: auto;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 4px;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  color: #ffffff;
  cursor: pointer;
  transition: background 0.15s;
}

.context-menu-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.context-menu-item--danger:hover {
  background: rgba(255, 70, 70, 0.25);
}

.context-icon {
  width: 16px;
  text-align: center;
  font-size: 12px;
  opacity: 0.7;
}

.context-menu-separator {
  height: 1px;
  margin: 3px 10px;
  background: rgba(255, 255, 255, 0.08);
}
</style>
