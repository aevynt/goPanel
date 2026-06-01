<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage, NTag, NSpace, NButton, NCard, NGrid, NGi, NSpin, NText, NIcon } from 'naive-ui'
import { CloudDownloadOutline, AlertCircleOutline } from '@vicons/ionicons5'
import { listApps, deployApp } from '@/api'
import type { AppCatalogItem } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'

const message = useMessage()
const apps = ref<AppCatalogItem[]>([])
const loading = ref(false)
const deployingMap = ref<Record<string, boolean>>({})
const selectedCategory = ref('All')

const categories = computed(() => {
  const cats = new Set<string>()
  cats.add('All')
  apps.value.forEach(app => cats.add(app.category))
  return Array.from(cats)
})

const filteredApps = computed(() => {
  if (selectedCategory.value === 'All') {
    return apps.value
  }
  return apps.value.filter(app => app.category === selectedCategory.value)
})

async function load() {
  loading.value = true
  try {
    apps.value = await listApps()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to load App Catalog')
  } finally {
    loading.value = false
  }
}

async function handleInstall(app: AppCatalogItem) {
  deployingMap.value[app.key] = true
  try {
    const result = await deployApp(app.key)
    if (result.warning) {
      message.warning(`App deployed on port ${result.port}, but proxy config had issues: ${result.warning}`, { duration: 10000 })
    } else {
      message.success(`Successfully deployed ${app.name} on port ${result.port}!`, { duration: 8000 })
    }
  } catch (err: any) {
    message.error(err.response?.data?.error || `Failed to deploy ${app.name}`)
  } finally {
    deployingMap.value[app.key] = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="App Store" />

    <n-space class="mb-6" :size="8">
      <n-button
        v-for="cat in categories"
        :key="cat"
        size="small"
        :secondary="selectedCategory !== cat"
        :type="selectedCategory === cat ? 'primary' : 'default'"
        @click="selectedCategory = cat"
      >
        {{ cat }}
      </n-button>
    </n-space>

    <n-spin :show="loading">
      <div v-if="filteredApps.length === 0 && !loading" class="empty-state">
        <n-icon :size="48"><alert-circle-outline /></n-icon>
        <p>No apps available in this category</p>
      </div>

      <n-grid cols="1 600:2 900:3" :x-gap="16" :y-gap="16" v-else>
        <n-gi v-for="app in filteredApps" :key="app.key">
          <n-card class="glass-card app-card" size="small">
            <template #header>
              <n-space align="center" justify="space-between" style="width: 100%;">
                <n-text strong class="app-title">{{ app.name }}</n-text>
                <n-tag size="small" type="info" :bordered="false">{{ app.category }}</n-tag>
              </n-space>
            </template>

            <div class="app-body">
              <p class="app-desc">{{ app.description }}</p>
              <div class="app-meta">
                <div class="meta-item">
                  <span class="meta-label">Docker Image:</span>
                  <span class="meta-val code">{{ app.image }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">Default Port:</span>
                  <span class="meta-val">{{ app.default_port }}</span>
                </div>
              </div>
            </div>

            <template #action>
              <n-button
                type="primary"
                block
                :loading="deployingMap[app.key]"
                @click="handleInstall(app)"
              >
                <template #icon>
                  <n-icon><cloud-download-outline /></n-icon>
                </template>
                1-Click Deploy
              </n-button>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<style scoped>
.glass-card {
  background: rgba(31, 31, 29, 0.6) !important;
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.app-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.app-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.app-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--claude-text-primary, #e8e6dc);
}

.app-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px 0;
}

.app-desc {
  font-size: 13px;
  color: var(--claude-text-secondary, #9a9990);
  line-height: 1.4;
  margin: 0;
  min-height: 40px;
}

.app-meta {
  background: rgba(0, 0, 0, 0.2);
  padding: 8px 12px;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
}

.meta-label {
  color: var(--claude-text-tertiary, #73726c);
}

.meta-val {
  color: var(--claude-text-primary, #e8e6dc);
  font-weight: 500;
}

.meta-val.code {
  font-family: monospace;
  font-size: 10px;
  background: rgba(255, 255, 255, 0.05);
  padding: 1px 4px;
  border-radius: 3px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  color: var(--claude-text-tertiary, #73726c);
  text-align: center;
  gap: 12px;
}
</style>
