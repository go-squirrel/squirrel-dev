<template>
  <div class="server-selector">
    <label class="selector-label">{{ $t('monitor.selectServer') }}</label>
    <select v-model="selectedId" class="selector-select" @change="handleChange">
      <option value="">{{ $t('monitor.noServerSelected') }}</option>
      <option
        v-for="server in servers"
        :key="server.id"
        :value="server.id"
        :disabled="server.status === 'offline'"
      >
        {{ server.server_alias || server.hostname }} ({{ server.ip_address }})
      </option>
    </select>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Server } from '@/types'

interface Props {
  servers: Server[]
  modelValue: number | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: number | null]
  change: [serverId: number]
}>()

const selectedId = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const handleChange = () => {
  if (selectedId.value) {
    emit('change', selectedId.value)
  }
}
</script>

<style scoped lang="scss">
.server-selector {
  display: flex;
  align-items: center;
  gap: 12px;
}

.selector-label {
  font-size: 14px;
  color: #64748b;
  white-space: nowrap;
}

.selector-select {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  background: #fff;
  min-width: 200px;
  cursor: pointer;
  color: #1e3a5f;

  &:focus {
    outline: none;
    border-color: #4fc3f7;
  }
}
</style>
