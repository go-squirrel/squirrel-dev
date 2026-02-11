<template>
  <div class="application-detail-overlay" @click.self="$emit('close')">
    <div class="application-detail-container">
      <div class="detail-header">
        <h3 class="detail-title">
          {{ $t('application.applicationDetail') }}
        </h3>
        <button class="close-btn" @click="$emit('close')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="detail-body">
        <div class="detail-section">
          <h4 class="section-title">{{ $t('application.basicInfo') }}</h4>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">{{ $t('application.applicationName') }}</span>
              <span class="info-value">{{ application?.name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('application.applicationType') }}</span>
              <span class="info-value type-badge" :class="`type-badge--${application?.type}`">
                {{ getTypeLabel(application?.type) }}
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('application.applicationVersion') }}</span>
              <span class="info-value">{{ application?.version }}</span>
            </div>
            <div class="info-item full-width">
              <span class="info-label">{{ $t('application.applicationDescription') }}</span>
              <span class="info-value">{{ application?.description || '-' }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h4 class="section-title">{{ $t('application.contentPreview') }}</h4>
          <div class="content-preview">
            <pre class="code-content">{{ application?.content }}</pre>
          </div>
        </div>
      </div>

      <div class="detail-footer">
        <Button type="secondary" @click="$emit('close')">
          {{ $t('common.close') }}
        </Button>
        <Button type="primary" @click="handleEdit">
          <Icon icon="lucide:edit-2" />
          {{ $t('application.editApplication') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from '@/components/Button/index.vue'
import type { ApplicationInstance } from '@/types'

const props = defineProps<{
  application: ApplicationInstance | null
}>()

const emit = defineEmits<{
  close: []
  edit: [application: ApplicationInstance]
}>()

const getTypeLabel = (type?: string) => {
  if (!type) return '-'
  const labels: Record<string, string> = {
    compose: 'Compose',
    script: '脚本',
    binary: '二进制'
  }
  return labels[type] || type
}

const handleEdit = () => {
  if (props.application) {
    emit('edit', props.application)
  }
}
</script>

<style scoped>
.application-detail-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.application-detail-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.detail-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e3a5f;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: #f1f5f9;
  color: #1e3a5f;
}

.detail-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.detail-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-label {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #1e3a5f;
  font-weight: 500;
}

.type-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.type-badge--compose {
  background: #dbeafe;
  color: #1e40af;
}

.type-badge--script {
  background: #dcfce7;
  color: #166534;
}

.type-badge--binary {
  background: #fef3c7;
  color: #92400e;
}

.content-preview {
  background: #1e293b;
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
}

.code-content {
  margin: 0;
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 12px;
  color: #e2e8f0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.detail-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}
</style>
