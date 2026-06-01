<script setup lang="ts">
import { CheckmarkCircleOutline } from '@vicons/ionicons5'
import type { PanelSettings } from '@/api/client'

const props = defineProps<{
  settings: PanelSettings
  loading: boolean
  saving: boolean
}>()

const emit = defineEmits<{
  (e: 'save-alerts'): void
}>()
</script>

<template>
  <div class="max-w-[480px]">
    <div class="mb-6">
      <h3 class="font-heading text-[1.15rem] font-medium text-text-primary mb-1">Alerting Notifications</h3>
      <p class="text-[13px] text-text-tertiary m-0">Receive high temperature or resource usage notifications directly to Discord or Telegram</p>
    </div>
    <n-form label-placement="top">
      <h4 class="text-sm font-semibold text-text-primary mb-3">Discord Integration</h4>
      <n-form-item label="Discord Webhook URL">
        <n-input
          v-model:value="props.settings.discord_webhook"
          placeholder="https://discord.com/api/webhooks/..."
          :disabled="props.loading"
          clearable
        />
      </n-form-item>

      <h4 class="text-sm font-semibold text-text-primary mt-6 mb-3">Telegram Integration</h4>
      <n-form-item label="Telegram Bot Token">
        <n-input
          v-model:value="props.settings.telegram_token"
          placeholder="1234567890:ABCdefGhI..."
          :disabled="props.loading"
          clearable
        />
      </n-form-item>
      <n-form-item label="Telegram Chat ID">
        <n-input
          v-model:value="props.settings.telegram_chat_id"
          placeholder="-100123456789"
          :disabled="props.loading"
          clearable
        />
      </n-form-item>

      <h4 class="text-sm font-semibold text-text-primary mt-6 mb-3">Thresholds</h4>
      <n-grid :cols="2" :x-gap="12">
        <n-gi>
          <n-form-item label="CPU Temp (°C)">
            <n-input-number
              v-model:value="props.settings.alert_temp_threshold"
              :min="30" :max="100"
              :disabled="props.loading"
              style="width: 100%;"
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="CPU Usage (%)">
            <n-input-number
              v-model:value="props.settings.alert_cpu_threshold"
              :min="10" :max="100"
              :disabled="props.loading"
              style="width: 100%;"
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="RAM Usage (%)">
            <n-input-number
              v-model:value="props.settings.alert_ram_threshold"
              :min="10" :max="100"
              :disabled="props.loading"
              style="width: 100%;"
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="Disk Usage (%)">
            <n-input-number
              v-model:value="props.settings.alert_disk_threshold"
              :min="10" :max="100"
              :disabled="props.loading"
              style="width: 100%;"
            />
          </n-form-item>
        </n-gi>
      </n-grid>

      <div class="flex justify-end pt-4">
        <n-button type="primary" :loading="props.saving" @click="emit('save-alerts')">
          <template #icon>
            <n-icon :size="16"><checkmark-circle-outline /></n-icon>
          </template>
          Save Alerts Settings
        </n-button>
      </div>
    </n-form>
  </div>
</template>
