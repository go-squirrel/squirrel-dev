<template>
  <div class="stats-section">
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon blue">
          <Icon icon="lucide:globe" />
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.website || 0 }}</span>
          <span class="stat-label">{{ $t('overview.website') }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon purple">
          <Icon icon="lucide:database" />
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.database || 0 }}</span>
          <span class="stat-label">{{ $t('overview.database') }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange">
          <Icon icon="lucide:clock" />
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.cron || 0 }}</span>
          <span class="stat-label">{{ $t('overview.cron') }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <Icon icon="lucide:package" />
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ stats.installedApps || 0 }}</span>
          <span class="stat-label">{{ $t('overview.installedApps') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { BaseStats } from '../types'

defineProps<{
  stats: BaseStats
}>()
</script>

<style scoped>
.stats-section {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 20px;
}

.stat-icon.blue {
  background: #e0f2fe;
  color: #0284c7;
}

.stat-icon.purple {
  background: #f3e8ff;
  color: #9333ea;
}

.stat-icon.orange {
  background: #ffedd5;
  color: #ea580c;
}

.stat-icon.green {
  background: #dcfce7;
  color: #16a34a;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1e3a5f;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
}
</style>
