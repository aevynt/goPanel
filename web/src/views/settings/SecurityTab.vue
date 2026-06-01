<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { LockClosedOutline, ShieldCheckmarkOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'
import { changePassword as apiChangePassword, setup2FA, enable2FA, disable2FA } from '@/api'

const auth = useAuthStore()
const message = useMessage()

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const is2FAEnabled = computed(() => auth.user?.totp_enabled || false)
const show2FAWizard = ref(false)
const totpSecret = ref('')
const totpUri = ref('')
const verifyCode = ref('')
const disablePassword = ref('')
const loading2FA = ref(false)
const saving = ref(false)

async function changePassword() {
  if (!currentPassword.value || !newPassword.value) {
    message.error('Please fill all fields')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    message.error('Passwords do not match')
    return
  }
  saving.value = true
  try {
    await apiChangePassword(currentPassword.value, newPassword.value)
    message.success('Password changed successfully')
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Failed to change password')
  }
  saving.value = false
}

async function start2FASetup() {
  loading2FA.value = true
  try {
    const res = await setup2FA()
    totpSecret.value = res.secret
    totpUri.value = res.provisioning_uri
    show2FAWizard.value = true
  } catch {
    message.error('Failed to start 2FA setup')
  }
  loading2FA.value = false
}

async function confirm2FA() {
  if (!verifyCode.value) {
    message.error('Please enter the 6-digit code')
    return
  }
  loading2FA.value = true
  try {
    await enable2FA(verifyCode.value)
    message.success('Two-Factor Authentication activated!')
    show2FAWizard.value = false
    verifyCode.value = ''
    await auth.checkAuth()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Verification failed')
  }
  loading2FA.value = false
}

async function deactivate2FA() {
  if (!disablePassword.value) {
    message.error('Password is required to disable 2FA')
    return
  }
  loading2FA.value = true
  try {
    await disable2FA(disablePassword.value)
    message.success('Two-Factor Authentication deactivated')
    disablePassword.value = ''
    await auth.checkAuth()
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Deactivation failed')
  }
  loading2FA.value = false
}
</script>

<template>
  <div class="max-w-[480px]">
    <!-- Change Password Form -->
    <div class="mb-6">
      <h3 class="font-heading text-[1.15rem] font-medium text-text-primary mb-1">Change Password</h3>
      <p class="text-[13px] text-text-tertiary m-0">Update your account password</p>
    </div>
    <n-form label-placement="top">
      <n-form-item label="Current Password">
        <n-input v-model:value="currentPassword" type="password" placeholder="Enter current password" />
      </n-form-item>
      <n-form-item label="New Password">
        <n-input v-model:value="newPassword" type="password" placeholder="Enter new password" />
      </n-form-item>
      <n-form-item label="Confirm Password">
        <n-input v-model:value="confirmPassword" type="password" placeholder="Re-enter new password" />
      </n-form-item>
      <div class="flex justify-end pt-1">
        <n-button type="primary" :loading="saving" @click="changePassword">
          <template #icon>
            <n-icon :size="16"><lock-closed-outline /></n-icon>
          </template>
          Update Password
        </n-button>
      </div>
    </n-form>

    <!-- Two-Factor Authentication Section -->
    <div class="border-t border-border-light my-8 pt-6">
      <h3 class="font-heading text-[1.15rem] font-medium text-text-primary mb-1">Two-Factor Authentication (2FA)</h3>
      <p class="text-[13px] text-text-tertiary mb-4">Protect your account with an extra layer of security using TOTP (Time-based One-Time Password)</p>

      <!-- 2FA is currently active -->
      <div v-if="is2FAEnabled" class="glass-card p-4 flex flex-col gap-4" style="background: rgba(31, 31, 29, 0.6); backdrop-filter: blur(12px); border: 1px solid rgba(255, 255, 255, 0.05); border-radius: 8px;">
        <div class="flex items-center gap-3">
          <n-icon :size="24" color="#16a34a"><shield-checkmark-outline /></n-icon>
          <div>
            <div class="text-sm font-semibold text-text-primary">2FA is active</div>
            <div class="text-[11px] text-text-tertiary">Your account is fully protected.</div>
          </div>
        </div>

        <div class="border-t border border-border-light pt-4 mt-2">
          <div class="text-xs font-semibold text-text-primary mb-2">Disable 2FA</div>
          <n-form-item label="Enter your Password to disable">
            <n-input v-model:value="disablePassword" type="password" placeholder="Confirm your password" />
          </n-form-item>
          <n-button type="error" :loading="loading2FA" @click="deactivate2FA">
            Disable 2FA
          </n-button>
        </div>
      </div>

      <!-- 2FA is NOT active -->
      <div v-else-if="!show2FAWizard" class="glass-card p-4" style="background: rgba(31, 31, 29, 0.6); backdrop-filter: blur(12px); border: 1px solid rgba(255, 255, 255, 0.05); border-radius: 8px;">
        <p class="text-xs text-text-secondary mb-4">2FA is currently disabled. We strongly recommend enabling it for administrators.</p>
        <n-button type="primary" :loading="loading2FA" @click="start2FASetup">
          Enable 2FA
        </n-button>
      </div>

      <!-- 2FA Setup Wizard -->
      <div v-else class="glass-card p-4 flex flex-col gap-4" style="background: rgba(31, 31, 29, 0.6); backdrop-filter: blur(12px); border: 1px solid rgba(255, 255, 255, 0.05); border-radius: 8px;">
        <div class="text-sm font-semibold text-text-primary">Setup 2FA Authenticator</div>
        <p class="text-xs text-text-secondary m-0">
          1. Enter the secret code below into your Authenticator App (e.g. Google Authenticator, Authy):
        </p>
        <div class="code-block" style="background: rgba(0, 0, 0, 0.3); padding: 10px; border-radius: 6px; font-family: monospace; font-size: 14px; text-align: center; letter-spacing: 2px; color: var(--claude-accent);">
          {{ totpSecret }}
        </div>
        <p class="text-xs text-text-secondary m-0">
          Or click here to open directly in your client:
          <a :href="totpUri" class="text-accent underline text-xs" style="color: var(--claude-accent, #c96442);">Open in Authenticator</a>
        </p>
        <p class="text-xs text-text-secondary m-0">
          2. Enter the 6-digit verification code generated by your app:
        </p>
        <n-form-item label="Verification Code">
          <n-input v-model:value="verifyCode" placeholder="e.g. 123456" maxlength="6" />
        </n-form-item>

        <n-space justify="end">
          <n-button secondary @click="show2FAWizard = false">Cancel</n-button>
          <n-button type="primary" :loading="loading2FA" @click="confirm2FA">Activate 2FA</n-button>
        </n-space>
      </div>
    </div>
  </div>
</template>
