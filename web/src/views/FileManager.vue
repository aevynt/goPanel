<script setup lang="ts">
import { h, ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, NButton, NIcon, NSpace, NBreadcrumb, NBreadcrumbItem, NDataTable, NModal, NInput, NProgress } from 'naive-ui'
import { FolderOpenOutline, DocumentOutline, TrashOutline, CloudUploadOutline, ArrowBackOutline, AddOutline, ArchiveOutline } from '@vicons/ionicons5'
import { listFiles, mkdir, uploadFile, removeFile, getFileBlob, zipFile, unzipFile } from '@/api'
import type { FileInfo } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'
import ModalFooter from '@/components/ModalFooter.vue'
import UploadProgressCard from '@/components/UploadProgressCard.vue'

const route = useRoute()
const router = useRouter()

function normalize(p: string): string {
  return '/' + p.replace(/^\/+|\/+$/g, '').split('/').filter(Boolean).join('/')
}

function pathFromRoute(): string {
  const p = route.params.path
  if (!p) return '/'
  const raw = Array.isArray(p) ? p.join('/') : p
  return normalize(raw)
}

interface UploadState {
  file: File
  path: string
  progress: number
  speed: string
  loaded: string
  total: string
  startTime: number
  lastLoaded: number
  lastTime: number
  status: 'uploading' | 'done' | 'error'
  error?: string
}

const message = useMessage()
const currentPath = ref(pathFromRoute())
const files = ref<FileInfo[]>([])
const loading = ref(false)
const breadcrumbs = ref<string[]>([])
const showNewFolderModal = ref(false)
const newFolderName = ref('')
const showEditor = ref(false)
const editorPath = ref('')
const editorContent = ref('')
const editorLoading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const uploads = ref<UploadState[]>([])

const showViewer = ref(false)
const viewerSrc = ref('')
const viewerType = ref<'image' | 'video' | ''>('')

const imageExts = new Set(['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg', '.bmp', '.ico'])
const videoExts = new Set(['.mp4', '.webm', '.ogg', '.avi', '.mov', '.mkv', '.wmv'])
const textExts = new Set(['.txt', '.json', '.xml', '.yaml', '.yml', '.toml', '.ini', '.cfg', '.conf', '.log', '.md', '.css', '.js', '.ts', '.vue', '.html', '.htm', '.php', '.py', '.rb', '.go', '.rs', '.java', '.c', '.cpp', '.h', '.hpp', '.sh', '.bat', '.ps1', '.sql', '.env', '.gitignore'])

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

function triggerUpload() {
  fileInput.value?.click()
}

async function load() {
  loading.value = true
  try {
    files.value = await listFiles(currentPath.value)
    breadcrumbs.value = currentPath.value === '/' ? [] : currentPath.value.split('/').filter(Boolean)
  } catch (err: any) {
    message.error('Failed to list files')
  }
  loading.value = false
  syncURL()
}

function syncURL() {
  const clean = currentPath.value === '/' ? '' : currentPath.value.replace(/^\//, '')
  const target = clean ? '/files/' + clean : '/files'
  router.replace(target)
}

function navigateTo(path: string) {
  currentPath.value = normalize(path)
  load()
}

function openDir(file: FileInfo) {
  if (file.is_dir) navigateTo(file.path)
}

function fileType(name: string): 'image' | 'video' | 'text' | 'other' {
  const ext = name.substring(name.lastIndexOf('.')).toLowerCase()
  if (imageExts.has(ext)) return 'image'
  if (videoExts.has(ext)) return 'video'
  if (textExts.has(ext)) return 'text'
  return 'other'
}

async function openFile(file: FileInfo) {
  if (file.is_dir) return
  const ft = fileType(file.name)
  if (ft === 'image' || ft === 'video') {
    try {
      const blob = await getFileBlob(file.path)
      viewerSrc.value = URL.createObjectURL(blob)
      viewerType.value = ft
      showViewer.value = true
    } catch {
      message.error('Failed to load file')
    }
    return
  }
  window.open(`/editor?path=${encodeURIComponent(file.path)}`, '_blank')
}

watch(showViewer, (v) => {
  if (!v && viewerSrc.value) {
    URL.revokeObjectURL(viewerSrc.value)
    viewerSrc.value = ''
  }
})

async function handleNewFolder() {
  if (!newFolderName.value) return
  const path = currentPath.value === '/'
    ? `/${newFolderName.value}`
    : `${currentPath.value}/${newFolderName.value}`
  try {
    await mkdir(path)
    message.success('Folder created')
    showNewFolderModal.value = false
    newFolderName.value = ''
    load()
  } catch {
    message.error('Failed to create folder')
  }
}

async function handleUpload(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return

  const pending: UploadState[] = []
  for (const file of input.files) {
    const path = currentPath.value === '/'
      ? `/${file.name}`
      : `${currentPath.value}/${file.name}`
    pending.push({
      file,
      path,
      progress: 0,
      speed: '0 B/s',
      loaded: '0 B',
      total: formatSize(file.size),
      startTime: Date.now(),
      lastLoaded: 0,
      lastTime: Date.now(),
      status: 'uploading',
    })
  }
  uploads.value = pending

  for (const state of pending) {
    try {
      await uploadFile(state.path, state.file, (pct) => {
        state.progress = pct
        const now = Date.now()
        const loaded = (state.file.size * pct) / 100
        const deltaLoaded = loaded - state.lastLoaded
        const deltaTime = now - state.lastTime
        if (deltaTime > 0) {
          state.speed = formatSpeed((deltaLoaded / deltaTime) * 1000)
        }
        state.loaded = formatSize(loaded)
        state.lastLoaded = loaded
        state.lastTime = now
      })
      state.status = 'done'
    } catch (err: any) {
      state.status = 'error'
      state.error = err.response?.data?.error || 'Upload failed'
    }
  }

  const hasError = pending.some((u) => u.status === 'error')
  if (hasError) {
    message.error('Some files failed to upload')
  } else {
    message.success('All files uploaded')
  }

  uploads.value = []
  input.value = ''
  load()
}

async function handleDelete(file: FileInfo) {
  try {
    await removeFile(file.path)
    message.success('Deleted')
    load()
  } catch {
    message.error('Failed to delete')
  }
}

async function handleZip(file: FileInfo) {
  const defaultZipPath = file.path + '.zip'
  loading.value = true
  try {
    await zipFile(file.path, defaultZipPath)
    message.success(`Zipped successfully to ${defaultZipPath}`)
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to zip')
  } finally {
    loading.value = false
  }
}

async function handleUnzip(file: FileInfo) {
  const parentPath = file.path.substring(0, file.path.lastIndexOf('/')) || '/'
  loading.value = true
  try {
    await unzipFile(file.path, parentPath)
    message.success(`Extracted successfully to ${parentPath}`)
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to extract')
  } finally {
    loading.value = false
  }
}

const columns = [
  {
    title: 'Name',
    key: 'name',
    render: (row: FileInfo) => {
      return h('a', {
        href: '#',
        onClick: (e: Event) => { e.preventDefault(); row.is_dir ? openDir(row) : openFile(row) },
        style: { cursor: 'pointer', fontWeight: row.is_dir ? '500' : 'normal' },
      }, [
        h(NIcon, { size: 16, style: { marginRight: '6px', verticalAlign: 'middle', color: row.is_dir ? 'var(--claude-accent)' : 'var(--claude-text-tertiary)' } }, {
          default: () => h(row.is_dir ? FolderOpenOutline : DocumentOutline),
        }),
        row.name,
      ])
    },
  },
  {
    title: 'Size',
    key: 'size',
    width: 100,
    render: (row: FileInfo) => row.is_dir ? '-' : formatSize(row.size),
  },
  { title: 'Mode', key: 'mode', width: 100 },
  { title: 'Modified', key: 'mod_time', width: 180 },
  {
    title: 'Actions',
    key: 'actions',
    width: 140,
    render: (row: FileInfo) => {
      const btns: any[] = []

      // Zip action
      btns.push(h(NButton, {
        size: 'tiny',
        quaternary: true,
        type: 'info',
        onClick: () => handleZip(row)
      }, {
        default: () => h(NIcon, { size: 16 }, { default: () => h(ArchiveOutline) }),
      }))

      // Unzip action if .zip file
      if (!row.is_dir && row.name.toLowerCase().endsWith('.zip')) {
        btns.push(h(NButton, {
          size: 'tiny',
          quaternary: true,
          type: 'warning',
          onClick: () => handleUnzip(row)
        }, {
          default: () => h(NIcon, { size: 16 }, { default: () => h(FolderOpenOutline) }),
        }))
      }

      // Delete action
      btns.push(h(NButton, {
        size: 'tiny',
        quaternary: true,
        type: 'error',
        onClick: () => handleDelete(row)
      }, {
        default: () => h(NIcon, { size: 16 }, { default: () => h(TrashOutline) }),
      }))

      return h(NSpace, { size: 'small' }, { default: () => btns })
    },
  },
]

watch(() => route.params.path, () => {
  const newPath = pathFromRoute()
  if (newPath !== currentPath.value) {
    currentPath.value = newPath
    load()
  }
})
onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="File Manager">
      <template #actions>
        <input ref="fileInput" type="file" multiple class="hidden-input" @change="handleUpload" />
        <n-button quaternary :loading="uploads.length > 0" @click="triggerUpload">
          <template #icon>
            <n-icon :size="18"><cloud-upload-outline /></n-icon>
          </template>
          Upload Files
        </n-button>
        <n-button type="primary" @click="showNewFolderModal = true">
          <template #icon>
            <n-icon :size="18"><add-outline /></n-icon>
          </template>
          New Folder
        </n-button>
      </template>
    </PageHeader>

    <n-breadcrumb class="mb-4">
      <n-breadcrumb-item>
        <a href="#" @click.prevent="navigateTo('/')">root</a>
      </n-breadcrumb-item>
      <n-breadcrumb-item v-for="(crumb, i) in breadcrumbs" :key="i">
        <a href="#" @click.prevent="navigateTo('/' + breadcrumbs.slice(0, i + 1).join('/'))">
          {{ crumb }}
        </a>
      </n-breadcrumb-item>
    </n-breadcrumb>

    <n-button
      v-if="currentPath !== '/'"
      quaternary size="tiny"
      @click="navigateTo(currentPath.split('/').slice(0, -1).join('/') || '/')"
      class="mb-2"
    >
      <template #icon>
        <n-icon :size="14"><arrow-back-outline /></n-icon>
      </template>
      Back
    </n-button>

    <div v-if="uploads.length > 0" class="flex flex-col gap-2 mb-4">
      <UploadProgressCard
        v-for="(u, i) in uploads"
        :key="i"
        :file-name="u.file.name"
        :file-total="u.total"
        :progress="u.progress"
        :speed="u.speed"
        :status="u.status"
        :error="u.error"
      >
        <template #icon>
          <n-icon :size="20"><document-outline /></n-icon>
        </template>
      </UploadProgressCard>
    </div>

    <div class="transition-opacity duration-200" :class="{ 'opacity-50 pointer-events-none': uploads.length > 0 }">
      <n-data-table
        :loading="loading"
        :columns="columns"
        :data="files"
        :single-line="false"
        :scroll-x="600"
      />
    </div>

    <n-modal v-model:show="showNewFolderModal" title="New Folder" preset="card" style="width: 400px; max-width: 90vw;" :mask-closable="false">
      <n-input v-model:value="newFolderName" placeholder="folder name" @keyup.enter="handleNewFolder" />
      <template #footer>
        <ModalFooter submit-label="Create" @cancel="showNewFolderModal = false" @submit="handleNewFolder" />
      </template>
    </n-modal>

    <n-modal v-model:show="showViewer" preset="card" style="width:90vw;max-width:960px;background:transparent;box-shadow:none;" :mask-closable="true" @mask-click="showViewer = false">
      <div style="display:flex;align-items:center;justify-content:center;min-height:60vh;padding:16px;">
        <img v-if="viewerType === 'image'" :src="viewerSrc" style="max-width:100%;max-height:80vh;border-radius:8px;object-fit:contain;" @click="showViewer = false" />
        <video v-if="viewerType === 'video'" :src="viewerSrc" controls style="max-width:100%;max-height:80vh;border-radius:8px;" @click="showViewer = false" />
      </div>
    </n-modal>
  </div>
</template>

<style scoped>
.hidden-input { display: none; }
</style>
