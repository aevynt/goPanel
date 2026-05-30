<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, NButton } from 'naive-ui'
import { readFile, writeFile } from '@/api'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const path = ref('')
const content = ref('')
const loading = ref(false)
const saving = ref(false)
const saved = ref(false)

async function load() {
  const p = route.query.path as string
  if (!p) {
    message.error('No path specified')
    return
  }
  path.value = p
  loading.value = true
  try {
    content.value = await readFile(p)
    document.title = `Editing: ${p.split('/').pop()}`
  } catch {
    message.error('Failed to read file')
  }
  loading.value = false
}

async function save() {
  saving.value = true
  try {
    await writeFile(path.value, content.value)
    saved.value = true
    message.success('File saved')
  } catch {
    message.error('Failed to save file')
  }
  saving.value = false
}

function close() {
  if (window.history.length > 1) {
    router.back()
  } else {
    window.close()
  }
}

onMounted(load)
</script>

<template>
  <div class="editor-page">
    <div class="editor-topbar">
      <div class="editor-path">{{ path || 'No file selected' }}</div>
      <div class="editor-actions">
        <n-button v-if="saved" quaternary size="tiny" style="color:#16a34a;" disabled>Saved</n-button>
        <n-button quaternary size="tiny" :loading="saving" @click="save">Save</n-button>
        <n-button quaternary size="tiny" @click="close">Close</n-button>
      </div>
    </div>
    <textarea
      v-model="content"
      class="editor-textarea"
      :disabled="loading"
      placeholder="Loading..."
      spellcheck="false"
    ></textarea>
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body { height: 100%; background: #141413; color: #e8e6dc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
.editor-page { display: flex; flex-direction: column; height: 100vh; }
.editor-topbar { display: flex; align-items: center; padding: 8px 16px; background: #1a1a18; border-bottom: 1px solid #2a2a28; gap: 12px; }
.editor-path { flex: 1; font-size: 13px; color: #9a9990; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.editor-actions { display: flex; gap: 8px; align-items: center; }
.editor-textarea { flex: 1; width: 100%; padding: 16px; background: #141413; color: #e8e6dc; border: none; outline: none; resize: none; font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace; font-size: 13px; line-height: 1.6; tab-size: 2; }
.editor-textarea:disabled { opacity: 0.5; }
.editor-textarea::placeholder { color: #5e5d59; }
</style>
