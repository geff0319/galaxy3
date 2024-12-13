<script setup lang="ts">
import {ref} from "vue";
import { OpenDirectoryDialog,Browser} from '@/bridge'
import {useYtdlpSettingsStore} from "@/stores";
import {message} from "ant-design-vue";
import YtDlp from "@/views/KernelView/YtDlp.vue";
import BiliQrCode from "@/views/SettingsView/components/BiliQrCode.vue";

const showForm = ref(false)
const ytdlpSettingsStore = useYtdlpSettingsStore()
const handelOpenFileDialog =async ()=>{
  try {
    const folder = await  OpenDirectoryDialog()
    if (folder.length!==0) {
      ytdlpSettingsStore.ytdlpConfig.downloadPath = folder
    }
  }catch (error :any){
    message.error("选择文件夹失败")
  }
}
const biliQrOpen = ()=>{
  showForm.value = true
}
const handelOpenDirDialog =()=>{
  Browser.OpenURL(ytdlpSettingsStore.ytdlpConfig.downloadPath)
}
</script>

<template>
  <div class="settings">
    <div class="settings-item">
      <div class="title">配置项</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div>下载目录</div>
          <a-tooltip :title="ytdlpSettingsStore.ytdlpConfig.downloadPath">
            <a-input
                style="width: 40%"
                v-model:value="ytdlpSettingsStore.ytdlpConfig.downloadPath"
                placeholder="请选择视频下载目录"
                disabled
            >
              <template #addonAfter>
                <div style="width: 60px"></div>
                <div class="addon-content" @click="handelOpenFileDialog" >
                  <Icon icon="edit"/>
                </div>
                <div class="addon-content1" @click="handelOpenDirDialog" >
                  <Icon icon="folder"/>
                </div>
              </template>
            </a-input>
          </a-tooltip>
        </div>
        <div class="gray-line" />
        <div class="card-item">
          <div>并发数量(重启生效)</div>
          <a-input-number style="width: 20%" v-model:value="ytdlpSettingsStore.ytdlpConfig.queueSize" :min="1" :max="10" />
        </div>
      </a-card>
    </div>

    <div class="settings-item">
      <div class="title">COOKIE管理</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div>BILIBILI</div>
          <button class="button" @click="biliQrOpen">
            <Icon icon="qr"/>
          </button>
        </div>
      </a-card>
    </div>
    <YtDlp />
  </div>

  <Modal v-model:open="showForm" max-height="80" :footer="false">
    <BiliQrCode />
  </Modal>
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
      }
    }
  }
}
.input-wrapper{
  display: flex;
  align-items: center; /* 垂直居中 */
  .label {
    flex: 0 0 auto; /* 不随内容变化而变化 */
    margin-right: 10px; /* 右侧留出一些空间 */
  }
  .input-container{
    display: flex;
    flex: 1/* 填充剩余空间 */
  }
}
.gray-line {
  height: 1px; /* 设置线的高度 */
  background-color: #cccccc; /* 设置灰色背景 */
  margin: 10px 0; /* 添加上下间距 */
}
.addon-content {
  position: absolute;
  top: 0;
  left: 0;
  width: 50%;
  height: 100%;
  cursor:pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  &:hover{
    background-color: rgb(220, 220, 220);
    color: black; /* 按钮文本颜色 */
  }
  //background-color: #389e0d;
}
.addon-content1 {
  position: absolute;
  top: 0;
  right: 0;
  width: 50%;
  height: 100%;
  border-left: 1px solid  #ccc;
  cursor:pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  &:hover{
    background-color: rgb(220, 220, 220);
    color: black; /* 按钮文本颜色 */
  }
  //background-color: #389e0d;
}
:deep(.ant-input-number){
  border: 1px solid  #ccc;
  box-shadow: none;
}
:deep(.ant-input) {
  border: 1px solid  #ccc;
  box-shadow: none;
}
.button {
  cursor:pointer;
  border: 1px solid  #ccc;
  border-radius: 4px;
  box-sizing: border-box;
  padding: 6px 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgb(250, 250, 250);
  &:hover{
    background-color: rgb(220, 220, 220);
    color: black; /* 按钮文本颜色 */
  }
}
</style>
