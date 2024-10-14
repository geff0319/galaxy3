<script setup lang="ts">
import { ref } from 'vue'

import { APP_TITLE, APP_VERSION, exitApp } from '@/utils'
import {type Menu, useAppSettingsStore, useKernelApiStore, useEnvStore} from '@/stores'
import {
  WML,
  Window,
} from '@/bridge'
import * as Stores from "@/stores";

const emits = defineEmits(['changeWin'])

const isPinned = ref(false)
const isFullScreen = ref(false)

const appSettingsStore = useAppSettingsStore()
const kernelApiStore = useKernelApiStore()
const envStore = useEnvStore()
const appStore = Stores.useAppStore()



const pinWindow = () => {
  isPinned.value = !isPinned.value
  Window.SetAlwaysOnTop(isPinned.value)
  // WindowSetAlwaysOnTop(isPinned.value)
}
const enableWidgets = async() => {
  const win = Window.Get("WidgetsWin")
  await win.SetRelativePosition(1400, 100).then(win.Show)
  // await win.OpenDevTools()
  // ShowWidgets()
}



const closeWindow = async () => {
  if (appSettingsStore.app.exitOnClose) {
    // await Window.Get("WidgetsWin").Close()
    await Window.Close()
    // exitApp()
  } else {
    await Window.Hide()
  }
}

const menus: Menu[] = [
  {
    label: 'titlebar.reload',
    handler: WML.Reload
  }
]
</script>

<template>
  <div
    v-if="envStore.env.os === 'windows'"
    class="titlebar"
    style="--wails-draggable: drag"
  >
    <img class="logo" draggable="false" src="@/assets/logo.png" />
    <div
      :style="{
        color: appSettingsStore.app.kernel.running ? 'var(--primary-color)' : 'var(--color)'
      }"
      class="appname"
    >
      {{ APP_TITLE }} {{ APP_VERSION }}
    </div>
    <Button v-if="kernelApiStore.loading" loading type="text" size="small" />
    <div v-menu="menus" class="menus"></div>
    <div class="action" style="--wails-draggable: disabled">
<!--      <Dropdown :trigger="['hover', 'click']">-->
<!--        <Button type="link" size="small">-->
<!--          悬浮窗-->
<!--        </Button>-->
<!--        <template #overlay>-->
<!--          <Button type="link" size="small" @click.stop="enableWidgets">-->
<!--            ⏰时钟-->
<!--          </Button>-->
<!--        </template>-->
<!--      </Dropdown>-->
      <Button @click.stop="pinWindow" type="text">
        <Icon :icon="isPinned ? 'pinFill' : 'pin'" />
      </Button>
      <Button @click.stop="Window.Hide" type="text">
        <Icon icon="minimize" />
      </Button>
      <Button @click.stop="enableWidgets" type="text">
        <Icon :icon="isFullScreen ? 'maximize2' : 'maximize'" />
      </Button>
      <Button
        @click.stop="closeWindow"
        :class="{ 'hover-red': appSettingsStore.app.exitOnClose }"
        type="text"
      >
        <Icon icon="close" />
      </Button>
    </div>
  </div>
  <div v-else class="placeholder" style="--wails-draggable: drag">
    <div
      :style="{
        color: appSettingsStore.app.kernel.running ? 'var(--primary-color)' : 'var(--color)'
      }"
      v-menu="menus"
      class="appname"
    >
      {{ APP_TITLE }} {{ APP_VERSION }}
    </div>
    <Button v-if="kernelApiStore.loading" loading type="text" size="small" />
  </div>
</template>

<style lang="less" scoped>
.titlebar {
  user-select: none;
  display: flex;
  padding: 4px 12px;
  align-items: center;
}
.logo {
  width: 24px;
  height: 24px;
  user-select: none;
}
.appname {
  font-size: 14px;
  margin-left: 8px;
  margin-top: 2px;
  font-weight: bold;
}

.menus {
  flex: 1;
  height: 100%;
}
.action {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  &-btn {
    width: 32px;
    height: 32px;
    line-height: 32px;
    text-align: center;
    border-radius: 4px;
    &:hover {
      background-color: var(--hover-bg-color);
    }
  }
}
.placeholder {
  user-select: none;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  .appname {
    font-size: 12px;
  }
}

.hover-red:hover {
  background: rgba(255, 0, 0, 0.6);
}
</style>
