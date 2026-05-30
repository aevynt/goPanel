<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { darkTheme } from 'naive-ui'
import '@/assets/claude.css'

const auth = useAuthStore()

onMounted(() => {
  auth.checkAuth()
})

document.documentElement.dataset.theme = 'dark'

const themeMeta = document.createElement('meta')
themeMeta.name = 'theme-color'
themeMeta.content = '#141413'
document.head.appendChild(themeMeta)

const darkOverrides = {
  common: {
    fontFamily: "'DM Sans', system-ui, -apple-system, sans-serif",
    fontFamilyMono: "'JetBrains Mono', 'Fira Code', monospace",
    fontSize: '15px',
    fontWeight: '400',
    fontWeightStrong: '500',
    primaryColor: '#c96442',
    primaryColorHover: '#d4733f',
    primaryColorPressed: '#b85a38',
    primaryColorSuppl: '#c96442',
    successColor: '#16a34a',
    warningColor: '#d97706',
    errorColor: '#dc2626',
    bodyColor: '#141413',
    cardColor: '#1f1f1d',
    modalColor: '#1f1f1d',
    popoverColor: '#1f1f1d',
    tableColor: '#1f1f1d',
    inputColor: '#2a2a28',
    actionColor: '#1f1f1d',
    tabColor: '#1f1f1d',
    borderRadius: '8px',
    borderRadiusSmall: '4px',
    textColor1: '#e8e6dc',
    textColor2: '#9c9a92',
    textColor3: '#73726c',
    borderColor: '#3d3d3a',
    dividerColor: '#2a2a28',
    hoverColor: 'rgba(201, 100, 66, 0.1)',
    pressedColor: 'rgba(201, 100, 66, 0.18)',
    avatarColor: 'rgba(201, 100, 66, 0.2)',
    avatarTextColor: '#c96442',
    name: 'claude-dark',
    fontSizeSmall: '13px',
    fontSizeMedium: '15px',
    fontSizeLarge: '17px',
    heightSmall: '30px',
    heightMedium: '36px',
    heightLarge: '42px',
    opacityDisabled: '0.5',
    actionColor2: '#2a2a28',
    progressRailColor: '#3d3d3a',
    railColor: '#3d3d3a',
    popoverBoxShadow: '0 4px 24px rgba(0,0,0,0.4)',
    tableHeaderColor: '#1a1a18',
  },
  Button: {
    textColorGhost: '#e8e6dc',
    borderPrimary: '#c96442', borderHoverPrimary: '#d4733f', borderPressedPrimary: '#b85a38',
    color: '#c96442', colorHover: '#d4733f', colorPressed: '#b85a38', colorFocus: '#c96442',
    colorDisabled: '#2a2a28', textColor: '#ffffff', textColorHover: '#ffffff',
    textColorPressed: '#ffffff', textColorDisabled: '#73726c',
    fontWeight: '500', borderRadiusSmall: '6px', borderRadiusMedium: '8px',
    borderRadiusLarge: '10px', rippleColor: 'rgba(201, 100, 66, 0.3)',
  },
  Card: {
    color: '#1f1f1d', borderColor: '#2a2a28', borderRadius: '8px',
    paddingMedium: '24px', titleTextColor: '#e8e6dc', titleFontWeight: '500', actionColor: '#1f1f1d',
  },
  DataTable: {
    borderColor: '#2a2a28', tdColor: '#1f1f1d', thColor: '#1a1a18',
    thTextColor: '#9c9a92', thFontWeight: '500', tdColorHover: 'rgba(201,100,66,0.06)',
    tdTextColor: '#e8e6dc', borderRadius: '8px', loadingColor: '#c96442',
  },
  Menu: {
    itemColorActive: 'rgba(201,100,66,0.15)', itemTextColorActive: '#c96442',
    itemTextColor: '#9c9a92', itemTextColorHover: '#e8e6dc',
    itemColorHover: 'rgba(255,255,255,0.04)', itemIconColor: '#73726c',
    itemIconColorActive: '#c96442', itemIconColorHover: '#9c9a92',
    itemHeight: '40px', borderRadius: '6px', fontSize: '14px', arrowColor: '#73726c',
  },
  Input: {
    color: '#2a2a28', borderColor: '#3d3d3a', borderHover: '#c96442',
    borderFocus: '#c96442', placeholderColor: '#73726c', textColor: '#e8e6dc', borderRadius: '6px',
  },
  Select: {
    menuColor: '#1f1f1d', menuBorderColor: '#3d3d3a',
    menuBoxShadow: '0 4px 24px rgba(0,0,0,0.4)', color: '#2a2a28',
    borderColor: '#3d3d3a', borderHover: '#c96442', borderFocus: '#c96442',
  },
  Tag: { color: '#2a2a28', textColor: '#9c9a92', borderRadius: '4px', fontWeight: '500' },
  Progress: { railColor: '#3d3d3a', fontWeight: '500' },
  Form: { labelTextColor: '#9c9a92', labelFontWeight: '500' },
  Modal: {
    color: '#1f1f1d', borderColor: '#2a2a28', borderRadius: '12px',
    boxShadow: '0 8px 32px rgba(0,0,0,0.4)', textColor: '#e8e6dc',
    titleTextColor: '#e8e6dc', titleFontWeight: '500',
  },
  Breadcrumb: {
    itemTextColor: '#73726c', itemTextColorHover: '#c96442',
    itemTextColorPressed: '#b85a38', separatorColor: '#3d3d3a',
    itemTextColorActive: '#e8e6dc', fontWeightActive: '500',
  },
  Descriptions: {
    thColor: '#1a1a18', tdColor: '#1f1f1d', thTextColor: '#9c9a92',
    tdTextColor: '#e8e6dc', borderColor: '#2a2a28', thFontWeight: '500',
  },
  Alert: { color: '#1f1f1d', borderColor: '#3d3d3a', borderRadius: '8px' },
  Log: { loaderColor: '#c96442', loaderTextColor: '#9c9a92', loaderBorder: '#3d3d3a' },
  Spin: { color: '#c96442' },
  Switch: { railColorActive: '#c96442', railColor: '#3d3d3a' },
  Checkbox: {
    checkMarkColor: '#141413', borderActive: '#c96442', border: '#3d3d3a',
    colorActive: '#c96442', color: '#2a2a28',
  },
  Upload: { borderColor: '#3d3d3a', borderColorHover: '#c96442', borderRadius: '8px' },
  Notification: { color: '#1f1f1d', borderColor: '#3d3d3a', borderRadius: '8px' },
  Message: { color: '#1f1f1d', borderColor: '#3d3d3a', borderRadius: '8px' },
  Dialog: { color: '#1f1f1d', borderColor: '#3d3d3a', borderRadius: '12px' },
  Popconfirm: { color: '#1f1f1d', borderColor: '#3d3d3a', borderRadius: '8px' },
  Scrollbar: { color: '#3d3d3a', colorHover: '#5e5d59' },
  Tooltip: { color: '#e8e6dc', textColor: '#141413', borderRadius: '6px' },
  Pagination: {
    itemColor: 'transparent', itemColorHover: 'rgba(201,100,66,0.1)',
    itemColorActive: 'rgba(201,100,66,0.2)', itemTextColor: '#9c9a92',
    itemTextColorHover: '#e8e6dc', itemTextColorActive: '#c96442',
    buttonColor: 'transparent', buttonColorHover: 'rgba(201,100,66,0.1)',
  },
  Dropdown: {
    color: '#1f1f1d', borderColor: '#3d3d3a', borderRadius: '8px',
    optionTextColor: '#e8e6dc', optionTextColorHover: '#e8e6dc',
    optionColorHover: 'rgba(201,100,66,0.1)', optionTextColorActive: '#c96442',
    optionColorActive: 'rgba(201,100,66,0.15)',
  },
  Badge: { color: '#c96442' },
  Empty: { textColor: '#73726c', iconColor: '#3d3d3a' },
}
</script>

<template>
  <n-config-provider :theme="darkTheme" :theme-overrides="darkOverrides">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <router-view />
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>
