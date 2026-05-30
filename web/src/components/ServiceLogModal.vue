<script setup lang="ts">
defineProps<{
  show: boolean
  serviceName: string
  logs: { timestamp: string; message: string }[]
  loading: boolean
}>()

const emit = defineEmits<{
  'update:show': [v: boolean]
}>()
</script>

<template>
  <n-modal
    :show="show"
    preset="card"
    style="width: 800px; max-width: 90vw;"
    :mask-closable="false"
    @update:show="emit('update:show', $event)"
  >
    <template #header>
      <span class="font-heading font-medium text-base">Logs: {{ serviceName }}</span>
    </template>
    <n-log
      :loading="loading"
      :rows="logs.map(l => l.timestamp ? `[${l.timestamp}] ${l.message}` : l.message).join('\n')"
      trim
      :font-size="12"
      style="height: 400px;"
    />
  </n-modal>
</template>
