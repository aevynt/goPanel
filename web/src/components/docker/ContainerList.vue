<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useMessage, NTag, NSpace, NButton, NDataTable } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { listContainers, startContainer, stopContainer, restartContainer } from '@/api'
import type { DockerContainer } from '@/api/client'
import LogsModal from './LogsModal.vue'

const message = useMessage()
const containers = ref<DockerContainer[]>([])
const loading = ref(false)

// Logs modal state
const showLogs = ref(false)
const logContainerId = ref('')
const logContainerName = ref('')

const stateTypeMap: Record<string, 'success' | 'default' | 'error' | 'warning'> = {
  running: 'success',
  exited: 'default',
  paused: 'warning',
  dead: 'error',
}

async function load() {
  loading.value = true
  try {
    containers.value = await listContainers()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to load Docker containers')
  }
  loading.value = false
}

async function handleStart(id: string) {
  try {
    await startContainer(id)
    message.success('Container started')
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to start container')
  }
}

async function handleStop(id: string) {
  try {
    await stopContainer(id)
    message.success('Container stopped')
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to stop container')
  }
}

async function handleRestart(id: string) {
  try {
    await restartContainer(id)
    message.success('Container restarted')
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to restart container')
  }
}

function handleLogs(id: string, name: string) {
  logContainerId.value = id
  logContainerName.value = name
  showLogs.value = true
}

const columns: DataTableColumns<DockerContainer> = [
  {
    title: 'Name',
    key: 'names',
    sorter: 'default',
  },
  {
    title: 'Image',
    key: 'image',
  },
  {
    title: 'State',
    key: 'state',
    render: (row: DockerContainer) => h(NTag, { type: stateTypeMap[row.state] || 'default', size: 'small', bordered: false }, { default: () => row.state }),
  },
  {
    title: 'Status',
    key: 'status',
  },
  {
    title: 'Ports',
    key: 'ports',
    render: (row: DockerContainer) => row.ports || '-',
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 250,
    render: (row: DockerContainer) => {
      const btns: any[] = []
      if (row.state === 'running') {
        btns.push(h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleStop(row.id) }, { default: () => 'Stop' }))
      } else {
        btns.push(h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => handleStart(row.id) }, { default: () => 'Start' }))
      }
      btns.push(h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleRestart(row.id) }, { default: () => 'Restart' }))
      btns.push(h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleLogs(row.id, row.names) }, { default: () => 'Logs' }))
      return h(NSpace, { size: 'small' }, btns)
    },
  },
]

defineExpose({
  load
})

onMounted(load)
</script>

<template>
  <n-space vertical :size="16">
    <n-space justify="end">
      <n-button secondary size="small" @click="load" :loading="loading">
        Refresh
      </n-button>
    </n-space>
    <n-data-table
      :loading="loading"
      :columns="columns"
      :data="containers"
      :pagination="{ pageSize: 15 }"
      :single-line="false"
      :scroll-x="600"
    />

    <LogsModal
      v-model:show="showLogs"
      :container-id="logContainerId"
      :container-name="logContainerName"
    />
  </n-space>
</template>
