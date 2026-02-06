<template>
  <div class="s-page-header">
    <div class="s-page-header-left">
      <h2 v-if="title" class="s-page-header-title">{{ title }}</h2>
      <slot name="title"></slot>
      <p v-if="subtitle" class="s-page-header-subtitle">{{ subtitle }}</p>
      <slot name="subtitle"></slot>
    </div>
    <div v-if="$slots.default || showBack" class="s-page-header-right">
      <Button v-if="showBack" type="ghost" @click="handleBack">
        <Icon icon="lucide:arrow-left" />
        {{ backText }}
      </Button>
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Icon } from '@iconify/vue'
import Button from '../Button/index.vue'

interface PageHeaderProps {
  title?: string
  subtitle?: string
  showBack?: boolean
  backText?: string
}

withDefaults(defineProps<PageHeaderProps>(), {
  title: '',
  subtitle: '',
  showBack: false,
  backText: '返回'
})

const router = useRouter()

const handleBack = () => {
  router.back()
}
</script>

<style scoped>
.s-page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.s-page-header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.s-page-header-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
  margin: 0;
}

.s-page-header-subtitle {
  font-size: 13px;
  color: #64748b;
  margin: 0;
}

.s-page-header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
