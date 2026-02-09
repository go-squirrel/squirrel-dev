<template>
  <div 
    class="metric-item" 
    @mouseenter="$emit('showDetail')" 
    @mouseleave="$emit('hideDetail')"
  >
    <div class="metric-circle">
      <svg viewBox="0 0 100 100">
        <circle class="circle-bg" cx="50" cy="50" r="42" />
        <circle
          class="circle-progress"
          cx="50"
          cy="50"
          r="42"
          :stroke-dasharray="getCircleProgress(value)"
        />
      </svg>
      <div class="metric-value">
        <span class="value">{{ value.toFixed(2) }}</span>
        <span v-if="unit" class="unit">{{ unit }}</span>
      </div>
    </div>
    <span class="metric-label">{{ label }}</span>
    <span class="metric-sub">{{ subLabel }}</span>

    <!-- 详情插槽 -->
    <div v-if="showTooltip" class="tooltip-panel">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  value: number
  label: string
  subLabel: string
  unit?: string
  showTooltip?: boolean
}>()

defineEmits<{
  showDetail: []
  hideDetail: []
}>()

const getCircleProgress = (percentage: number) => {
  const radius = 42
  const circumference = 2 * Math.PI * radius
  const progress = Math.min(percentage, 100) / 100
  const dasharray = `${circumference * progress} ${circumference}`
  return dasharray
}
</script>

<style scoped>
.metric-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  position: relative;
  padding: 10px;
}

.metric-circle {
  position: relative;
  width: 80px;
  height: 80px;
}

.metric-circle svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.circle-bg {
  fill: none;
  stroke: #f1f5f9;
  stroke-width: 4;
}

.circle-progress {
  fill: none;
  stroke: url(#gradient);
  stroke-width: 4;
  stroke-linecap: round;
  transition: stroke-dasharray 0.5s ease;
}

.metric-value {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.metric-value .value {
  font-size: 16px;
  font-weight: 700;
  color: #1e3a5f;
}

.metric-value .unit {
  font-size: 12px;
  color: #64748b;
}

.metric-label {
  font-size: 13px;
  font-weight: 500;
  color: #1e3a5f;
}

.metric-sub {
  font-size: 11px;
  color: #94a3b8;
}

.tooltip-panel {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  width: 280px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  z-index: 100;
  margin-top: 8px;
}
</style>
