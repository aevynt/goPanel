<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NButton, NPopconfirm, NInput, NModal } from 'naive-ui'
import { listShares, createShare, deleteShare, getPublicDomain, setPublicDomain } from '@/api'
import type { Share } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'

const router = useRouter()
const message = useMessage()
const shares = ref<Share[]>([])
const loading = ref(false)
const showCreate = ref(false)
const publicDomain = ref('')

const newShare = ref({ id: '', title: '', description: '' })

async function load() {
  loading.value = true
  try {
    shares.value = await listShares()
    const d = await getPublicDomain()
    publicDomain.value = d.public_domain
  } catch {}
  loading.value = false
}

async function handleCreate() {
  try {
    await createShare(newShare.value)
    message.success('Share created')
    showCreate.value = false
    newShare.value = { id: '', title: '', description: '' }
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to create')
  }
}

async function handleDelete(id: string) {
  try {
    await deleteShare(id)
    message.success('Deleted')
    load()
  } catch {
    message.error('Failed to delete')
  }
}

async function handleSaveDomain() {
  try {
    await setPublicDomain(publicDomain.value)
    message.success('Public domain updated')
    if (publicDomain.value) {
      message.info(`Caddy proxy configured: ${publicDomain.value} -> :3637`)
    }
  } catch {
    message.error('Failed to update domain')
  }
}

const shareColumns = [
  { title: 'ID', key: 'id' },
  { title: 'Folder', key: 'folder' },
  {
    title: 'Type', key: 'type', width: 80,
    render: (row: Share) => h('span', { style: { fontSize: '13px' } }, row.type || 'album'),
  },
  {
    title: 'Title', key: 'title',
    render: (row: Share) => row.title || row.id,
  },
  {
    title: 'Link', key: 'id',
    render: (row: Share) => {
      const href = publicDomain.value
        ? `https://${publicDomain.value}/p/${row.id}`
        : `/p/${row.id}`
      return h('a', { href, target: '_blank', style: { color: '#c96442', fontSize: '13px', textDecoration: 'none' } }, href)
    },
  },
  {
    title: '', key: 'actions', width: 100,
    render: (row: Share) => [
      h(NButton, { size: 'tiny', quaternary: true, onClick: () => router.push(`/files/public/${row.id}`) }, { default: () => 'Files' }),
      h(NPopconfirm, {
        onPositiveClick: () => handleDelete(row.id),
      }, {
        default: () => `Delete ${row.id}?`,
        trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error', style: 'margin-left:4px' }, { default: () => 'Delete' }),
      }),
    ],
  },
]

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Public">
      <template #actions>
        <n-button type="primary" @click="showCreate = true">New Album</n-button>
      </template>
    </PageHeader>

    <div class="mb-4 p-4 rounded-lg" style="background:#1a1a18;border:1px solid #2a2a28;">
      <div class="flex items-center gap-3">
        <span style="font-size:13px;color:#9a9990;white-space:nowrap;">Public Domain:</span>
        <n-input v-model:value="publicDomain" placeholder="files.example.com" style="max-width:300px;" />
        <n-button quaternary size="tiny" @click="handleSaveDomain">Save</n-button>
      </div>
      <div v-if="publicDomain" style="font-size:12px;color:#73726c;margin-top:6px;">
        Caddy proxies {{ publicDomain }} → :3637 &middot; Files served at /p/{id}
      </div>
    </div>

    <n-data-table
      :loading="loading"
      :columns="shareColumns"
      :data="shares"
      :pagination="{ pageSize: 20 }"
      :single-line="false"
      :scroll-x="700"
    />

    <!-- Create share modal -->
    <n-modal v-model:show="showCreate" title="New Album" preset="card" style="width:440px;max-width:90vw;" :mask-closable="false">
      <n-form>
        <n-form-item label="ID (leave empty for random)">
          <n-input v-model:value="newShare.id" placeholder="my-album" />
        </n-form-item>
        <n-form-item label="Title">
          <n-input v-model:value="newShare.title" placeholder="My Photo Album" />
        </n-form-item>
        <n-form-item label="Description">
          <n-input v-model:value="newShare.description" type="textarea" rows="2" placeholder="Optional description" />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <n-button quaternary @click="showCreate = false">Cancel</n-button>
          <n-button type="primary" @click="handleCreate">Create</n-button>
        </div>
      </template>
    </n-modal>

  </div>
</template>
