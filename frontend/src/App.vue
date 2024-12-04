<script setup lang="ts">
import {ref} from 'vue'
//
import * as Stores from '@/stores'
import {
  Events, ExitApp,
  Window
} from '@/bridge'
import { ignoredError, sleep } from '@/utils'
import { useMessage, usePicker, useConfirm, usePrompt, useAlert } from '@/hooks'

import AboutView from '@/views/AboutView.vue'
import SplashView from '@/views/SplashView.vue'
import CommandView from './views/CommandView.vue'
import { NavigationBar,MainPage, TitleBar } from '@/components'
import WidgetsView from "@/views/WidgetsView.vue";
import { useRouter,useRoute } from 'vue-router';
import {message as antmessage} from "ant-design-vue/es/components";
const route = useRoute();
const router = useRouter();

const loading =ref(false)
const envStore = Stores.useEnvStore()
const appStore = Stores.useAppStore()
const pluginsStore = Stores.usePluginsStore()
const appSettings = Stores.useAppSettingsStore()
const scheduledTasksStore = Stores.useScheduledTasksStore()
const wsClientStore = Stores.useWsClientStore()
const mqttClientStore = Stores.useMqttClientStore()

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

Events.Once("appInit", function(event:any) {
  window.addEventListener('beforeunload', scheduledTasksStore.removeScheduledTasks)
  appSettings.setupAppSettings().then(async () => {
    await Promise.all([
      ignoredError(envStore.setupEnv),
      ignoredError(pluginsStore.setupPlugins),
      ignoredError(scheduledTasksStore.setupScheduledTasks),
    ])
    await sleep(1000)

    try {
      await pluginsStore.onStartupTrigger()
    } catch (error: any) {
      message.error(error)
    }
  })
  Events.On('beforeClose',  () => {
    // exitApp()
    ExitApp()
  })
  // 0：是否切换页面 1：等级 2：消息
  Events.On('notify',(event:any)=>{
    if(event.data[0]){
      router.push({path: '/ytdlp'})
      const w = Window.Get('MainWin')
      w.IsFocused().then((status)=>{
        console.log(status)
        if (!status){
          w.Show()
          w.Focus()
        }
      })
    }
    console.log(event)
    switch (event.data[1]) {
      case "info":
        antmessage.info(event.data[2],3);
        break;
      case "error":
        antmessage.error(event.data[2],3)
        break;
      case "warn":
        antmessage.warn(event.data[2],3);
        break;
      default:
        antmessage.info(event.data[2],3)
        break;
    }
  })
  //ws
  // wsClientStore.setupWsSettings().then(()=>{
  //   if(wsClientStore.ws.autoConnect){
  //     try {
  //       wsClientStore.connectWs()
  //     }catch (error){
  //       antmessage.error("ws连接失败：" + error)
  //     }
  //   }
  // })
  // mqttClientStore.setupMqttSettings().then(()=>{
  //     if(mqttClientStore.mqttInfo.autoConnect){
  //       try {
  //         console.log("mqtt连接vue")
  //         mqttClientStore.connectMqtt()
  //       }catch (error){
  //         antmessage.error("mqtt连接失败：" + error)
  //       }
  //     }
  // })
  mqttClientStore.setupMqttSettings()
  loading.value = false
})


</script>

<template>
<!--  <SplashView v-if="loading && route.path !== '/ytdlpWidgets'" />-->
<!--  <template v-else-if="appStore.widgetsEnable"><widgetsView /></template>-->
<!--  <template v-else-if="!appStore.widgetsEnable">-->
  <meta name="referrer" content="no-referrer">
  <template v-if="!route.meta.newlayout">
    <TitleBar />
    <div class="main">
      <NavigationBar />
      <MainPage />
    </div>

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
  <template v-else><WidgetsView /></template>
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
