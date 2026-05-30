<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useMessage, NButton, NPopconfirm, NTag, NRadioGroup, NRadio } from 'naive-ui'
import { listSites, addSite, removeSite, caddyHealth } from '@/api'
import type { Site } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'
import ModalFooter from '@/components/ModalFooter.vue'

const message = useMessage()
const sites = ref<Site[]>([])
const loading = ref(false)
const caddyOnline = ref(false)
const showAddModal = ref(false)
const siteType = ref<'proxy' | 'static'>('proxy')

const newSite = ref<Site>({
  domain: '',
  service_port: 8080,
  tls_enabled: true,
  tls_email: '',
  extra_config: '',
  type: 'proxy',
  root: '',
})

async function load() {
  loading.value = true
  try {
    sites.value = await listSites()
    caddyOnline.value = true
  } catch {
    caddyOnline.value = false
  }
  loading.value = false
}

async function checkHealth() {
  try {
    await caddyHealth()
    caddyOnline.value = true
    message.success('Caddy is online')
  } catch {
    caddyOnline.value = false
    message.error('Caddy is not reachable')
  }
}

function handleTypeChange(val: 'proxy' | 'static') {
  siteType.value = val
  newSite.value.type = val
}

async function handleAdd() {
  try {
    await addSite(newSite.value)
    message.success('Site added')
    showAddModal.value = false
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to add site')
  }
}

async function handleRemove(domain: string) {
  try {
    await removeSite(domain)
    message.success('Site removed')
    load()
  } catch (err: any) {
    message.error('Failed to remove site')
  }
}

const columns = [
  { title: 'Domain', key: 'domain' },
  {
    title: 'Type',
    key: 'type',
    width: 90,
    render: (row: Site) => {
      const isStatic = row.type === 'static'
      return h(NTag, {
        size: 'small',
        type: isStatic ? 'success' : 'info',
      }, { default: () => isStatic ? 'Static' : 'Proxy' })
    },
  },
  {
    title: 'Target',
    key: 'service_port',
    render: (row: Site) => {
      if (row.type === 'static') {
        return h('span', { style: { fontSize: '13px', color: '#9a9990' } }, row.root || '')
      }
      return h('span', { style: { fontSize: '13px' } }, `:${row.service_port}`)
    },
  },
  {
    title: 'TLS',
    key: 'tls_enabled',
    width: 70,
    render: (row: Site) => {
      const color = row.tls_enabled ? '#16a34a' : '#5e5d59'
      return h('span', { style: { color, fontSize: '13px' } }, row.tls_enabled ? 'Yes' : 'No')
    },
  },
  {
    title: '',
    key: 'actions',
    width: 80,
    render: (row: Site) => h(NPopconfirm, { onPositiveClick: () => handleRemove(row.domain) }, {
      default: () => `Delete site ${row.domain}?`,
      trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => 'Delete' }),
    }),
  },
]

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Sites">
      <template #actions>
        <n-tag :type="caddyOnline ? 'success' : 'error'" size="small">
          {{ caddyOnline ? 'Online' : 'Offline' }}
        </n-tag>
        <n-button quaternary size="tiny" @click="checkHealth">Check</n-button>
        <n-button type="primary" @click="showAddModal = true">Add Site</n-button>
      </template>
    </PageHeader>

    <n-data-table
      :loading="loading"
      :columns="columns"
      :data="sites"
      :pagination="{ pageSize: 20 }"
      :single-line="false"
      :scroll-x="600"
    />

    <n-modal v-model:show="showAddModal" title="Add Site" preset="card" style="width: 520px; max-width: 90vw;" :mask-closable="false">
      <n-form>
        <n-form-item label="Domain">
          <n-input v-model:value="newSite.domain" placeholder="example.com" />
        </n-form-item>
        <n-form-item label="Type">
          <n-radio-group :value="siteType" @update:value="handleTypeChange">
            <n-radio value="proxy" style="color:#e8e6dc;">Reverse Proxy</n-radio>
            <n-radio value="static" style="color:#e8e6dc;">Static Site</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="siteType === 'proxy'" label="Service Port">
          <n-input-number v-model:value="newSite.service_port" :min="1" :max="65535" />
        </n-form-item>
        <n-form-item v-if="siteType === 'static'" label="Site Directory">
          <n-input :value="`public-sites/${newSite.domain}/`" disabled placeholder="auto" />
        </n-form-item>
        <n-checkbox v-model:checked="newSite.tls_enabled">
          Enable TLS (Let's Encrypt)
        </n-checkbox>
        <n-form-item v-if="newSite.tls_enabled" label="TLS Email">
          <n-input v-model:value="newSite.tls_email" placeholder="admin@example.com" />
        </n-form-item>
      </n-form>
      <template #footer>
        <ModalFooter submit-label="Add" @cancel="showAddModal = false" @submit="handleAdd" />
      </template>
    </n-modal>
  </div>
</template>
