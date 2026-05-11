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
</template>

<style scoped>
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
