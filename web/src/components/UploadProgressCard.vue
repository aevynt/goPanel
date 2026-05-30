<script setup lang="ts">
defineProps<{
  fileName: string
  fileTotal: string
  progress: number
  speed: string
  status?: 'uploading' | 'done' | 'error'
  error?: string
}>()
</script>

<template>
  <div
    class="flex items-center justify-between gap-4 px-4 py-3 mb-4 rounded-lg opacity-[0.85] transition-opacity duration-200"
    :class="[
      'bg-surface border border-border-light',
      status === 'done' ? 'opacity-50' : '',
      status === 'error' ? 'border-danger' : '',
    ]"
  >
    <div class="flex items-center gap-3 min-w-0 flex-1">
      <n-icon :size="22" class="shrink-0 text-accent"><slot name="icon" /></n-icon>
      <div class="flex flex-col gap-[2px] min-w-0">
        <span class="text-sm font-medium text-text-primary truncate">{{ fileName }}</span>
        <span class="text-xs text-text-tertiary">{{ fileTotal }}</span>
      </div>
    </div>

    <div class="shrink-0 w-[180px]">
      <template v-if="status === 'uploading' || !status">
        <n-progress
          type="line"
          :percentage="progress"
          :height="5"
          :border-radius="3"
          indicator-placement="inside"
          :show-indicator="false"
        />
        <div class="flex justify-between mt-1 text-[11px] text-text-tertiary">
          <span class="tabular-nums">{{ speed }}</span>
          <span class="tabular-nums font-semibold text-accent">{{ progress }}%</span>
        </div>
      </template>
      <span v-else-if="status === 'done'" class="text-sm font-medium text-success">Done</span>
      <span v-else class="text-sm font-medium text-danger">{{ error || 'Failed' }}</span>
    </div>
  </div>
</template>
