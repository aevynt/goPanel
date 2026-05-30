<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useMessage, NButton, NTag, NInputNumber } from 'naive-ui'
import { listPorts, checkPort, findPort } from '@/api'
import type { PortInfo } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'

const message = useMessage()
const ports = ref<PortInfo[]>([])
const loading = ref(false)
const portCheck = ref<number | null>(null)
const portAvailable = ref<boolean | null>(null)
const suggestedPort = ref<number | null>(null)

async function load() {
  loading.value = true
  try {
    ports.value = await listPorts()
  } catch (err: any) {
    message.error('Failed to scan ports')
  }
  loading.value = false
}

async function handleCheck() {
  if (!portCheck.value) return
  try {
    const res = await checkPort(portCheck.value)
    portAvailable.value = res.available
  } catch {
    message.error('Check failed')
  }
}

async function handleFind() {
  try {
    const res = await findPort()
    suggestedPort.value = res.port
    message.success(`Suggested port: ${res.port}`)
  } catch {
    message.error('No ports available')
  }
}

const columns = [
  { title: 'Port', key: 'port', sorter: true, width: 100 },
  { title: 'Protocol', key: 'protocol', width: 80 },
  {
    title: 'State',
    key: 'state',
    render: (row: PortInfo) => h(NTag, { type: row.state === 'listening' ? 'success' : 'default', size: 'small' }, { default: () => row.state }),
    width: 100,
  },
]

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Ports">
      <template #actions>
        <n-button quaternary @click="handleFind">Find Available Port</n-button>
        <n-button type="primary" @click="load" :loading="loading">Scan</n-button>
      </template>
    </PageHeader>

    <div v-if="suggestedPort" class="alert">
      Suggested available port: <strong>{{ suggestedPort }}</strong>
    </div>

    <div class="card mb-4">
      <div class="card-title">Check Port</div>
      <n-space>
        <n-input-number v-model:value="portCheck" :min="1" :max="65535" placeholder="Enter port" class="w-40" />
        <n-button quaternary @click="handleCheck">Check</n-button>
        <span v-if="portAvailable !== null">
          <n-tag :type="portAvailable ? 'success' : 'error'" size="small">
            {{ portAvailable ? 'Available' : 'In Use' }}
          </n-tag>
        </span>
      </n-space>
    </div>

    <n-data-table
      :loading="loading"
      :columns="columns"
      :data="ports"
      :pagination="{ pageSize: 50 }"
      :single-line="false"
      :scroll-x="400"
    />
  </div>
</template>

<style scoped>
.alert {
  background: var(--claude-accent-muted, rgba(201, 100, 66, 0.08));
  color: var(--claude-accent, #c96442);
  border: 1px solid var(--claude-border-light, #f0eee6);
  border-radius: var(--radius, 8px);
  padding: 12px 16px;
  margin-bottom: 16px;
  font-size: 14px;
}

.card {
  transition: border-color 0.2s ease;
}

.card:hover {
  border-color: var(--claude-border, #3d3d3a);
}
</style>
