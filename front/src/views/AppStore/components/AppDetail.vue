<template>
  <div class="app-detail-overlay" @click.self="$emit('close')">
    <div class="app-detail-container">
      <div class="detail-header">
        <div class="header-info">
          <img v-if="app?.icon" :src="app.icon" class="detail-icon" />
          <Icon v-else icon="lucide:package" class="detail-icon-default" />
          <div class="header-text">
            <h3 class="detail-title">
              {{ app?.name }}
              <span v-if="app?.is_official" class="official-badge">
                {{ $t('appStore.isOfficial') }}
              </span>
            </h3>
            <p class="detail-description">{{ app?.description }}</p>
          </div>
        </div>
        <button class="close-btn" @click="$emit('close')">
          <Icon icon="lucide:x" />
        </button>
      </div>

      <div class="detail-body">
        <!-- 基本信息 -->
        <div class="detail-section">
          <h4 class="section-title">{{ $t('appStore.basicInfo') }}</h4>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">{{ $t('appStore.type') }}</span>
              <span class="info-value">
                <span class="type-tag" :class="`type-${app?.type}`">
                  {{ getTypeLabel(app?.type || '') }}
                </span>
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('appStore.category') }}</span>
              <span class="info-value">{{ getCategoryLabel(app?.category || '') }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('appStore.version') }}</span>
              <span class="info-value version">{{ app?.version }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('appStore.author') }}</span>
              <span class="info-value">{{ app?.author }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('appStore.downloads') }}</span>
              <span class="info-value">
                <Icon icon="lucide:download" class="inline-icon" />
                {{ app?.downloads }}
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ $t('appStore.status') }}</span>
              <span class="info-value">
                <span class="status-badge" :class="`status-${app?.status}`">
                  {{ getStatusLabel(app?.status || '') }}
                </span>
              </span>
            </div>
          </div>
        </div>

        <!-- 标签 -->
        <div v-if="app?.tags" class="detail-section">
          <h4 class="section-title">{{ $t('appStore.tags') }}</h4>
          <div class="tags-container">
            <span v-for="tag in tagList" :key="tag" class="tag-item">{{ tag }}</span>
          </div>
        </div>

        <!-- 链接 -->
        <div v-if="app?.repo_url || app?.homepage_url" class="detail-section">
          <h4 class="section-title">{{ $t('appStore.repoUrl') }}</h4>
          <div class="links-container">
            <a v-if="app?.repo_url" :href="app.repo_url" target="_blank" class="link-item">
              <Icon icon="lucide:github" class="link-icon" />
              <span>{{ $t('appStore.repoUrl') }}</span>
              <Icon icon="lucide:external-link" class="external-icon" />
            </a>
            <a v-if="app?.homepage_url" :href="app.homepage_url" target="_blank" class="link-item">
              <Icon icon="lucide:globe" class="link-icon" />
              <span>{{ $t('appStore.homepageUrl') }}</span>
              <Icon icon="lucide:external-link" class="external-icon" />
            </a>
          </div>
        </div>

        <!-- 模板内容 -->
        <div class="detail-section">
          <div class="content-header">
            <h4 class="section-title">{{ $t('appStore.templateContent') }}</h4>
            <button class="copy-btn" @click="handleCopy">
              <Icon icon="lucide:copy" />
              {{ $t('appStore.copy') }}
            </button>
          </div>
          <pre class="content-code"><code>{{ app?.content }}</code></pre>
        </div>
      </div>

      <div class="detail-footer">
        <Button type="secondary" @click="$emit('close')">
          {{ $t('common.close') }}
        </Button>
        <Button type="secondary" icon="lucide:download-cloud" @click="$emit('import', app)">
          {{ $t('appStore.importToApplication') }}
        </Button>
        <Button type="primary" icon="lucide:download" @click="$emit('download', app)">
          {{ $t('appStore.download') }}
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import Button from '@/components/Button/index.vue'
import type { AppStore } from '@/types'

const props = defineProps<{
  app: AppStore | null
}>()

defineEmits<{
  close: []
  download: [app: AppStore | null]
  import: [app: AppStore | null]
}>()

const { t } = useI18n()

const tagList = computed(() => {
  return props.app?.tags?.split(',').map(tag => tag.trim()).filter(Boolean) || []
})

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    compose: t('appStore.typeCompose'),
    k8s_manifest: t('appStore.typeK8s'),
    helm_chart: t('appStore.typeHelm')
  }
  return labels[type] || type
}

const getCategoryLabel = (category: string) => {
  const labels: Record<string, string> = {
    web: t('appStore.categoryWeb'),
    database: t('appStore.categoryDatabase'),
    middleware: t('appStore.categoryMiddleware'),
    devops: t('appStore.categoryDevops')
  }
  return labels[category] || category
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    active: t('appStore.statusActive'),
    deprecated: t('appStore.statusDeprecated')
  }
  return labels[status] || status
}

const handleCopy = async () => {
  if (!props.app?.content) return
  try {
    await navigator.clipboard.writeText(props.app.content)
    alert(t('appStore.copySuccess'))
  } catch (err) {
    console.error('Copy failed:', err)
  }
}
</script>

<style scoped>
.app-detail-overlay {
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

.app-detail-container {
  background: #ffffff;
  border-radius: 16px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.detail-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 24px;
  border-bottom: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.header-info {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.detail-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  object-fit: cover;
  flex-shrink: 0;
}

.detail-icon-default {
  width: 56px;
  height: 56px;
  color: #94a3b8;
  flex-shrink: 0;
}

.header-text {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.detail-title {
  font-size: 20px;
  font-weight: 600;
  color: #1e3a5f;
  display: flex;
  align-items: center;
  gap: 10px;
}

.official-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: #ffffff;
  font-size: 11px;
  font-weight: 600;
  border-radius: 4px;
  text-transform: uppercase;
}

.detail-description {
  font-size: 14px;
  color: #64748b;
  line-height: 1.5;
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.close-btn:hover {
  background: #f1f5f9;
  color: #1e3a5f;
}

.detail-body {
  padding: 24px;
  overflow-y: auto;
  max-height: calc(90vh - 180px);
}

.detail-section {
  margin-bottom: 28px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a5f;
  margin-bottom: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-label {
  font-size: 12px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-size: 14px;
  color: #1e3a5f;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.info-value.version {
  font-family: 'SF Mono', Monaco, Consolas, monospace;
}

.inline-icon {
  width: 14px;
  height: 14px;
  color: #94a3b8;
}

.type-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.type-compose {
  background: #dbeafe;
  color: #1d4ed8;
}

.type-k8s_manifest {
  background: #dcfce7;
  color: #15803d;
}

.type-helm_chart {
  background: #f3e8ff;
  color: #7c3aed;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.status-active {
  background: #dcfce7;
  color: #15803d;
}

.status-deprecated {
  background: #f1f5f9;
  color: #64748b;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  background: #f1f5f9;
  color: #475569;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.links-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.link-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: #f8fafc;
  border-radius: 8px;
  color: #0284c7;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.2s ease;
}

.link-item:hover {
  background: #eff6ff;
}

.link-icon {
  width: 18px;
  height: 18px;
}

.external-icon {
  width: 14px;
  height: 14px;
  margin-left: auto;
  color: #94a3b8;
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.content-header .section-title {
  margin-bottom: 0;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f1f5f9;
  border: none;
  border-radius: 6px;
  color: #64748b;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  background: #e2e8f0;
  color: #1e3a5f;
}

.content-code {
  background: #1e293b;
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  margin: 0;
}

.content-code code {
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  color: #e2e8f0;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.detail-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
  flex-shrink: 0;
}
</style>
