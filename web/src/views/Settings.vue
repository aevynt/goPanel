<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage, NIcon } from 'naive-ui'
import { PersonOutline, ServerOutline, LockClosedOutline, CheckmarkCircleOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'
import { getSettings, updateSettings } from '@/api'
import type { PanelSettings } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'

const auth = useAuthStore()
const message = useMessage()

const settings = ref<PanelSettings>({ panel_domain: '', port: 8080, log_level: 'info', public_domain: '', public_port: 3637 })
const saving = ref(false)
const loading = ref(false)
const activeSection = ref('profile')

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

interface NavItem {
  key: string
  label: string
  icon: any
  adminOnly?: boolean
}

const navItems: NavItem[] = [
  { key: 'profile', label: 'Profile', icon: PersonOutline },
  { key: 'panel', label: 'Panel', icon: ServerOutline, adminOnly: true },
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
    settings.value = await updateSettings({ panel_domain: settings.value.panel_domain })
    message.success('Domain updated')
  } catch {
    message.error('Failed to save settings')
  }
  saving.value = false
}

function changePassword() {
  if (!currentPassword.value || !newPassword.value) {
    message.error('Please fill all fields')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    message.error('Passwords do not match')
    return
  }
  message.success('Password changed (API not implemented yet)')
}

onMounted(() => {
  if (auth.user?.role === 'admin') loadSettings()
})
</script>

<template>
  <div>
    <PageHeader title="Settings" />

    <div class="flex flex-col md:flex-row gap-0">
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

      <div class="flex-1 min-w-0 pl-0 md:pl-7 pt-5 md:pt-0">
        <!-- Profile -->
        <div v-if="activeSection === 'profile'" class="max-w-[480px]">
          <div class="flex items-center gap-4 mb-6">
            <div class="w-13 h-13 rounded-full bg-accent-muted flex items-center justify-center shrink-0">
              <span class="font-heading text-[1.35rem] font-semibold text-accent">{{ (auth.user?.username || '?')[0].toUpperCase() }}</span>
            </div>
            <div class="flex flex-col gap-1.5">
              <div class="text-[1.1rem] font-semibold text-text-primary">{{ auth.user?.username }}</div>
              <span
                class="text-[11px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded w-fit"
                :class="auth.user?.role === 'admin'
                  ? 'text-danger bg-[rgba(239,68,68,0.1)]'
                  : 'text-text-secondary bg-[rgba(255,255,255,0.05)]'"
              >
                {{ auth.user?.role }}
              </span>
            </div>
          </div>

          <n-descriptions label-placement="left" :column="1" size="small">
            <n-descriptions-item label="Username">
              {{ auth.user?.username }}
            </n-descriptions-item>
            <n-descriptions-item label="Role">
              <span
                class="font-medium"
                :style="{ color: auth.user?.role === 'admin' ? 'var(--claude-danger, #ef4444)' : 'var(--claude-text-secondary, #9c9a92)' }"
              >
                {{ auth.user?.role }}
              </span>
            </n-descriptions-item>
          </n-descriptions>
        </div>

        <!-- Panel -->
        <div v-if="activeSection === 'panel'" class="max-w-[480px]">
          <div class="mb-6">
            <h3 class="font-heading text-[1.15rem] font-medium text-text-primary mb-1">Panel Domain</h3>
            <p class="text-[13px] text-text-tertiary m-0">Configure the domain name used to access this panel</p>
          </div>
          <n-form label-placement="top">
            <n-form-item label="Domain">
              <n-input
                v-model:value="settings.panel_domain"
                placeholder="panel.example.com"
                :disabled="loading"
                clearable
              >
                <template #prefix>
                  <n-icon :size="14"><server-outline /></n-icon>
                </template>
              </n-input>
            </n-form-item>
            <div class="flex justify-end pt-1">
              <n-button type="primary" :loading="saving" @click="saveDomain">
                <template #icon>
                  <n-icon :size="16"><checkmark-circle-outline /></n-icon>
                </template>
                Save Changes
              </n-button>
            </div>
          </n-form>
        </div>

        <!-- Security -->
        <div v-if="activeSection === 'security'" class="max-w-[480px]">
          <div class="mb-6">
            <h3 class="font-heading text-[1.15rem] font-medium text-text-primary mb-1">Change Password</h3>
            <p class="text-[13px] text-text-tertiary m-0">Update your account password</p>
          </div>
          <n-form label-placement="top">
            <n-form-item label="Current Password">
              <n-input v-model:value="currentPassword" type="password" placeholder="Enter current password" />
            </n-form-item>
            <n-form-item label="New Password">
              <n-input v-model:value="newPassword" type="password" placeholder="Enter new password" />
            </n-form-item>
            <n-form-item label="Confirm Password">
              <n-input v-model:value="confirmPassword" type="password" placeholder="Re-enter new password" />
            </n-form-item>
            <div class="flex justify-end pt-1">
              <n-button type="primary" @click="changePassword">
                <template #icon>
                  <n-icon :size="16"><lock-closed-outline /></n-icon>
                </template>
                Update Password
              </n-button>
            </div>
          </n-form>
        </div>
      </div>
    </div>
  </div>
</template>
