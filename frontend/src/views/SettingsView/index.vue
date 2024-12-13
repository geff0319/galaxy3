<script setup lang="ts">
import {ref} from 'vue'
import {useI18n} from 'vue-i18n'

import {useAppStore, useYtdlpSettingsStore} from '@/stores'
import KernelView from '@/views/KernelView/index.vue'
import GeneralSetting from './components/GeneralSetting.vue'
import TranslateSetting from './components/TranslateSetting.vue'
import YtdlpSetting from "./components/YtdlpSetting.vue";
import WsSetting from "@/views/SettingsView/components/WsSetting.vue";
import MqttSetting from "@/views/SettingsView/components/MqttSetting.vue";
import icons from '@/components/Icon/icons'


type TabItemType = {
  // icon:(typeof icons)[number]
  icon:string
  key: string
  tab: string
}

const settings:TabItemType[] = [
  { icon:'system',key: 'general', tab: 'settings.general' },
  { icon:'video',key: 'ytdlp', tab: '视频配置' },
  { icon:'server',key: 'mqtt', tab: '通信配置' },
  // { key: 'websocket', tab: 'WebSocket' },
  // { key: 'kernel', tab: 'router.kernel' }
]

const activeKey = ref(settings[0].key)

const { t } = useI18n()
const appStore = useAppStore()

</script>

<template>
  <TabsWithIcon v-model:active-key="activeKey" :items="settings" height="100%">
    <template #general>
      <GeneralSetting />
    </template>

    <template #translate>
      <TranslateSetting />
    </template>

    <template #ytdlp>
      <YtdlpSetting />
    </template>

<!--    <template #websocket>-->
<!--      <WsSetting />-->
<!--    </template>-->
    <template #mqtt>
      <MqttSetting />
    </template>

    <!--    <template #kernel>-->
<!--      <KernelView />-->
<!--    </template>-->

    <template #extra>
      <button class="custom-btn" @click="appStore.showAbout = true">
<!--        <icon icon="about"></icon>-->
        <svg> <use href="#info"></use></svg>
        <span :style="{marginLeft: '8px'}">{{ t('router.about') }}</span>
      </button>
    </template>
  </TabsWithIcon>
</template>

<style lang="less" scoped>
@import "@/assets/main";

.custom-btn {
  width: 80%;
  padding: 8px 12px;
  margin: 3px 4px;
  font-size: 14px;
  border: 2px #ccc; /* 设置边框 */
  cursor: pointer; /* 鼠标悬停时显示指针 */
  transition: all 0.3s ease; /* 为所有变化添加过渡效果 */
  color: var(--btn-text-color);
  background-color: var(--btn-text-bg);
  border: none;
  display: flex;
  justify-content: flex-start; /* 左对齐 */
  align-items: center;
  border-radius: 6px;
  .icon_hover();
  &:hover {
    color: var(--btn-text-hover-color);
    background-color: var(--btn-text-hover-bg);
  }
  &:active {
    color: var(--btn-text-active-color);
    background-color: var(--btn-text-active-bg);
  }
  svg {
    width: 22px;
    height: 22px;
  }
}
button.active-tab {
  font-weight: bold; /* 如果选中则加粗 */
  color: var(--btn-text-active-color);
  background-color: var(--btn-text-active-bg);
}
.custom-btn .span{
  font-family: "幼圆", "Yu Yuan", sans-serif;
}
</style>

