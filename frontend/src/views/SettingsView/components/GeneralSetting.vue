<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useMessage } from '@/hooks'
import { useAppSettingsStore, useEnvStore } from '@/stores'
import { Theme, Lang, WindowStartState, Color } from '@/constant'
import { APP_TITLE, APP_VERSION, getTaskSchXmlString } from '@/utils'
import {GetEnv, Writefile, Removefile, Browser, OpenDirectoryDialog} from '@/bridge'
import {
  QuerySchTask,
  CreateSchTask,
  DeleteSchTask,
  CheckPermissions,
  SwitchPermissions
} from '@/utils'
import {AppChangeLog} from "@/bridge/utils";
import TranslateSetting from "@/views/SettingsView/components/TranslateSetting.vue";

const isAdmin = ref(false)
const isTaskScheduled = ref(false)

const { t } = useI18n()
const { message } = useMessage()
const appSettings = useAppSettingsStore()
const envStore = useEnvStore()


const themes = [
  {
    label: 'settings.theme.dark',
    value: Theme.Dark
  },
  {
    label: 'settings.theme.light',
    value: Theme.Light
  },
  {
    label: 'settings.theme.auto',
    value: Theme.Auto
  }
]

const colors = [
  {
    label: 'settings.color.default',
    value: Color.Default
  },
  {
    label: 'settings.color.orange',
    value: Color.Orange
  },
  {
    label: 'settings.color.pink',
    value: Color.Pink
  },
  {
    label: 'settings.color.red',
    value: Color.Red
  },
  {
    label: 'settings.color.skyblue',
    value: Color.Skyblue
  },
  {
    label: 'settings.color.green',
    value: Color.Green
  }
]

const langs = [
  {
    label: 'settings.lang.zh',
    value: Lang.ZH
  },
  {
    label: 'settings.lang.en',
    value: Lang.EN
  }
]

const windowStates = [
  { label: 'settings.windowState.normal', value: WindowStartState.Normal },
  // { label: 'settings.windowState.maximised', value: WindowStartState.Maximised },
  { label: 'settings.windowState.minimised', value: WindowStartState.Minimised }
  // { label: 'settings.windowState.fullscreen', value: WindowStartState.Fullscreen }
]

const resetFontFamily = () => {
  appSettings.app['font-family'] = '"Microsoft Yahei", "Arial", sans-serif, "Twemoji Mozilla"'
}

const resetUserAgent = () => {
  appSettings.app.userAgent = APP_TITLE + '/' + APP_VERSION
}

const onPermChange = async (v: boolean) => {
  try {
    await SwitchPermissions(v)
    message.success('success')
  } catch (error: any) {
    message.error(error)
    console.log(error)
  }
}

const handleOpenFolder = async () => {
  const { basePath } = await GetEnv()
  await Browser.OpenURL(basePath)
}
const handleOpenLogFolder = async () => {
  const { basePath } = await GetEnv()
  if(appSettings.app.logPath === ''){
    await Browser.OpenURL(basePath+ "/logs")
  }else {
    await Browser.OpenURL(appSettings.app.logPath)
  }
}
const handelSelectLogFolderDialog =async ()=>{
  try {
    const folder = await  OpenDirectoryDialog()
    if (folder.length!==0) {
      appSettings.app.logPath = folder
      await AppChangeLog(0,folder)
    }
  }catch (error :any){
    message.error("选择文件夹失败")
  }
}
const checkSchtask = async () => {
  try {
    await QuerySchTask(APP_TITLE)
    isTaskScheduled.value = true
  } catch (error) {
    isTaskScheduled.value = false
  }
}

const onTaskSchChange = async (v: boolean) => {
  isTaskScheduled.value = !v
  try {
    if (v) {
      await createSchTask(appSettings.app.startupDelay)
    } else {
      await DeleteSchTask(APP_TITLE)
    }
    isTaskScheduled.value = v
  } catch (error: any) {
    console.error(error)
    message.error(error)
  }
}

const onStartupDelayChange = async (delay: number) => {
  try {
    await createSchTask(delay)
  } catch (error: any) {
    console.error(error)
    message.error(error)
  }
}

const createSchTask = async (delay = 30) => {
  const xmlPath = 'data/.cache/tasksch.xml'
  const xmlContent = await getTaskSchXmlString(delay)
  await Writefile(xmlPath, xmlContent)
  await CreateSchTask(APP_TITLE, xmlPath)
  await Removefile(xmlPath)
}

const onThemeClick = (e: MouseEvent) => {
  document.documentElement.style.setProperty('--x', e.clientX + 'px')
  document.documentElement.style.setProperty('--y', e.clientY + 'px')
}

if (envStore.env.os === 'windows') {
  checkSchtask()

  CheckPermissions().then((admin) => {
    isAdmin.value = admin
  })
}
</script>

<template>
  <div class="settings">
<!--    <div class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.theme.name') }}-->
<!--      </div>-->
<!--      <Radio v-model="appSettings.app.theme" @click="onThemeClick" :options="themes" />-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.color.name') }}-->
<!--      </div>-->
<!--      <Radio v-model="appSettings.app.color" :options="colors" />-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">{{ t('settings.lang.name') }}</div>-->
<!--      <Radio v-model="appSettings.app.lang" :options="langs" />-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">{{ t('settings.fontFamily') }}</div>-->
<!--      <div style="display: flex; align-items: center">-->
<!--        <Button @click="resetFontFamily" v-tips="'settings.resetFont'" type="text">-->
<!--          <Icon icon="reset" />-->
<!--        </Button>-->
<!--        <Input v-model="appSettings.app['font-family']" editable style="margin-left: 8px" />-->
<!--      </div>-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">{{ t('settings.appFolder.name') }}</div>-->
<!--      <Button @click="handleOpenFolder" type="primary">-->
<!--        <FolderOpenOutlined />-->
<!--        <span style="margin-left: 8px">{{ t('settings.appFolder.open') }}</span>-->
<!--      </Button>-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">日志</div>-->
<!--      <Button @click="handleOpenLogFolder" type="primary">-->
<!--        <FolderOpenOutlined />-->
<!--        <span style="margin-left: 8px">打开日志目录</span>-->
<!--      </Button>-->
<!--      <Button @click="handelSelectLogFolderDialog" type="primary">-->
<!--        <FolderOpenOutlined />-->
<!--        <span style="margin-left: 8px">设置目录路径</span>-->
<!--      </Button>-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.exitOnClose') }}-->
<!--      </div>-->
<!--      <Switch v-model="appSettings.app.exitOnClose" />-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.closeKernelOnExit') }}-->
<!--      </div>-->
<!--      <Switch v-model="appSettings.app.closeKernelOnExit" />-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.autoSetSystemProxy') }}-->
<!--      </div>-->
<!--      <Switch v-model="appSettings.app.autoSetSystemProxy" />-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.autoStartKernel') }}-->
<!--      </div>-->
<!--      <Switch v-model="appSettings.app.autoStartKernel" />-->
<!--    </div>-->
<!--    <div v-if="envStore.env.os === 'windows'" class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.admin') }}-->
<!--        <span class="tips">({{ t('settings.needRestart') }})</span>-->
<!--      </div>-->
<!--      <Switch v-model="isAdmin" @change="onPermChange" />-->
<!--    </div>-->
<!--    <div v-if="envStore.env.os === 'windows'" class="settings-item">-->
<!--      <div class="title">-->
<!--        {{ t('settings.startup.name') }}-->
<!--        <span class="tips">({{ t('settings.needAdmin') }})</span>-->
<!--      </div>-->
<!--      <div style="display: flex; align-items: center">-->
<!--        <Switch v-model="isTaskScheduled" @change="onTaskSchChange" style="margin-right: 16px" />-->
<!--        <template v-if="isTaskScheduled">-->
<!--          <Radio v-model="appSettings.app.windowStartState" :options="windowStates" type="number" />-->
<!--          <span style="margin: 0 8px">{{ t('settings.startup.delay') }}</span>-->
<!--          <Input-->
<!--            v-model="appSettings.app.startupDelay"-->
<!--            @submit="onStartupDelayChange"-->
<!--            :min="0"-->
<!--            type="number"-->
<!--          />-->
<!--        </template>-->
<!--      </div>-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">{{ t('settings.addToMenu') }}</div>-->
<!--      <Switch v-model="appSettings.app.addPluginToMenu" />-->
<!--    </div>-->
<!--    <div class="settings-item">-->
<!--      <div class="title">{{ t('settings.userAgent.name') }}</div>-->
<!--      <div style="display: flex; align-items: center">-->
<!--        <Button @click="resetUserAgent" v-tips="'settings.userAgent.reset'" type="text">-->
<!--          <Icon icon="reset" />-->
<!--        </Button>-->
<!--        <Input v-model.lazy="appSettings.app.userAgent" editable style="margin-left: 8px" />-->
<!--      </div>-->
<!--    </div>-->
  </div>
  <div class="settings">
    <div class="settings-item">
      <div class="title">系统</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div>{{ t('settings.appFolder.name') }}</div>
          <button class="button" @click="handleOpenFolder">
            <Icon icon="folder"/>
          </button>
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <div>{{ t('settings.exitOnClose') }}</div>
          <a-switch size="large" v-model:checked="appSettings.app.exitOnClose" />
        </div>
        <div class="gray-line" v-if="envStore.env.os === 'windows'"></div>
        <div class="card-item" v-if="envStore.env.os === 'windows'">
          <div>{{ t('settings.admin') }}<span>({{ t('settings.needAdmin') }})</span></div>
          <a-switch size="large" v-model:checked="isAdmin" @change="onPermChange" />
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <div>日志</div>
          <div style="display: flex">
            <button class="button" @click="handelSelectLogFolderDialog">
              <Icon icon="edit"/>
            </button>
            <button class="button" @click="handleOpenLogFolder">
              <Icon icon="folder"/>
            </button>
          </div>
        </div>
      </a-card>
    </div>
  </div>
  <translate-setting></translate-setting>
</template>

<style lang="less" scoped>
.settings {
  &-item {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 8px 16px;
    .title {
      align-self: flex-start;
      font-size: 18px;
      font-weight: bold;
      padding: 8px 0 16px 0;
      margin-bottom: 10px;
      .tips {
        font-weight: normal;
        font-size: 12px;
      }
    }
    .card{
      width: 100%;
      box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
      transform: translateY(-10px); /* 让元素上移 */
      &-item{
        height: 30px;
        margin: 0 5px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-family: "幼圆", "Yu Yuan", sans-serif;
      }
    }
  }
}
.gray-line {
  height: 1px; /* 设置线的高度 */
  background-color: #cccccc; /* 设置灰色背景 */
  margin: 10px 0; /* 添加上下间距 */
}
.button {
  cursor:pointer;
  border: 1px solid  #ccc;
  border-radius: 4px;
  box-sizing: border-box;
  padding: 6px 12px;
  margin-left: 20px;
  display: flex;
  align-content: center;
  justify-items: center;
  background-color: rgb(250, 250, 250);
  &:hover{
    background-color: rgb(220, 220, 220);
    color: black; /* 按钮文本颜色 */
  }
}
</style>
