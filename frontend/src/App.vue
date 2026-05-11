<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { main } from '../wailsjs/go/models'
import { GetInitialItems, StartMockStream, SelectItem } from '../wailsjs/go/main/App'
import { EventsOn } from '../wailsjs/runtime/runtime'

type ListItem = main.ListItem

const items = ref<ListItem[]>([])

onMounted(async () => {
  items.value = await GetInitialItems()

  EventsOn('ui:item:add', (item: ListItem) => {
    items.value.unshift(item)

    if (items.value.length > 30) {
      items.value.pop()
    }
  })

  await StartMockStream()
})

async function handleClick(item: ListItem) {
  await SelectItem(item.id, item.text)
}
</script>

<template>
  <div class="panel">
    <div class="title-bar">
      <span class="title-text">SYSTEM MENU</span>
      <span class="status-dot"></span>
    </div>

    <div class="scroll-area">
      <div
        v-for="item in items"
        :key="item.id"
        class="menu-item item-enter"
        :style="{ color: item.color }"
        @click="handleClick(item)"
      >
        <span class="dot"></span>
        <span class="item-text">{{ item.text }}</span>
      </div>
    </div>

    <div class="footer-bar">
      <span class="footer-text">{{ items.length }} items</span>
    </div>
  </div>
</template>

<style scoped>
.panel {
  width: 220px;
  height: 420px;
  padding: 10px;
  background: rgba(0, 0, 0, 0.38);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(120, 255, 255, 0.25);
  border-radius: 12px;
  box-shadow: 0 0 24px rgba(0, 255, 255, 0.12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.title-bar {
  --wails-draggable: drag;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 28px;
  padding: 0 6px;
  margin-bottom: 6px;
  flex-shrink: 0;
}

.title-text {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 2px;
  color: rgba(255, 255, 255, 0.6);
  text-shadow: 0 0 8px rgba(0, 255, 255, 0.3);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #00ff88;
  box-shadow: 0 0 6px #00ff88, 0 0 12px rgba(0, 255, 136, 0.4);
}

.scroll-area {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.scroll-area::-webkit-scrollbar {
  width: 3px;
}

.scroll-area::-webkit-scrollbar-track {
  background: transparent;
}

.scroll-area::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 22px;
  min-height: 22px;
  padding: 0 6px;
  border-radius: 6px;
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
  text-shadow: 0 0 6px currentColor;
  transition: background 0.15s ease, transform 0.15s ease;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.08);
  transform: translateX(2px);
}

.dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 6px currentColor;
  flex-shrink: 0;
}

.item-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.footer-bar {
  flex-shrink: 0;
  height: 22px;
  display: flex;
  align-items: center;
  padding: 0 6px;
  margin-top: 4px;
  border-top: 1px solid rgba(120, 255, 255, 0.08);
}

.footer-text {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 1px;
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
</style>
