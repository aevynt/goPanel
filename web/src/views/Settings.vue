<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage, NIcon } from 'naive-ui'
import { PersonOutline, ServerOutline, LockClosedOutline, NotificationsOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'
import { getSettings, updateSettings } from '@/api'
import type { PanelSettings } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'

// Import Modular Tab Components
import ProfileTab from './settings/ProfileTab.vue'
import PanelTab from './settings/PanelTab.vue'
import AlertsTab from './settings/AlertsTab.vue'
import SecurityTab from './settings/SecurityTab.vue'

const auth = useAuthStore()
const message = useMessage()

const settings = ref<PanelSettings>({
  panel_domain: '',
  port: 3636,
  log_level: 'info',
  public_domain: '',
  public_port: 3637,
  discord_webhook: '',
  telegram_token: '',
  telegram_chat_id: '',
  alert_temp_threshold: 80,
  alert_cpu_threshold: 90,
  alert_ram_threshold: 90,
  alert_disk_threshold: 90
})
const saving = ref(false)
const loading = ref(false)
const activeSection = ref('profile')

interface NavItem {
  key: string
  label: string
  icon: any
  adminOnly?: boolean
}

const navItems: NavItem[] = [
  { key: 'profile', label: 'Profile', icon: PersonOutline },
  { key: 'panel', label: 'Panel', icon: ServerOutline, adminOnly: true },
  { key: 'alerts', label: 'Alerts', icon: NotificationsOutline, adminOnly: true },
  { key: 'security', label: 'Security', icon: LockClosedOutline },
]

const visibleNavItems = computed(() => navItems.filter((item) => !item.adminOnly || auth.user?.role === 'admin'))

async function loadSettings() {
  loading.value = true
  try {
    settings.value = await getSettings()
  } catch {
    message.error('Failed to load settings')
  }
  loading.value = false
}

async function saveDomain() {
  saving.value = true
  try {
    settings.value = await updateSettings({
      panel_domain: settings.value.panel_domain,
      public_domain: settings.value.public_domain,
      discord_webhook: settings.value.discord_webhook,
      telegram_token: settings.value.telegram_token,
      telegram_chat_id: settings.value.telegram_chat_id,
      alert_temp_threshold: settings.value.alert_temp_threshold,
      alert_cpu_threshold: settings.value.alert_cpu_threshold,
      alert_ram_threshold: settings.value.alert_ram_threshold,
      alert_disk_threshold: settings.value.alert_disk_threshold,
    })
    message.success('Domain updated successfully')
  } catch {
    message.error('Failed to save settings')
  }
  saving.value = false
}

async function saveAlerts() {
  saving.value = true
  try {
    settings.value = await updateSettings({
      panel_domain: settings.value.panel_domain,
      public_domain: settings.value.public_domain,
      discord_webhook: settings.value.discord_webhook,
      telegram_token: settings.value.telegram_token,
      telegram_chat_id: settings.value.telegram_chat_id,
      alert_temp_threshold: settings.value.alert_temp_threshold,
      alert_cpu_threshold: settings.value.alert_cpu_threshold,
      alert_ram_threshold: settings.value.alert_ram_threshold,
      alert_disk_threshold: settings.value.alert_disk_threshold,
    })
    message.success('Alerts settings saved successfully')
  } catch {
    message.error('Failed to save alerts settings')
  }
  saving.value = false
}

onMounted(() => {
  if (auth.user?.role === 'admin') loadSettings()
  auth.checkAuth()
})
</script>

<template>
  <div>
    <PageHeader title="Settings" />

    <div class="flex flex-col md:flex-row gap-0">
      <!-- Sidebar Navigation -->
      <div class="w-full md:w-[200px] shrink-0 flex flex-row md:flex-col gap-[2px] pb-3 md:pb-0 pr-0 md:pr-3 border-b md:border-b-0 md:border-r border-border-light overflow-x-auto">
        <div
          v-for="item in visibleNavItems"
          :key="item.key"
          class="flex items-center gap-[10px] px-3 py-[10px] rounded-md cursor-pointer text-sm font-medium select-none transition-all duration-150 whitespace-nowrap md:whitespace-normal"
          :class="activeSection === item.key
            ? 'text-accent bg-accent-muted'
            : 'text-text-secondary hover:text-text-primary hover:bg-border-light'"
          @click="activeSection = item.key"
        >
          <n-icon :size="18"><component :is="item.icon" /></n-icon>
          <span>{{ item.label }}</span>
        </div>
      </div>

      <!-- Active Content Section -->
      <div class="flex-1 min-w-0 pl-0 md:pl-7 pt-5 md:pt-0">
        <!-- Render Tabs Dynamically -->
        <ProfileTab v-if="activeSection === 'profile'" />

        <PanelTab
          v-else-if="activeSection === 'panel'"
          :settings="settings"
          :loading="loading"
          :saving="saving"
          @save-domain="saveDomain"
        />

        <AlertsTab
          v-else-if="activeSection === 'alerts'"
          :settings="settings"
          :loading="loading"
          :saving="saving"
          @save-alerts="saveAlerts"
        />

        <SecurityTab v-else-if="activeSection === 'security'" />
      </div>
    </div>
  </div>
</template>
