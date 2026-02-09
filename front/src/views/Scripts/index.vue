<template>
  <div class="scripts-page">
    <PageHeader :title="$t('scripts.listTitle')">
      <div class="header-actions">
        <div class="search-wrapper">
          <Icon icon="lucide:search" class="search-icon" />
          <input
            v-model="searchKeyword"
            :placeholder="$t('scripts.searchPlaceholder')"
            class="search-input"
          />
          <button v-if="searchKeyword" class="clear-btn" @click="searchKeyword = ''">
            <Icon icon="lucide:x" />
          </button>
        </div>
        <Button type="primary" @click="handleAdd">
          <Icon icon="lucide:plus" />
          {{ $t('scripts.addScript') }}
        </Button>
      </div>
    </PageHeader>

    <Loading v-if="loading" :text="$t('common.loading')" />

    <Empty v-else-if="filteredScripts.length === 0" :description="$t('common.noData')" icon="lucide:file-code">
      <template #action>
        <Button type="primary" @click="handleAdd">
          {{ $t('scripts.addScript') }}
        </Button>
      </template>
    </Empty>

    <ScriptTable
      v-else
      :scripts="filteredScripts"
      @edit="handleEdit"
      @delete="handleDelete"
      @execute="handleExecute"
      @view-results="handleViewResults"
    />

    <ScriptEditor
      v-if="showEditor"
      :script="editingScript"
      @submit="handleFormSubmit"
      @cancel="showEditor = false"
    />

    <ExecuteDialog
      v-if="showExecuteDialog"
      :script="selectedScript"
      @execute="confirmExecute"
      @cancel="showExecuteDialog = false"
    />

    <ResultList
      v-if="showResultList"
      :script="selectedScript"
      :results="scriptResults"
      @close="showResultList = false"
      @view-log="handleViewLog"
    />

    <ExecutionLog
      v-if="showExecutionLog"
      :script="selectedScript"
      :result="selectedResult"
      @close="showExecutionLog = false"
      @refresh="refreshExecutionLog"
    />

    <DeleteConfirm
      v-if="showDeleteConfirm && deletingScript"
      :script="deletingScript"
      :deleting="deleting"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />

    <Toast
      :visible="toastVisible"
      :message="toastMessage"
      type="success"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useI18n } from 'vue-i18n'
import type { Script, CreateScriptRequest, UpdateScriptRequest, ScriptResult } from '@/types'
import PageHeader from '@/components/PageHeader/index.vue'
import Button from '@/components/Button/index.vue'
import Loading from '@/components/Loading/index.vue'
import Empty from '@/components/Empty/index.vue'
import Toast from '@/components/Toast/index.vue'
import ScriptTable from './components/ScriptTable.vue'
import ScriptEditor from './components/ScriptEditor.vue'
import ExecuteDialog from './components/ExecuteDialog.vue'
import ResultList from './components/ResultList.vue'
import ExecutionLog from './components/ExecutionLog.vue'
import DeleteConfirm from './components/DeleteConfirm.vue'
import { useScript } from './composables/useScript'

const { t } = useI18n()
const { scripts, loading, loadScripts, addScript, editScript, removeScript, runScript, getScriptResults } = useScript()

const showEditor = ref(false)
const showExecuteDialog = ref(false)
const showResultList = ref(false)
const showExecutionLog = ref(false)
const showDeleteConfirm = ref(false)
const searchKeyword = ref('')
const editingScript = ref<Script | null>(null)
const selectedScript = ref<Script | null>(null)
const selectedResult = ref<ScriptResult | null>(null)
const deletingScript = ref<Script | null>(null)
const scriptResults = ref<ScriptResult[]>([])
const deleting = ref(false)
const toastVisible = ref(false)
const toastMessage = ref('')

const filteredScripts = computed(() => {
  if (!searchKeyword.value) {
    return scripts.value
  }
  const keyword = searchKeyword.value.toLowerCase()
  return scripts.value.filter(script =>
    script.name.toLowerCase().includes(keyword)
  )
})

const handleAdd = () => {
  editingScript.value = null
  showEditor.value = true
}

const handleEdit = (script: Script) => {
  editingScript.value = script
  showEditor.value = true
}

const handleDelete = (script: Script) => {
  deletingScript.value = script
  showDeleteConfirm.value = true
}

const handleExecute = (script: Script) => {
  selectedScript.value = script
  showExecuteDialog.value = true
}

const handleViewResults = async (script: Script) => {
  selectedScript.value = script
  scriptResults.value = await getScriptResults(script.id)
  showResultList.value = true
}

const handleViewLog = (result: ScriptResult) => {
  selectedResult.value = result
  showExecutionLog.value = true
}

const refreshExecutionLog = async () => {
  if (!selectedScript.value || !selectedResult.value) return
  const results = await getScriptResults(selectedScript.value.id)
  const updatedResult = results.find(r => r.id === selectedResult.value?.id)
  if (updatedResult) {
    selectedResult.value = updatedResult
    // 更新结果列表中的数据
    const index = scriptResults.value.findIndex(r => r.id === updatedResult.id)
    if (index !== -1) {
      scriptResults.value[index] = updatedResult
    }
  }
}

const handleFormSubmit = async (data: CreateScriptRequest | UpdateScriptRequest) => {
  let success = false
  if (editingScript.value) {
    success = await editScript(editingScript.value.id, data as UpdateScriptRequest)
    if (success) {
      toastMessage.value = t('scripts.updateSuccess')
    }
  } else {
    success = await addScript(data as CreateScriptRequest)
    if (success) {
      toastMessage.value = t('scripts.createSuccess')
    }
  }

  if (success) {
    showEditor.value = false
    await loadScripts()
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2000)
  }
}

const confirmExecute = async (serverId: number) => {
  if (!selectedScript.value) return

  const success = await runScript({
    script_id: selectedScript.value.id,
    server_id: serverId
  })

  if (success) {
    showExecuteDialog.value = false
    toastMessage.value = t('scripts.executeSuccess')
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2000)
  }
}

const confirmDelete = async () => {
  if (!deletingScript.value) return

  deleting.value = true
  const success = await removeScript(deletingScript.value.id)
  deleting.value = false

  if (success) {
    showDeleteConfirm.value = false
    await loadScripts()
    toastMessage.value = t('scripts.deleteSuccess')
    toastVisible.value = true
    setTimeout(() => {
      toastVisible.value = false
    }, 2002)
  }
}

onMounted(() => {
  loadScripts()
})
</script>

<style scoped>
.scripts-page {
  padding: 20px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-wrapper {
  display: flex;
  align-items: center;
  position: relative;
}

.search-input {
  width: 280px;
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
