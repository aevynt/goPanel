<script setup lang="ts">
import { ref, h, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { MenuOption } from 'naive-ui'
import { NIcon, NButton, NTag } from 'naive-ui'
import {
  SpeedometerOutline,
  ServerOutline,
  DocumentTextOutline,
  FolderOpenOutline,
  GitNetworkOutline,
  GlobeOutline,
  AlbumsOutline,
  PeopleOutline,
  SettingsOutline,
  LogOutOutline,
  CubeOutline,
  AppsOutline,
} from '@vicons/ionicons5'
import type { UpdateInfo } from '@/api/client'
import { checkUpdate } from '@/api'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const isMobile = ref(false)
const updateInfo = ref<UpdateInfo | null>(null)

onMounted(async () => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  try {
    updateInfo.value = await checkUpdate()
  } catch {}
})

function checkMobile() {
  isMobile.value = window.innerWidth < 768
}

onUnmounted(() => window.removeEventListener('resize', checkMobile))

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions: MenuOption[] = [
  { label: 'Dashboard', key: 'dashboard', icon: renderIcon(SpeedometerOutline) },
  { label: 'Services', key: 'services', icon: renderIcon(ServerOutline) },
  { label: 'Binaries', key: 'binaries', icon: renderIcon(DocumentTextOutline) },
  { label: 'Files', key: 'files', icon: renderIcon(FolderOpenOutline) },
  { label: 'Ports', key: 'ports', icon: renderIcon(GitNetworkOutline) },
  { label: 'Sites', key: 'sites', icon: renderIcon(GlobeOutline) },
  { label: 'Public', key: 'public', icon: renderIcon(AlbumsOutline) },
  ...(auth.user?.role === 'admin'
    ? [
        { label: 'Docker', key: 'docker', icon: renderIcon(CubeOutline) },
        { label: 'App Store', key: 'apps', icon: renderIcon(AppsOutline) },
        { label: 'Users', key: 'users', icon: renderIcon(PeopleOutline) }
      ]
    : []),
  { label: 'Settings', key: 'settings', icon: renderIcon(SettingsOutline) },
]

function navigate(key: string) {
  router.push(`/${key}`)
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <n-layout class="claude-layout" has-sider position="absolute">
    <!-- Desktop sidebar -->
    <n-layout-sider
      v-if="!isMobile"
      bordered
      collapse-mode="width"
      :collapsed-width="60"
      :width="220"
      show-trigger="bar"
      trigger-style="top: 50%; background: transparent; border: none; color: #73726c;"
      :native-scrollbar="false"
    >
      <div class="sider-header">
        <span class="logo-text">goPanel</span>
      </div>

      <n-menu
        :options="menuOptions"
        :value="route.name as string"
        @update:value="(key: string) => router.push(`/${key}`)"
        :indent="18"
      />

      <div class="sider-footer">
        <div class="user-name">{{ auth.user?.username }}</div>
        <div class="user-role">{{ auth.user?.role }}</div>
        <div class="flex items-center gap-2" style="margin: 4px 0;">
          <span class="version-badge">v{{ updateInfo?.current_version || '-' }}</span>
          <span v-if="updateInfo?.has_update" class="update-dot" :title="updateInfo.latest_version + ' available'">↑</span>
        </div>
        <n-button quaternary size="tiny" class="logout-btn" @click="handleLogout">
          <template #icon>
            <n-icon :size="16"><log-out-outline /></n-icon>
          </template>
          Logout
        </n-button>
      </div>
    </n-layout-sider>

    <n-layout>
      <n-layout-header class="topbar">
        <span class="current-page">{{ route.name ? String(route.name).charAt(0).toUpperCase() + String(route.name).slice(1) : '' }}</span>
        <div class="topbar-spacer" />
        <n-button quaternary size="tiny" class="logout-btn-top" @click="handleLogout">
          <template #icon>
            <n-icon :size="16"><log-out-outline /></n-icon>
          </template>
        </n-button>
      </n-layout-header>
      <n-layout-content class="layout-content" :native-scrollbar="false">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </n-layout-content>
    </n-layout>
  </n-layout>

  <!-- Mobile bottom navigation -->
  <div v-if="isMobile" class="footbar">
    <button
      v-for="item in menuOptions"
      :key="item.key"
      class="footbar-item"
      :class="{ active: route.name === item.key }"
      @click="navigate(item.key as string)"
    >
      <component :is="item.icon" />
      <span class="footbar-label">{{ item.label }}</span>
    </button>
  </div>
</template>

<style scoped>
.claude-layout {
  background: var(--claude-bg, #141413);
}

.sider-header {
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-family: var(--font-heading, 'Playfair Display', serif);
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--claude-text-primary, #e8e6dc);
  letter-spacing: -0.02em;
}

.sider-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.sider-footer .user-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--claude-text-primary, #e8e6dc);
}

.sider-footer .user-role {
  font-size: 11px;
  color: var(--claude-text-tertiary, #73726c);
  text-transform: capitalize;
}

.logout-btn {
  width: 100%;
  justify-content: flex-start;
  padding: 4px 8px;
  margin-top: 4px;
}

.version-badge {
  font-size: 10px;
  font-weight: 500;
  padding: 1px 6px;
  border-radius: 4px;
  background: rgba(255,255,255,0.06);
  color: var(--claude-text-tertiary, #73726c);
  letter-spacing: 0.02em;
}

.update-dot {
  font-size: 12px;
  font-weight: 700;
  color: var(--claude-accent, #c96442);
  animation: pulse-dot 2s ease-in-out infinite;
  cursor: help;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.topbar {
  height: 48px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  background: var(--claude-bg, #141413);
  border-bottom: 1px solid var(--claude-border-light, #2a2a28);
  gap: 8px;
}

.topbar-spacer {
  flex: 1;
}

.logout-btn-top {
  padding: 4px !important;
}

.current-page {
  font-family: var(--font-heading, 'Playfair Display', serif);
  font-size: 1rem;
  font-weight: 500;
  color: var(--claude-text-primary, #e8e6dc);
}

.layout-content {
  padding: 24px;
  background: var(--claude-bg, #141413);
  min-height: calc(100vh - 48px);
}

@media (max-width: 767px) {
  .layout-content {
    padding: 16px;
    padding-bottom: 72px;
    min-height: calc(100vh - 48px - 56px);
  }
}

/* ---- Mobile bottom nav ---- */

.footbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: var(--claude-surface, #1f1f1d);
  border-top: 1px solid var(--claude-border-light, #2a2a28);
  display: flex;
  align-items: stretch;
  z-index: 1000;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.footbar-item {
  flex: 1 0 0;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 4px 2px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: var(--claude-text-tertiary, #73726c);
  transition: color 0.15s ease;
  -webkit-tap-highlight-color: transparent;
}

.footbar-item.active {
  color: var(--claude-accent, #c96442);
}

.footbar-item :deep(.n-icon) {
  font-size: 20px;
}

.footbar-label {
  font-size: 10px;
  font-weight: 500;
  line-height: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
