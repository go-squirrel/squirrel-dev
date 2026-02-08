<template>
  <div class="stats-section">
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-content">
          <span class="stat-value">{{ stats.script || 0 }}</span>
          <span class="stat-label">{{ $t('overview.script') }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-content">
          <span class="stat-value">{{ stats.deployment || 0 }}</span>
          <span class="stat-label">{{ $t('overview.deployment') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { BaseStats } from '../types'

defineProps<{
  stats: BaseStats
}>()
</script>

<style scoped>
.stats-section {
  background: #ffffff;
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-card {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  background: #f8fafc;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-value {
  font-size: 12px;
  font-weight: 700;
  color: #1e3a5f;
}

.stat-label {
  font-size: 14px;
  color: #64748b;
}
</style>
