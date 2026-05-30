<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useMessage, NButton, NSpace, NTag, NPopconfirm, NUpload, NProgress, NIcon } from 'naive-ui'
import { CloudUploadOutline, DocumentOutline, TrashOutline } from '@vicons/ionicons5'
import { listBinaries, uploadBinary, deleteBinary } from '@/api'
import type { Binary } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'
import UploadProgressCard from '@/components/UploadProgressCard.vue'

interface UploadState {
  file: File
  progress: number
  speed: string
  loaded: string
  total: string
  startTime: number
  lastLoaded: number
  lastTime: number
}

const message = useMessage()
const binaries = ref<Binary[]>([])
const loading = ref(false)
const currentUpload = ref<UploadState | null>(null)

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatSpeed(bytesPerSec: number): string {
  if (bytesPerSec === 0) return '0 B/s'
  const k = 1024
  const sizes = ['B/s', 'KB/s', 'MB/s', 'GB/s']
  const i = Math.floor(Math.log(bytesPerSec) / Math.log(k))
  return parseFloat((bytesPerSec / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function handleUpload({ file }: any) {
  const f = file.file || file
  currentUpload.value = {
    file: f,
    progress: 0,
    speed: '0 B/s',
    loaded: '0 B',
    total: formatSize(f.size),
    startTime: Date.now(),
    lastLoaded: 0,
    lastTime: Date.now(),
  }

  uploadBinary(f, undefined, (pct) => {
    const u = currentUpload.value
    if (!u) return
    u.progress = pct
    const now = Date.now()
    const loaded = (f.size * pct) / 100
    const elapsed = (now - u.startTime) / 1000
    const deltaLoaded = loaded - u.lastLoaded
    const deltaTime = now - u.lastTime
    if (deltaTime > 0) {
      const instantSpeed = (deltaLoaded / deltaTime) * 1000
      u.speed = formatSpeed(instantSpeed)
    } else {
      u.speed = formatSpeed(elapsed > 0 ? loaded / elapsed : 0)
    }
    u.loaded = formatSize(loaded)
    u.lastLoaded = loaded
    u.lastTime = now
  })
    .then(() => {
      message.success('Binary uploaded')
      load()
    })
    .catch((err: any) => {
      message.error(err.response?.data?.error || 'Upload failed')
    })
    .finally(() => {
      currentUpload.value = null
    })
}

async function load() {
  loading.value = true
  try {
    binaries.value = await listBinaries()
  } catch (err: any) {
    message.error('Failed to load binaries')
  }
  loading.value = false
}

async function handleDelete(id: number) {
  try {
    await deleteBinary(id)
    message.success('Binary deleted')
    load()
  } catch (err: any) {
    message.error('Failed to delete binary')
  }
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: 'Name', key: 'name' },
  { title: 'Size', key: 'size', render: (row: Binary) => formatSize(row.size) },
  { title: 'Version', key: 'version' },
  { title: 'Path', key: 'path', ellipsis: true },
  {
    title: 'Created',
    key: 'created_at',
    render: (row: any) => row.created_at || '-',
  },
  {
    title: '',
    key: 'actions',
    width: 80,
    render: (row: Binary) => h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
      default: () => 'Delete this binary?',
      trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, {
        default: () => h(NIcon, { size: 16 }, { default: () => h(TrashOutline) }),
      }),
    }),
  },
]

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Binaries">
      <template #actions>
        <div class="ml-auto inline-flex items-center">
          <n-upload
            :show-file-list="false"
            :custom-request="({ file }) => handleUpload({ file: file.file })"
          >
            <n-button type="primary" :loading="!!currentUpload">
              <template #icon>
                <n-icon :size="18"><cloud-upload-outline /></n-icon>
              </template>
              Upload Binary
            </n-button>
          </n-upload>
        </div>
      </template>
    </PageHeader>

    <UploadProgressCard
      v-if="currentUpload"
      :file-name="currentUpload.file.name"
      :file-total="currentUpload.total"
      :progress="currentUpload.progress"
      :speed="currentUpload.speed"
    >
      <template #icon>
        <n-icon :size="22"><document-outline /></n-icon>
      </template>
    </UploadProgressCard>

    <div class="transition-opacity duration-200" :class="{ 'opacity-50 pointer-events-none': !!currentUpload }">
      <n-data-table
        :loading="loading"
        :columns="columns"
        :data="binaries"
        :pagination="{ pageSize: 20 }"
        :single-line="false"
        :scroll-x="700"
      />
    </div>
  </div>
</template>
