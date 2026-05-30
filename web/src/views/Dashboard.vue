<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { NSpin, NGrid, NGi, NProgress, NButton } from 'naive-ui'
import { getDashboardStats, checkUpdate } from '@/api'
import type { DashboardStats, UpdateInfo } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'

const stats = ref<DashboardStats | null>(null)
const loading = ref(true)
const updateInfo = ref<UpdateInfo | null>(null)
let ws: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null

function connectWs() {
	const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
	ws = new WebSocket(`${protocol}//${location.host}/ws`)

	ws.onmessage = (event) => {
		try {
			stats.value = JSON.parse(event.data)
			loading.value = false
		} catch {}
	}

	ws.onclose = () => {
		reconnectTimer = setTimeout(connectWs, 3000)
	}
}

function formatBytes(bytes: number): string {
	if (bytes === 0) return '0 B'
	const k = 1024
	const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
	const i = Math.floor(Math.log(bytes) / Math.log(k))
	return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(async () => {
	try {
		stats.value = await getDashboardStats()
	} catch {}
	try {
		updateInfo.value = await checkUpdate()
	} catch {}
	loading.value = false
	connectWs()
})

onUnmounted(() => {
	if (ws) ws.close()
	if (reconnectTimer) clearTimeout(reconnectTimer)
})
</script>

<template>
  <div>
    <PageHeader title="Dashboard" />

    <div v-if="updateInfo?.has_update" class="update-banner">
      <div class="update-banner-content">
        <div class="update-icon">↑</div>
        <div>
          <div class="update-title">
            Update available: <strong>{{ updateInfo.latest_version }}</strong>
          </div>
          <div class="update-sub">
            {{ updateInfo.release?.name || '' }} &middot; Current: {{ updateInfo.current_version }}
          </div>
        </div>
      </div>
      <a :href="updateInfo.release?.html_url || '#'" target="_blank" class="update-btn">
        View Release
      </a>
    </div>

    <div v-if="updateInfo && !updateInfo.has_update && !updateInfo.error" class="update-banner update-banner-ok">
      <div class="update-banner-content">
        <div class="update-icon" style="color:var(--claude-accent);">✓</div>
        <div>
          <div class="update-title">goPanel is up to date</div>
          <div class="update-sub">{{ updateInfo.current_version }} &middot; checked {{ new Date(updateInfo.checked_at).toLocaleString() }}</div>
        </div>
      </div>
    </div>

    <n-spin :show="loading">
      <n-grid cols="1 600:2 900:4" :x-gap="16" :y-gap="16">
        <n-gi>
          <div class="card h-card">
            <div class="card-title">System</div>
            <div class="h-card-body">
              <div class="stat-row">
                <div class="stat-label">Uptime</div>
                <div class="stat-value">{{ stats?.uptime || '-' }}</div>
              </div>
              <div class="stat-row">
                <div class="stat-label">Hostname</div>
                <div class="stat-value">{{ stats?.hostname || '-' }}</div>
              </div>
              <div class="stat-row">
                <div class="stat-label">OS</div>
                <div class="stat-value">{{ stats?.os || '-' }}</div>
              </div>
              <div class="stat-row">
                <div class="stat-label">Kernel</div>
                <div class="stat-value">{{ stats?.kernel || '-' }}</div>
              </div>
              <div class="stat-row">
                <div class="stat-label">Go Version</div>
                <div class="stat-value">{{ stats?.go_version || '-' }}</div>
              </div>
              <div class="stat-row">
                <div class="stat-label">Panel Version</div>
                <div class="stat-value">{{ stats?.version || '-' }}</div>
              </div>
            </div>
          </div>
        </n-gi>
        <n-gi>
          <div class="card h-card">
            <div class="card-title">CPU</div>
            <div class="h-card-center">
              <n-progress type="circle" :percentage="Math.round(stats?.cpu_percent || 0)" :width="100">
                {{ Math.round(stats?.cpu_percent || 0) }}%
              </n-progress>
            </div>
          </div>
        </n-gi>
        <n-gi>
          <div class="card h-card">
            <div class="card-title">Memory</div>
            <div class="h-card-center">
              <n-progress type="circle" :percentage="Math.round(stats?.memory.used_percent || 0)" :width="100" status="warning">
                {{ Math.round(stats?.memory.used_percent || 0) }}%
              </n-progress>
              <div class="stat-detail">
                {{ formatBytes(stats?.memory.used || 0) }} / {{ formatBytes(stats?.memory.total || 0) }}
              </div>
            </div>
          </div>
        </n-gi>
        <n-gi>
          <div class="card h-card">
            <div class="card-title">Disk</div>
            <div class="h-card-center">
              <n-progress type="circle" :percentage="Math.round(stats?.disk.used_percent || 0)" :width="100" status="error">
                {{ Math.round(stats?.disk.used_percent || 0) }}%
              </n-progress>
              <div class="stat-detail">
                {{ formatBytes(stats?.disk.used || 0) }} / {{ formatBytes(stats?.disk.total || 0) }}
              </div>
            </div>
          </div>
        </n-gi>
      </n-grid>

      <n-grid cols="1 600:3" :x-gap="16" :y-gap="16" class="mt-4">
        <n-gi>
          <div class="card h-card">
            <div class="h-card-center gap-1">
              <div class="stat-label">Services</div>
              <div class="stat-value font-heading text-[1.75rem]">{{ stats?.services_count || 0 }}</div>
            </div>
          </div>
        </n-gi>
        <n-gi>
          <div class="card h-card">
            <div class="h-card-center gap-1">
              <div class="stat-label">Ports Open</div>
              <div class="stat-value font-heading text-[1.75rem]">{{ stats?.ports_open || 0 }}</div>
            </div>
          </div>
        </n-gi>
        <n-gi>
          <div class="card h-card">
            <div class="h-card-center gap-1">
              <div class="stat-label">Sites</div>
              <div class="stat-value font-heading text-[1.75rem]">{{ stats?.sites_count || 0 }}</div>
            </div>
          </div>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<style scoped>
.h-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.h-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.h-card-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 6px;
  padding: 8px 0;
}

.h-card-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px 0;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2px 0;
}

.stat-label {
  font-size: 13px;
  color: var(--claude-text-secondary, #9c9a92);
}

.stat-value {
  font-size: 14px;
  font-weight: 500;
  color: var(--claude-text-primary, #e8e6dc);
  text-align: right;
}

.stat-detail {
	text-align: center;
	font-size: 12px;
	color: var(--claude-text-tertiary, #73726c);
}

.update-banner {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 12px;
	margin-bottom: 20px;
	padding: 12px 16px;
	border-radius: 8px;
	background: linear-gradient(135deg, rgba(201,100,66,0.15), rgba(201,100,66,0.05));
	border: 1px solid rgba(201,100,66,0.3);
}

.update-banner-ok {
	background: rgba(80,180,100,0.08);
	border-color: rgba(80,180,100,0.2);
}

.update-banner-content {
	display: flex;
	align-items: center;
	gap: 12px;
}

.update-icon {
	font-size: 20px;
	font-weight: 700;
	width: 36px;
	height: 36px;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	background: rgba(201,100,66,0.2);
	color: var(--claude-accent, #c96442);
	flex-shrink: 0;
}

.update-title {
	font-size: 14px;
	font-weight: 500;
	color: var(--claude-text-primary, #e8e6dc);
}

.update-sub {
	font-size: 12px;
	color: var(--claude-text-tertiary, #73726c);
	margin-top: 2px;
}

.update-btn {
	padding: 6px 14px;
	border-radius: 6px;
	font-size: 13px;
	font-weight: 500;
	background: var(--claude-accent, #c96442);
	color: #fff;
	text-decoration: none;
	white-space: nowrap;
	transition: opacity 0.15s;
}

.update-btn:hover {
	opacity: 0.85;
}
</style>
