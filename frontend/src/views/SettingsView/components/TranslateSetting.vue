<script setup lang="ts">
import { ref ,reactive} from 'vue'
import { useI18n } from 'vue-i18n'

import { useMessage } from '@/hooks'
import { useAppSettingsStore, useEnvStore } from '@/stores'


const { t } = useI18n()
const { message } = useMessage()
const appSettings = useAppSettingsStore()
const envStore = useEnvStore()
const hidestr = ref('***********************')
const showOptions =
    reactive<Record<string, boolean>>({
      "tencentTanslateSecretId": false,
      "tencentTanslateSecretKey": false
    })

const changeShow = (key:string) => {
  showOptions[key] = !showOptions[key];
  console.log( showOptions[key] )
}

</script>

<template>
  <div class="settings">
    <div class="settings-item">
      <div class="title">{{ t('settings.translate.tencentTanslate') }}</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div>{{ t('settings.translate.tencentTanslateSecretId') }}</div>
          <a-input-password
              class="password"
              v-model:value="appSettings.app.translate.tencentTanslateSecretId"
              placeholder="请输入SecretId"
              :visible="false"
          ></a-input-password>
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <div>{{ t('settings.translate.tencentTanslateSecretKey') }}</div>
          <a-input-password
              class="password"
              v-model:value="appSettings.app.translate.tencentTanslateSecretKey"
              placeholder="请输入SecretKey"
              :visible="false"
          ></a-input-password>
        </div>
      </a-card>
    </div>
  </div>
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
.password{
  width: 40%;
}
:deep(.ant-input-affix-wrapper:hover) {
  border: 1px solid  #ccc;
  box-shadow: none;
}
:deep(.ant-input-affix-wrapper) {
  border: 1px solid  #ccc;
  box-shadow: none;
}
</style>
