<script setup lang="ts">
import { h, ref, computed, onMounted } from 'vue'
import { useMessage, NTag, NSpace, NButton, NPopconfirm, NTabPane, NTabs, NTooltip } from 'naive-ui'
import { listServices, listBinaries, createService, startService, stopService, restartService, removeService, getServiceLogs } from '@/api'
import type { Service, ServiceSpec, Binary } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import PageHeader from '@/components/PageHeader.vue'
import ServiceCreateModal from '@/components/ServiceCreateModal.vue'
import ServiceLogModal from '@/components/ServiceLogModal.vue'

const auth = useAuthStore()
const message = useMessage()
const services = ref<Service[]>([])
const binaries = ref<Binary[]>([])
const loading = ref(false)
const showCreateModal = ref(false)
const showLogModal = ref(false)
const currentLogs = ref<{ timestamp: string; message: string }[]>([])
const logLoading = ref(false)
const currentService = ref('')
const tabValue = ref('panel')

const panelServices = computed(() => services.value.filter((s: Service) => s.panel_managed))

const statusMap: Record<string, 'success' | 'default' | 'error' | 'warning'> = {
  active: 'success',
  inactive: 'default',
  failed: 'error',
}

function isSelf(row: Service) {
  return row.name.toLowerCase() === 'gopanel'
}

async function load() {
  loading.value = true
  try {
    const [svc, bins] = await Promise.all([listServices(), listBinaries()])
    services.value = svc
    binaries.value = bins
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to load services')
  }
  loading.value = false
}

async function handleCreate(spec: ServiceSpec) {
  try {
    await createService(spec)
    message.success('Service created')
    showCreateModal.value = false
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to create service')
  }
}

async function handleStart(name: string) {
  await startService(name)
  message.success('Service started')
  load()
}

async function handleStop(name: string) {
  await stopService(name)
  message.success('Service stopped')
  load()
}

async function handleRestart(name: string) {
  await restartService(name)
  message.success('Service restarted')
  load()
}

async function handleRemove(name: string) {
  await removeService(name)
  message.success('Service removed')
  load()
}

async function handleLogs(name: string) {
  currentService.value = name
  logLoading.value = true
  showLogModal.value = true
  try {
    currentLogs.value = await getServiceLogs(name, 100)
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to load logs')
  }
  logLoading.value = false
}

const actionColumn = {
  title: '',
  key: 'actions' as const,
  width: 280,
  render: (row: Service) => {
    const self = isSelf(row)
    const btns: any[] = []
    if (row.status === 'active') {
      btns.push(self
        ? h(NTooltip, { trigger: 'hover' }, { default: () => 'Cannot stop the panel itself', trigger: () => h(NButton, { size: 'tiny', quaternary: true, disabled: true }, { default: () => 'Stop' }) })
        : h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleStop(row.name) }, { default: () => 'Stop' }),
      )
    } else {
      btns.push(h(NButton, { size: 'tiny', quaternary: true, type: 'primary', onClick: () => handleStart(row.name) }, { default: () => 'Start' }))
    }
    btns.push(
      self
        ? h(NTooltip, { trigger: 'hover' }, { default: () => 'Cannot restart the panel itself', trigger: () => h(NButton, { size: 'tiny', quaternary: true, disabled: true }, { default: () => 'Restart' }) })
        : h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleRestart(row.name) }, { default: () => 'Restart' }),
    )
    btns.push(h(NButton, { size: 'tiny', quaternary: true, onClick: () => handleLogs(row.name) }, { default: () => 'Logs' }))
    if (auth.user?.role === 'admin') {
      btns.push(
        self
          ? h(NTooltip, { trigger: 'hover' }, { default: () => 'Cannot delete the panel itself', trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error', disabled: true }, { default: () => 'Delete' }) })
          : h(NPopconfirm, { onPositiveClick: () => handleRemove(row.name) }, {
              default: () => 'Delete this service?',
              trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => 'Delete' }),
            }),
      )
    }
    return h(NSpace, { size: 'small' }, btns)
  },
}

const columns = [
  {
    title: 'Name',
    key: 'name' as const,
    sorter: true,
    render: (row: Service) => {
      const children: any[] = [row.name]
      if (isSelf(row)) {
        children.push(' ')
        children.push(h(NTag, { size: 'tiny', type: 'info', bordered: false }, { default: () => 'Panel' }))
      }
      return children
    },
  },
  {
    title: 'Status',
    key: 'status' as const,
    render: (row: Service) => h(NTag, { type: (statusMap[row.status] || 'default') as any, size: 'small' }, { default: () => row.status }),
  },
  {
    title: 'Port',
    key: 'port' as const,
    render: (row: Service) => row.port?.toString() || '-',
  },
  actionColumn,
]

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Services">
      <template #actions>
        <n-button v-if="auth.user?.role === 'admin'" type="primary" @click="showCreateModal = true">
          Create Service
        </n-button>
      </template>
    </PageHeader>

    <n-tabs v-model:value="tabValue" type="line">
      <n-tab-pane name="panel" tab="Panel Services">
        <n-data-table
          :loading="loading"
          :columns="columns"
          :data="panelServices"
          :pagination="{ pageSize: 20 }"
          :single-line="false"
          :scroll-x="600"
        />
      </n-tab-pane>
      <n-tab-pane name="all" tab="All Services">
        <n-data-table
          :loading="loading"
          :columns="columns"
          :data="services"
          :pagination="{ pageSize: 20 }"
          :single-line="false"
          :scroll-x="600"
        />
      </n-tab-pane>
    </n-tabs>

    <ServiceCreateModal
      v-model:show="showCreateModal"
      :binaries="binaries"
      @create="handleCreate"
    />

    <ServiceLogModal
      v-model:show="showLogModal"
      :service-name="currentService"
      :logs="currentLogs"
      :loading="logLoading"
    />
  </div>
</template>
