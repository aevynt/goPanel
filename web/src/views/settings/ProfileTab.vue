<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
</script>

<template>
  <div class="max-w-[480px]">
    <div class="flex items-center gap-4 mb-6">
      <div class="w-13 h-13 rounded-full bg-accent-muted flex items-center justify-center shrink-0" style="width: 52px; height: 52px; border-radius: 50%; background: rgba(201, 100, 66, 0.1); display: flex; align-items: center; justify-content: center;">
        <span class="font-heading text-[1.35rem] font-semibold text-accent" style="color: var(--claude-accent, #c96442); font-weight: 600;">{{ (auth.user?.username || '?')[0].toUpperCase() }}</span>
      </div>
      <div class="flex flex-col gap-1.5">
        <div class="text-[1.1rem] font-semibold text-text-primary">{{ auth.user?.username }}</div>
        <span
          class="text-[11px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded w-fit"
          :style="{
            color: auth.user?.role === 'admin' ? '#ef4444' : '#9c9a92',
            background: auth.user?.role === 'admin' ? 'rgba(239,68,68,0.1)' : 'rgba(255,255,255,0.05)'
          }"
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
</template>
