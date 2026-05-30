<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <div class="flex justify-between items-center mb-6">
    <h2 class="font-heading font-medium text-[1.5rem] m-0">{{ title }}</h2>
    <div v-if="$slots.actions" class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
