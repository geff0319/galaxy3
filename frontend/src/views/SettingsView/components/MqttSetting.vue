<script setup lang="ts">

import { LoadingOutlined } from '@ant-design/icons-vue';
import { ref } from 'vue';
import {useMqttClientStore} from "@/stores/mqttClient";
import {message} from "ant-design-vue";
const appMqttClient = useMqttClientStore()

const loading = ref(false)
const connectMqtt =()=>{
  if (appMqttClient.mqttInfo.broker === ''||appMqttClient.mqttInfo.clientID === ''){
    message.info("必填参数为空")
    return
  }
  loading.value = true
  appMqttClient.connectMqtt()
  loading.value = false
}

appMqttClient.connectionStatus()
</script>

<template>
  <div class="settings">
    <div class="settings-item">
      <div class="title">MQTT服务端</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div style="display:flex;align-items:center;">Broker<Icon icon="must" /></div>
          <input
              class="custom-input"
              type="text"
              placeholder="broker"
              v-model="appMqttClient.mqttInfo.broker"
              :style="{cursor:appMqttClient.isConnected?'not-allowed':'text','background-color':appMqttClient.isConnected?'rgb(245, 245, 245)':'rgb(255, 255, 255)'}"
          />
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <div style="display:flex;align-items:center;">Port<Icon icon="must" /></div>
          <a-input-number style="width: 30%" v-model:value="appMqttClient.mqttInfo.port" :min="1" :max="65535" :disabled="appMqttClient.isConnected"  />
        </div>
      </a-card>
    </div>

    <div class="settings-item">
      <div class="title">MQTT客户端</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div style="display:flex;align-items:center;">Client ID<Icon icon="must" /></div>
          <input
              class="custom-input"
              type="text"
              placeholder="clientid"
              v-model="appMqttClient.mqttInfo.clientID"
              :style="{cursor:appMqttClient.isConnected?'not-allowed':'text','background-color':appMqttClient.isConnected?'rgb(245, 245, 245)':'rgb(255, 255, 255)'}"
          />
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <div>UserName</div>
          <input
              class="custom-input"
              type="text"
              placeholder="username"
              v-model="appMqttClient.mqttInfo.userName"
              :style="{cursor:appMqttClient.isConnected?'not-allowed':'text','background-color':appMqttClient.isConnected?'rgb(245, 245, 245)':'rgb(255, 255, 255)'}"
          />
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <div>Password</div>
          <input
              class="custom-input"
              type="text"
              placeholder="password"
              v-model="appMqttClient.mqttInfo.password"
              :style="{cursor:appMqttClient.isConnected?'not-allowed':'text','background-color':appMqttClient.isConnected?'rgb(245, 245, 245)':'rgb(255, 255, 255)'}"
          />
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <div>AutoConnect</div>
          <a-switch v-model:checked="appMqttClient.mqttInfo.autoConnect" />
        </div>
        <div class="gray-line"></div>
        <div class="card-item">
          <Icon v-if="!loading && appMqttClient.isConnected" icon="success" />
          <Icon  v-else-if="!loading && !appMqttClient.isConnected" icon="failed" />
          <loading-outlined v-else-if="loading" />
          <Button v-if="!appMqttClient.isConnected" class="connect-button" size="small" @click="connectMqtt">连 接</Button>
          <Button v-else class="connect-button" size="small" @click="appMqttClient.disConnectMqtt">断开连接</Button>
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
      .tips {
        font-weight: normal;
        font-size: 12px;
      }
    }
    .card{
      width: 100%;
      &-item{
        height: 30px;
        margin: 0 5px;
        display: flex;
        justify-content: space-between;
        align-items: center;
      }
    }
  }
}
.gray-line {
  height: 1px; /* 设置线的高度 */
  background-color: #cccccc; /* 设置灰色背景 */
  margin: 10px 0; /* 添加上下间距 */
}
.custom-input {
  //cursor:pointer; /* 禁用点击时的光标样式 */
  //background-color: #dcdcdc;
  width: 30%;
  padding: 7px;
  border: 1px solid  #ccc;
  border-radius: 4px;
  box-sizing: border-box;
  font-family: "幼圆", "Yu Yuan", sans-serif;
  outline: none; /* 去除默认聚焦边框 */
}
:deep(.ant-input-number){
  border: 1px solid  #ccc;
  box-shadow: none;
}
.connect-button{
  width: 80px;
  border-radius: 10px;
  background-color: rgb(210, 210, 210);
  font-size: 12px;
  font-weight:400;
  &:hover{
    background-color: rgb(200, 200, 200);
    color: black; /* 按钮文本颜色 */
  }
}
</style>
