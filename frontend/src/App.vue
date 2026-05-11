<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { main } from '../wailsjs/go/models'
import { GetInitialItems, StartMockStream } from '../wailsjs/go/main/App'
import { EventsOn, Quit } from '../wailsjs/runtime/runtime'

type ListItem = main.ListItem

const items = ref<ListItem[]>([])

onMounted(async () => {
  console.log('=== [DIAG] Vue mounted ===')
  console.log('html bg:', getComputedStyle(document.documentElement).background)
  console.log('body bg:', getComputedStyle(document.body).background)
  console.log('#app bg:', getComputedStyle(document.getElementById('app')!).background)

  items.value = await GetInitialItems()
  console.log('=== [DIAG] Initial items loaded:', items.value.length)

  EventsOn('ui:item:add', (item: ListItem) => {
    items.value.unshift(item)

    if (items.value.length > 30) {
      items.value.pop()
    }
  })

  await StartMockStream()
  console.log('=== [DIAG] Mock stream started ===')
})
</script>

<template>
  <div class="scroll-area" @contextmenu.prevent="Quit()">
    <div
      v-for="item in items"
      :key="item.id"
      class="menu-item item-enter"
      :style="{ color: item.color }"
      @contextmenu.prevent="Quit()"
    >
      <span class="bar">|</span>
      <span class="item-text">{{ item.text }}</span>
    </div>
  </div>
</template>

<style scoped>
.scroll-area {
  --wails-draggable: drag;
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
</style>
