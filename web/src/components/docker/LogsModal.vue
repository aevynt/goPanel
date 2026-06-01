<script setup lang="ts">
import { ref, watch } from 'vue'
import { useMessage, NModal, NSpin } from 'naive-ui'
import { getContainerLogs } from '@/api'

const props = defineProps<{
  show: boolean
  containerId: string
  containerName: string
}>()

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void
}>()

const message = useMessage()
const loading = ref(false)
const logs = ref('')

async function loadLogs() {
  if (!props.containerId) return
  loading.value = true
  logs.value = ''
  try {
    logs.value = await getContainerLogs(props.containerId)
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to load logs')
  }
  loading.value = false
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    loadLogs()
  }
})
</script>

<template>
  <n-modal
    :show="props.show"
    @update:show="(val) => emit('update:show', val)"
    preset="card"
    style="width: 80%; max-width: 800px; background: #141413;"
    :title="'Container Logs: ' + props.containerName"
  >
    <n-spin :show="loading">
      <div class="log-container">
        <pre v-if="logs">{{ logs }}</pre>
        <div v-else class="empty-logs">No logs found</div>
      </div>
    </n-spin>
  </n-modal>
</template>

<style scoped>
.log-container {
  max-height: 500px;
  overflow-y: auto;
  background: #0d0d0e;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #2a2a28;
}

.log-container pre {
  margin: 0;
  font-family: 'Fira Code', 'Courier New', Courier, monospace;
  font-size: 12px;
  color: #e8e6dc;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.empty-logs {
  color: #73726c;
  text-align: center;
  padding: 24px;
}
</style>
