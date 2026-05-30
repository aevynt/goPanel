<script setup lang="ts">
defineProps<{
  submitLabel?: string
  submitLoading?: boolean
  cancelLabel?: string
}>()

const emit = defineEmits<{
  cancel: []
  submit: []
}>()
</script>

<template>
  <div class="flex justify-end gap-2">
    <n-button quaternary @click="emit('cancel')">{{ cancelLabel || 'Cancel' }}</n-button>
    <n-button type="primary" :loading="submitLoading" @click="emit('submit')">{{ submitLabel || 'Submit' }}</n-button>
  </div>
</template>
