<template>
  <div class="category-filter">
    <div class="filter-group">
      <label class="filter-label">{{ $t('appStore.category') }}</label>
      <select v-model="selectedCategory" class="filter-select" @change="handleCategoryChange">
        <option value="">{{ $t('appStore.allCategory') }}</option>
        <option value="web">{{ $t('appStore.categoryWeb') }}</option>
        <option value="database">{{ $t('appStore.categoryDatabase') }}</option>
        <option value="middleware">{{ $t('appStore.categoryMiddleware') }}</option>
        <option value="devops">{{ $t('appStore.categoryDevops') }}</option>
      </select>
    </div>

    <div class="filter-group">
      <label class="filter-label">{{ $t('appStore.type') }}</label>
      <select v-model="selectedType" class="filter-select" @change="handleTypeChange">
        <option value="">{{ $t('appStore.allType') }}</option>
        <option value="compose">{{ $t('appStore.typeCompose') }}</option>
        <option value="k8s_manifest">{{ $t('appStore.typeK8s') }}</option>
        <option value="helm_chart">{{ $t('appStore.typeHelm') }}</option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  category?: string
  type?: string
}>()

const emit = defineEmits<{
  'update:category': [value: string]
  'update:type': [value: string]
}>()

const selectedCategory = ref(props.category || '')
const selectedType = ref(props.type || '')

watch(() => props.category, (newVal) => {
  selectedCategory.value = newVal || ''
})

watch(() => props.type, (newVal) => {
  selectedType.value = newVal || ''
})

const handleCategoryChange = () => {
  emit('update:category', selectedCategory.value)
}

const handleTypeChange = () => {
  emit('update:type', selectedType.value)
}
</script>

<style scoped>
.category-filter {
  display: flex;
  gap: 16px;
  align-items: center;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
}

.filter-select {
  padding: 8px 12px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 13px;
  color: #1e3a5f;
  background: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 140px;
}

.filter-select:focus {
  outline: none;
  border-color: #4fc3f7;
}

.filter-select:hover {
  border-color: #cbd5e1;
}
</style>
