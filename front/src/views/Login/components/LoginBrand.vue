<template>
  <div class="brand">
    <Icon icon="lucide:squirrel" class="brand-icon" />
    <h1 class="brand-title">Squirrel</h1>
    <p class="brand-subtitle">{{ $t('login.subtitle') }}</p>
  </div>
  <div class="features">
    <div class="feature-item">
      <Icon icon="lucide:zap" class="feature-icon" />
      <span>{{ $t('login.feature1') }}</span>
    </div>
    <div class="feature-item">
      <Icon icon="lucide:shield-check" class="feature-icon" />
      <span>{{ $t('login.feature2') }}</span>
    </div>
    <div class="feature-item">
      <Icon icon="lucide:bar-chart-3" class="feature-icon" />
      <span>{{ $t('login.feature3') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
</script>

<style scoped>
.brand {
  text-align: center;
  margin-bottom: 60px;
  z-index: 1;
}

.brand-icon {
  width: 64px;
  height: 64px;
  color: #4fc3f7;
  margin-bottom: 16px;
}

.brand-title {
  font-size: 36px;
  font-weight: 700;
  color: #ffffff;
  margin-bottom: 8px;
  letter-spacing: 1px;
}

.brand-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.6);
}

.features {
  display: flex;
  flex-direction: column;
  gap: 20px;
  z-index: 1;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.feature-icon {
  width: 20px;
  height: 20px;
  color: #4fc3f7;
}
</style>
