<template>
  <div class="deployment-page">
    <PageHeader :title="$t('deployment.listTitle')">
      <div class="header-actions">
        <div class="filter-wrapper">
          <div class="select-wrapper">
            <select v-model="filterServerId" class="filter-select" @change="handleServerFilter">
              <option :value="undefined">{{ $t('deployment.allServers') }}</option>
              <option v-for="server in servers" :key="server.id" :value="server.id">
                {{ server.ip_address }}{{ server.server_alias ? ` (${server.server_alias})` : '' }}
              </option>
            </select>
            <Icon icon="lucide:chevron-down" class="select-icon" />
          </div>
        </div>
        <div class="search-wrapper">
          <Icon icon="lucide:search" class="search-icon" />
          <input
            v-model="searchKeyword"
            :placeholder="$t('deployment.searchPlaceholder')"
            class="search-input"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="searchKeyword = ''">
            <Icon icon="lucide:x" />
          </button>
        </div>
        <Button type="primary" @click="handleAdd">
          <Icon icon="lucide:plus" />
          {{ $t('deployment.addDeployment') }}
        </Button>
      </div>
    </PageHeader>

    <Loading v-if="loading" :text="$t('common.loading')" />

    <Empty v-else-if="filteredDeployments.length === 0" :description="$t('common.noData')" icon="lucide:rocket">
      <template #action>
        <Button type="primary" @click="handleAdd">
          {{ $t('deployment.addDeployment') }}
        </Button>
      </template>
    </Empty>

    <DeploymentTable
      v-else
      :deployments="filteredDeployments"
      @view="handleView"
      @start="handleStart"
      @stop="handleStop"
      @undeploy="handleUndeploy"
    />

    <DeploymentForm
      v-model:visible="showForm"
      :applications="applications"
      :servers="servers"
      @submit="handleFormSubmit"
    />

    <DeploymentDetail
      v-if="showDetail && viewingDeployment"
      :visible="showDetail"
      :deployment="viewingDeployment"
      @close="showDetail = false"
      @start="handleStartFromDetail"
      @stop="handleStopFromDetail"
      @undeploy="handleUndeployFromDetail"
    />

    <UndeployConfirm
      v-if="showUndeployConfirm && undeployingDeployment"
      :visible="showUndeployConfirm"
      :deployment="undeployingDeployment"
      @confirm="confirmUndeploy"
      @cancel="showUndeployConfirm = false"
    />

    <Toast
      :visible="toastVisible"
      :message="toastMessage"
      :type="toastType"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import { fetchDeployments, createDeployment, startDeployment, stopDeployment, undeployDeployment } from '@/api/deployment'
import { fetchApplications } from '@/api/application'
import { fetchServers } from '@/api/server'
import type { ApplicationInstance, Server } from '@/types'
import type { Deployment, CreateDeploymentRequest } from './types'
import PageHeader from '@/components/PageHeader/index.vue'
import Button from '@/components/Button/index.vue'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import Toast from '@/components/Toast/index.vue'
import DeploymentTable from './components/DeploymentTable.vue'
import DeploymentForm from './components/DeploymentForm.vue'
import DeploymentDetail from './components/DeploymentDetail.vue'
import UndeployConfirm from './components/UndeployConfirm.vue'
import { useLoading } from '@/composables/useLoading'

const { t } = useI18n()
const { loading, withLoading } = useLoading()

const deployments = ref<Deployment[]>([])
const applications = ref<ApplicationInstance[]>([])
const servers = ref<Server[]>([])
const showForm = ref(false)
const showUndeployConfirm = ref(false)
const showDetail = ref(false)
const viewingDeployment = ref<Deployment | null>(null)
const undeployingDeployment = ref<Deployment | null>(null)
const searchKeyword = ref('')
const filterServerId = ref<number | undefined>(undefined)
const toastVisible = ref(false)
const toastMessage = ref('')
const toastType = ref<'success' | 'error'>('success')

const filteredDeployments = computed(() => {
  let result = deployments.value || []

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(d =>
      d.application.name.toLowerCase().includes(keyword) ||
      d.deploy_id.toString().includes(keyword) ||
      d.server.ip_address.includes(keyword)
    )
  }

  return result
})

const loadDeployments = async () => {
  await withLoading(async () => {
    deployments.value = await fetchDeployments(filterServerId.value)
  })
}

const loadApplications = async () => {
  applications.value = await fetchApplications()
}

const loadServers = async () => {
  servers.value = await fetchServers()
}

const handleAdd = () => {
  showForm.value = true
}

const handleView = (deployment: Deployment) => {
  viewingDeployment.value = deployment
  showDetail.value = true
}

const handleStart = async (deployment: Deployment) => {
  try {
    await startDeployment(deployment.id)
    await loadDeployments()
    showToast(t('deployment.startSuccess'), 'success')
  } catch (error) {
    console.error('Failed to start deployment:', error)
    showToast(t('deployment.operationFailed'), 'error')
  }
}

const handleStop = async (deployment: Deployment) => {
  try {
    await stopDeployment(deployment.id)
    await loadDeployments()
    showToast(t('deployment.stopSuccess'), 'success')
  } catch (error) {
    console.error('Failed to stop deployment:', error)
    showToast(t('deployment.operationFailed'), 'error')
  }
}

const handleUndeploy = (deployment: Deployment) => {
  undeployingDeployment.value = deployment
  showUndeployConfirm.value = true
}

const handleStartFromDetail = async (deployment: Deployment) => {
  showDetail.value = false
  await handleStart(deployment)
}

const handleStopFromDetail = async (deployment: Deployment) => {
  showDetail.value = false
  await handleStop(deployment)
}

const handleUndeployFromDetail = (deployment: Deployment) => {
  showDetail.value = false
  handleUndeploy(deployment)
}

const confirmUndeploy = async () => {
  if (!undeployingDeployment.value) return

  try {
    await undeployDeployment(undeployingDeployment.value.id)
    showUndeployConfirm.value = false
    await loadDeployments()
    showToast(t('deployment.undeploySuccess'), 'success')
  } catch (error) {
    console.error('Failed to undeploy:', error)
    showToast(t('deployment.operationFailed'), 'error')
  }
}

const handleFormSubmit = async (data: CreateDeploymentRequest) => {
  try {
    await createDeployment(data.application_id, data)
    await loadDeployments()
    showToast(t('deployment.createSuccess'), 'success')
  } catch (error) {
    console.error('Failed to create deployment:', error)
    showToast(t('deployment.operationFailed'), 'error')
  }
}

const handleServerFilter = () => {
  loadDeployments()
}

const showToast = (message: string, type: 'success' | 'error') => {
  toastMessage.value = message
  toastType.value = type
  toastVisible.value = true
  setTimeout(() => {
    toastVisible.value = false
  }, 2000)
}

onMounted(() => {
  loadDeployments()
  loadApplications()
  loadServers()
})
</script>

<style scoped>
.deployment-page {
  padding: 20px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-wrapper {
  display: flex;
  align-items: center;
}

.select-wrapper {
  position: relative;
}

.filter-select {
  width: 180px;
  padding: 8px 32px 8px 12px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  color: #1e3a5f;
  background: #f8fafc;
  cursor: pointer;
  appearance: none;
  transition: all 0.2s ease;
}

.filter-select:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.select-icon {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  width: 14px;
  height: 14px;
  color: #64748b;
  pointer-events: none;
}

.search-wrapper {
  display: flex;
  align-items: center;
  position: relative;
}

.search-input {
  width: 240px;
  padding: 8px 12px 8px 36px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
  color: #1e3a5f;
  background: #f8fafc;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #4fc3f7;
  background: #ffffff;
}

.search-icon {
  position: absolute;
  left: 12px;
  width: 16px;
  height: 16px;
  color: #94a3b8;
  pointer-events: none;
}

.clear-btn {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  background: #f1f5f9;
  color: #64748b;
}
</style>
