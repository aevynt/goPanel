<script setup lang="ts">
import { ServerOutline, CheckmarkCircleOutline } from '@vicons/ionicons5'
import type { PanelSettings } from '@/api/client'

const props = defineProps<{
  settings: PanelSettings
  loading: boolean
  saving: boolean
}>()

const emit = defineEmits<{
  (e: 'save-domain'): void
}>()
</script>

<template>
  <div class="max-w-[480px]">
    <div class="mb-6">
      <h3 class="font-heading text-[1.15rem] font-medium text-text-primary mb-1">Panel Domain</h3>
      <p class="text-[13px] text-text-tertiary m-0">Configure the domain name used to access this panel</p>
    </div>
    <n-form label-placement="top">
      <n-form-item label="Domain">
        <n-input
          v-model:value="props.settings.panel_domain"
          placeholder="panel.example.com"
          :disabled="props.loading"
          clearable
        >
          <template #prefix>
            <n-icon :size="14"><server-outline /></n-icon>
          </template>
        </n-input>
      </n-form-item>
      <div class="flex justify-end pt-1">
        <n-button type="primary" :loading="props.saving" @click="emit('save-domain')">
          <template #icon>
            <n-icon :size="16"><checkmark-circle-outline /></n-icon>
          </template>
          Save Changes
        </n-button>
      </div>
    </n-form>
  </div>
</template>
