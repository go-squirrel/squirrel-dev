<template>
  <div class="s-empty">
    <Icon :icon="icon" class="empty-icon" />
    <p v-if="description" class="description">{{ description }}</p>
    <slot name="action"></slot>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

interface EmptyProps {
  icon?: string
  description?: string
}

withDefaults(defineProps<EmptyProps>(), {
  icon: 'lucide:inbox',
  description: ''
})
</script>

<style scoped>
.s-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: #64748b;
}

.empty-icon {
  width: 64px;
  height: 64px;
  color: #cbd5e1;
  margin-bottom: 16px;
}

.description {
  font-size: 14px;
  margin-bottom: 20px;
}
</style>
