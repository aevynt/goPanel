<script setup lang="ts">
import { ref } from 'vue'
import { NTabPane, NTabs } from 'naive-ui'
import PageHeader from '@/components/PageHeader.vue'
import ContainerList from '@/components/docker/ContainerList.vue'
import ComposeDeploy from '@/components/docker/ComposeDeploy.vue'

const tabValue = ref('containers')
const listRef = ref<InstanceType<typeof ContainerList> | null>(null)

function onDeployed() {
  tabValue.value = 'containers'
  // reload container list after a short delay to allow background docker compose setup to launch
  setTimeout(() => {
    listRef.value?.load()
  }, 2000)
}
</script>

<template>
  <div>
    <PageHeader title="Docker Manager" />

    <n-tabs v-model:value="tabValue" type="line">
      <n-tab-pane name="containers" tab="Containers">
        <ContainerList ref="listRef" />
      </n-tab-pane>

      <n-tab-pane name="deploy" tab="Deploy Compose">
        <ComposeDeploy @deployed="onDeployed" />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>
