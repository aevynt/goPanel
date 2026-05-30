<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { useMessage, NButton, NSpace, NPopconfirm } from 'naive-ui'
import { listUsers, createUser, deleteUser } from '@/api'
import type { User } from '@/api/client'
import PageHeader from '@/components/PageHeader.vue'
import ModalFooter from '@/components/ModalFooter.vue'

const message = useMessage()
const users = ref<User[]>([])
const loading = ref(false)
const showCreateModal = ref(false)
const newUser = ref({ username: '', password: '', role: 'viewer' })

async function load() {
  loading.value = true
  try {
    users.value = await listUsers()
  } catch {
    message.error('Failed to load users')
  }
  loading.value = false
}

async function handleCreate() {
  try {
    await createUser(newUser.value)
    message.success('User created')
    showCreateModal.value = false
    newUser.value = { username: '', password: '', role: 'viewer' }
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to create user')
  }
}

async function handleDelete(id: number) {
  try {
    await deleteUser(id)
    message.success('User deleted')
    load()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to delete user')
  }
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: 'Username', key: 'username' },
  {
    title: 'Role',
    key: 'role',
    render: (row: User) => {
      const color = row.role === 'admin' ? 'var(--claude-danger, #dc2626)' : 'var(--claude-text-secondary, #5e5d59)'
      return h('span', { style: { color, fontSize: '13px', fontWeight: '500' } }, row.role)
    },
  },
  {
    title: '',
    key: 'actions',
    width: 80,
    render: (row: User) => row.id === 1
      ? null
      : h(NPopconfirm, { onPositiveClick: () => handleDelete(row.id) }, {
          default: () => 'Delete this user?',
          trigger: () => h(NButton, { size: 'tiny', quaternary: true, type: 'error' }, { default: () => 'Delete' }),
        }),
  },
]

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Users">
      <template #actions>
        <n-button type="primary" @click="showCreateModal = true">Create User</n-button>
      </template>
    </PageHeader>

    <n-data-table
      :loading="loading"
      :columns="columns"
      :data="users"
      :single-line="false"
      :scroll-x="400"
    />

    <n-modal v-model:show="showCreateModal" title="Create User" preset="card" style="width: 400px; max-width: 90vw;" :mask-closable="false">
      <n-form>
        <n-form-item label="Username">
          <n-input v-model:value="newUser.username" placeholder="username" />
        </n-form-item>
        <n-form-item label="Password">
          <n-input v-model:value="newUser.password" type="password" placeholder="••••••" />
        </n-form-item>
        <n-form-item label="Role">
          <n-select v-model:value="newUser.role" :options="[
            { label: 'Admin', value: 'admin' },
            { label: 'Viewer', value: 'viewer' },
          ]" />
        </n-form-item>
      </n-form>
      <template #footer>
        <ModalFooter submit-label="Create" @cancel="showCreateModal = false" @submit="handleCreate" />
      </template>
    </n-modal>
  </div>
</template>
