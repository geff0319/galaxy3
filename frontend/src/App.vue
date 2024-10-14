<script setup lang="ts">
import {ref} from 'vue'

import * as Stores from '@/stores'
import {
  Events,
  Window
} from '@/bridge'
import { exitApp, ignoredError, sleep } from '@/utils'
import { useMessage, usePicker, useConfirm, usePrompt, useAlert } from '@/hooks'

import AboutView from '@/views/AboutView.vue'
import SplashView from '@/views/SplashView.vue'
import CommandView from './views/CommandView.vue'
import { NavigationBar, MainPage, TitleBar } from '@/components'
import WidgetsView from "@/views/WidgetsView.vue";
import { useRouter,useRoute } from 'vue-router';

const route = useRoute()

const loading =ref(true)
const envStore = Stores.useEnvStore()
const appStore = Stores.useAppStore()
const pluginsStore = Stores.usePluginsStore()
const appSettings = Stores.useAppSettingsStore()
const scheduledTasksStore = Stores.useScheduledTasksStore()
const wsClientStore = Stores.useWsClientStore()

const { message } = useMessage()
const { picker } = usePicker()
const { confirm } = useConfirm()
const { prompt } = usePrompt()
const { alert } = useAlert()

window.Plugins.message = message
window.Plugins.picker = picker
window.Plugins.confirm = confirm
window.Plugins.prompt = prompt
window.Plugins.alert = alert


// EventsOn('launchArgs', async (args: string[]) => {
//   console.log('launchArgs', args)
//   const url = new URL(args[0])
//   if (url.pathname === '//install-config/') {
//     const _url = url.searchParams.get('url')
//     const _name = url.searchParams.get('name') || sampleID()
//
//     if (!_url) {
//       message.error('URL missing')
//       return
//     }
//
//     try {
//       await subscribesStore.importSubscribe(_name, _url)
//       message.success('common.success')
//     } catch (error: any) {
//       message.error(error)
//     }
//   }
// })

Events.On('beforeClose', async () => {
  exitApp()
})
console.log(route.path)

window.addEventListener('beforeunload', scheduledTasksStore.removeScheduledTasks)


appSettings.setupAppSettings().then(async () => {
  await Promise.all([
    ignoredError(envStore.setupEnv),
    ignoredError(pluginsStore.setupPlugins),
    ignoredError(scheduledTasksStore.setupScheduledTasks),
    ignoredError(wsClientStore.setupWsSettings)
  ])
  await sleep(1000)

  loading.value = false

  try {
    await pluginsStore.onStartupTrigger()
  } catch (error: any) {
    message.error(error)
  }
})

</script>

<template>
  <SplashView v-if="loading && route.path !== '/ytdlpWidgets'" />
<!--  <template v-else-if="appStore.widgetsEnable"><widgetsView /></template>-->
<!--  <template v-else-if="!appStore.widgetsEnable">-->
  <template v-else-if="!route.meta.newlayout">
    <TitleBar />
    <div class="main">
      <NavigationBar />
      <MainPage />
    </div>
  </template>
  <template v-else-if="route.meta.newlayout"><WidgetsView /></template>

  <Modal
    v-model:open="appStore.showAbout"
    :cancel="false"
    :submit="false"
    mask-closable
    min-width="50"
  >
    <AboutView />
  </Modal>

  <Menu
    v-model="appStore.menuShow"
    :position="appStore.menuPosition"
    :menu-list="appStore.menuList"
  />

  <Tips
    v-model="appStore.tipsShow"
    :position="appStore.tipsPosition"
    :message="appStore.tipsMessage"
  />

  <CommandView />
</template>

<style scoped>
.main {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  padding: 8px;
}
</style>
