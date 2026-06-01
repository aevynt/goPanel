<script setup lang="ts">
import { ref } from 'vue'
import { useMessage, NInput, NCard, NButton, NSpace } from 'naive-ui'
import { deployCompose } from '@/api'

const emit = defineEmits<{
  (e: 'deployed'): void
}>()

const message = useMessage()
const projectName = ref('')
const composeContent = ref(`version: '3.8'
services:
  nginx-test:
    image: nginx:alpine
    ports:
      - "8080:80"
    restart: unless-stopped
`)
const deploying = ref(false)

async function handleDeploy() {
  if (!projectName.value.trim()) {
    message.warning('Project name is required')
    return
  }
  if (!composeContent.value.trim()) {
    message.warning('Compose content is required')
    return
  }
  deploying.value = true
  try {
    await deployCompose(projectName.value, composeContent.value)
    message.success('Deployment started successfully')
    projectName.value = ''
    emit('deployed')
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to deploy compose project')
  }
  deploying.value = false
}
</script>

<template>
  <n-card class="glass-card" title="New Docker Compose Project">
    <n-space vertical :size="16">
      <div>
        <div class="field-label">Project Name</div>
        <n-input
          v-model:value="projectName"
          placeholder="e.g. my-web-app"
          maxlength="64"
          show-count
        />
      </div>
      <div>
        <div class="field-label">docker-compose.yml</div>
        <n-input
          v-model:value="composeContent"
          type="textarea"
          :autosize="{ minRows: 10, maxRows: 25 }"
          placeholder="Paste docker-compose.yml here..."
          class="code-editor"
        />
      </div>
      <n-button type="primary" :loading="deploying" @click="handleDeploy">
        Deploy Project
      </n-button>
    </n-space>
  </n-card>
</template>

<style scoped>
.field-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--claude-text-secondary, #9a9990);
  margin-bottom: 6px;
}

.code-editor :deep(textarea) {
  font-family: 'Fira Code', 'Courier New', Courier, monospace;
  font-size: 13px;
  background: #0d0d0e !important;
  color: #a7c080 !important;
  line-height: 1.5;
}

.glass-card {
  background: rgba(31, 31, 29, 0.6) !important;
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
</style>
