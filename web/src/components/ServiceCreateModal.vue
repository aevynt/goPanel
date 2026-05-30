<script setup lang="ts">
import { ref } from 'vue'
import type { ServiceSpec, Binary } from '@/api/client'
import ModalFooter from '@/components/ModalFooter.vue'

const props = defineProps<{
  show: boolean
  binaries: Binary[]
}>()

const emit = defineEmits<{
  'update:show': [v: boolean]
  create: [spec: ServiceSpec]
}>()

const form = ref<ServiceSpec>({
  name: '',
  description: '',
  binary_path: '',
  working_dir: '',
  port: 8080,
  env_vars: '',
  args: '',
  auto_start: true,
  run_as: '',
})

const selectedBinaryId = ref<number | null>(null)
const envItems = ref<{ key: string; value: string }[]>([])

function resetForm() {
  form.value = {
    name: '',
    description: '',
    binary_path: '',
    working_dir: '',
    port: 8080,
    env_vars: '',
    args: '',
    auto_start: true,
    run_as: '',
  }
  selectedBinaryId.value = null
  envItems.value = []
}

function onBinaryChange(val: number | null) {
  selectedBinaryId.value = val
  const bin = props.binaries.find(b => b.id === val)
  if (bin) {
    form.value.binary_path = bin.path
  }
}

function addEnv() {
  envItems.value.push({ key: '', value: '' })
}

function removeEnv(idx: number) {
  envItems.value.splice(idx, 1)
}

function handleCreate() {
  form.value.env_vars = envItems.value
    .filter(e => e.key && e.value)
    .map(e => `${e.key}=${e.value}`)
    .join('\n')
  emit('create', { ...form.value })
}
</script>

<template>
  <n-modal
    :show="show"
    title="Create Service"
    preset="card"
    style="width: 560px; max-width: 90vw;"
    :mask-closable="false"
    @after-leave="resetForm"
    @update:show="emit('update:show', $event)"
  >
    <n-form label-placement="top">
      <n-form-item label="Name">
        <n-input v-model:value="form.name" placeholder="my-service" />
      </n-form-item>
      <n-form-item label="Description">
        <n-input v-model:value="form.description" placeholder="My Go service" />
      </n-form-item>
      <n-form-item label="Binary">
        <n-select
          v-model:value="selectedBinaryId"
          :options="binaries.map(b => ({ label: b.name, value: b.id }))"
          placeholder="Select a binary"
          filterable
          @update:value="onBinaryChange"
        />
      </n-form-item>
      <n-form-item label="Port" class="max-w-[140px]">
        <n-input-number v-model:value="form.port" :min="1" :max="65535" class="w-full" />
      </n-form-item>
      <n-form-item label="Environment Variables">
        <div class="w-full">
          <div
            v-for="(item, idx) in envItems"
            :key="idx"
            class="flex gap-2 mb-2 items-center"
          >
            <n-input v-model:value="item.key" placeholder="KEY" class="flex-1 max-w-[160px]" />
            <n-input v-model:value="item.value" placeholder="value" class="flex-[2]" />
            <n-button quaternary circle size="tiny" type="error" @click="removeEnv(idx)">
              <template #icon>
                <n-icon>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </n-icon>
              </template>
            </n-button>
          </div>
          <n-button size="tiny" quaternary type="primary" @click="addEnv">
            + Add variable
          </n-button>
        </div>
      </n-form-item>

      <n-collapse :default-expanded-names="[]">
        <n-collapse-item title="Advanced" name="advanced">
          <n-form-item label="Args">
            <n-input v-model:value="form.args" placeholder="--port 8080" />
          </n-form-item>
          <n-form-item label="Working Directory">
            <n-input v-model:value="form.working_dir" placeholder="/var/lib/gopanel" />
          </n-form-item>
          <n-form-item label="Run As">
            <n-input v-model:value="form.run_as" placeholder="root" />
          </n-form-item>
          <n-form-item label="Binary Path">
            <n-input v-model:value="form.binary_path" placeholder="/var/lib/gopanel/binaries/myapp" />
          </n-form-item>
          <n-form-item label=" ">
            <n-checkbox v-model:checked="form.auto_start">
              Auto start on boot
            </n-checkbox>
          </n-form-item>
        </n-collapse-item>
      </n-collapse>
    </n-form>
    <template #footer>
      <ModalFooter submit-label="Create" @cancel="emit('update:show', false)" @submit="handleCreate" />
    </template>
  </n-modal>
</template>
