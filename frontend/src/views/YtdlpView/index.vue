<script setup lang="ts">
import { onUnmounted, ref } from 'vue'
import { useI18n, I18nT } from 'vue-i18n'
import { FileTextOutlined,FolderOpenOutlined,LinkOutlined,DeleteOutlined,PlayCircleOutlined } from '@ant-design/icons-vue';
import { View } from '@/constant'
import {
  formatResolution, formatSpeedMiB, type Menu, type PluginType, type ProcessType,
  useAppSettingsStore,
  useYtdlpStore
} from '@/stores'
import {setIntervalImmediately} from "@/utils";
import YtdlpForm from "@/views/YtdlpView/components/YtdlpForm.vue";
import {AppCheckBiliLogin, deleteProcess,Browser} from "@/bridge";
import {message} from "ant-design-vue";


const { t } = useI18n()
const appSettingsStore = useAppSettingsStore()
const ytdlpStore = useYtdlpStore()
const showForm = ref(false)
const loginLoad =ref(false)
ytdlpStore.getAllVideoInfo()
const timer = setIntervalImmediately(ytdlpStore.getAllVideoInfo, 1000)

const handleAddSub = async () => {
  showForm.value = true
}
const check = async () =>{
  loginLoad.value = true
  try {
    const data = await AppCheckBiliLogin()
    message.info(data)
    loginLoad.value = false
  }catch (err:any){
    message.error(err)
    loginLoad.value = false
  }

}
onUnmounted(() => {
  clearInterval(timer)
})

const menuList: Menu[] = [
  {
    label: '删除',
    handler: async (id: string) => {
      try {
        const data = await deleteProcess(id)
        await  ytdlpStore.getAllVideoInfo()
        message.info(data)
      }catch (error :any) {
        message.error(error)
      }
    }
  }
]
const generateMenus = (p: ProcessType) => {
  let builtInMenus: Menu[] = menuList.map((v) => ({ ...v, handler: () => v.handler?.(p.id) }))
  if (p.progress.process_status == 3) {
    builtInMenus.unshift({
      label: '重试',
      handler: async () => {
        try {
          ytdlpStore.resProcess = p
          await ytdlpStore.downloadYoutube(true,true)
          // message.info(p.id + "开始下载")
        }catch (error :any) {
          message.error(error)
        }
      }
    })
  }
  return builtInMenus
}
const reTry = async (p:ProcessType) => {
  try {
    ytdlpStore.resProcess = p
    await ytdlpStore.downloadYoutube(true,true)
    message.info(p.id + "开始下载")
  }catch (error :any) {
    message.error(error)
  }
}
const deleteVideo = async (id: string) => {
  try {
    const data = await deleteProcess(id)
    await  ytdlpStore.getAllVideoInfo()
    message.info(data)
  }catch (error :any) {
    message.error(error)
  }
}

</script>

<template>
  <div v-if="ytdlpStore.process.length === 0" class="grid-list-empty">
    <Empty>
      <template #description>
        <I18nT keypath="ytdlp.empty" tag="p" scope="global">
          <template #action>
            <Button @click="handleAddSub" type="link">{{ t('common.add') }}</Button>
          </template>
        </I18nT>
      </template>
    </Empty>
  </div>

  <div v-else class="grid-list-header">
<!--    <Radio-->
<!--        v-model="appSettingsStore.app.ytdlpView"-->
<!--        :options="[-->
<!--        { label: 'common.grid', value: View.Grid },-->
<!--        { label: 'common.list', value: View.List }-->
<!--      ]"-->
<!--        class="mr-auto"-->
<!--    />-->

    <div class="mr-auto">
<!--      <a-card style="width: 100px" size="small">-->
<!--        <a-avatar>-->
<!--          <template #icon><UserOutlined /></template>-->
<!--        </a-avatar>-->
<!--      </a-card>-->
      <Button :disable="loginLoad" @click="check" type="primary">
        校验登录状态
      </Button>
    </div>
    <Button @click="handleAddSub" type="primary">
      {{ t('common.add') }}
    </Button>
  </div>

  <div :class="'grid-list-' + appSettingsStore.app.ytdlpView">
    <Card v-if="appSettingsStore.app.ytdlpView === View.Grid" hoverable size="small" v-for="(p, index) in ytdlpStore.process" :class="'item'" :body-style="{height:'66px'}" :key="p.id + p.progress.process_status" v-menu="generateMenus(p)" >
      <template v-if="appSettingsStore.app.ytdlpView === View.Grid" #cover>
<!--        <a-image-->
<!--            v-if = "p.info.thumbnail.length !== 0 "-->
<!--            class="unselectable-image"-->
<!--            :height="125"-->
<!--            :preview="false"-->
<!--            :src="p.info.thumbnail"-->
<!--            fallback="data:image/png;base64,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"-->
<!--        />-->
<!--        <a-image-->
<!--            v-else-->
<!--            class="unselectable-image"-->
<!--            :height="125"-->
<!--            :preview="false"-->
<!--            src="data:image/png;base64,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"-->
<!--        />-->
            <a-image
                class="unselectable-image"
                :height="125"
                :preview="false"
                src="data:image/png;base64,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"
            />
      </template>
      <div class="overlay" v-show="p.progress.process_status === 3">
        <!-- 纯圆形遮罩 -->
        <div :class='"circle-"+ appSettingsStore.app.ytdlpView'><Icon class="icon" icon="fail"/></div>
      </div>
      <a-flex vertical gap="small" >
        <a-row :gutter="16">
          <a-col :span="18" >
            <div :class="'title-' + appSettingsStore.app.ytdlpView" :title="p.info.title">{{ p.info.title }}</div>
          </a-col>
          <a-col :span="6" >
            <a-tag style="height: 20px" color="processing">{{  p.biliMeta.SelectedVideoQuality || formatResolution(p.info.resolution) }}</a-tag>
          </a-col>
        </a-row>
        <a-row :gutter="4">
          <a-col :span="18">
            <a-progress style="display: flex;align-items: center" :percent="p.progress.process_status===2 ? 100: parseInt(p.progress.percentage,10)" size="small" :show-info="false" />
          </a-col>
          <a-col :span="6" style="display: flex;align-items: center">
            <div>{{ formatSpeedMiB(p.progress.speed) }}</div>
          </a-col>
        </a-row>
      </a-flex>
    </Card>
    <a-card v-if="appSettingsStore.app.ytdlpView === View.List" v-for="(p, index) in ytdlpStore.process" size="small" class="a-card">
      <!--      View.List-->
      <div :class="'card-' + appSettingsStore.app.ytdlpView">
        <div style="width: 16%;border: 1px solid  rgba(232,232,232,0.99);border-radius: 8px;margin-right: 15px">
          <img
              class="unselectable-image"
              height="100%"
              width="100%"
              :src="p.info.thumbnail"
              alt=""
          />
        </div>
        <div style="width: 84%">
          <div style="height: 80%">
            <div style="height: 40%; display: flex; align-items: center;">
              <div style="width: 600px;display: flex; ">
                <div :class="'title-' + appSettingsStore.app.ytdlpView" :title="p.info.title">
                  {{ p.info.title }}
                </div>
                <a-tag color="#108ee9">{{ p.biliMeta.SelectedVideoQuality || formatResolution(p.info.resolution) }}</a-tag>
              </div>
              <div style="flex: 1;display: flex;align-items: center;">
                <div style="width: 73px">
                  {{ p.info.created_at.substring(0, 10) }}
                </div>
                <div class="icon-container">
                  <a-tooltip color="#ffffff">
                    <template #title>
                      <span class="custom-tooltip">打开链接</span>
                    </template>
                    <LinkOutlined class="icon-wrapper" @click="Browser.OpenURL(p.url)"/>
                  </a-tooltip>
                  <a-tooltip color="#ffffff">
                    <template #title>
                      <span class="custom-tooltip">播放视频</span>
                    </template>
                    <PlayCircleOutlined class="icon-wrapper" @click="Browser.OpenURL(p.output.savedFilePath)"/>
                  </a-tooltip>
                  <a-tooltip color="#ffffff">
                    <template #title>
                      <span class="custom-tooltip">打开文件目录</span>
                    </template>
                    <FolderOpenOutlined class="icon-wrapper" @click="Browser.OpenURL(p.output.Path)"/>
                  </a-tooltip>
                  <a-tooltip color="#ffffff">
                    <template #title>
                      <span class="custom-tooltip">删除记录</span>
                    </template>
                    <DeleteOutlined class="icon-wrapper" @click="deleteVideo(p.id)"/>
                  </a-tooltip>
                </div>
              </div>
            </div>
          </div>
          <div style="display: flex">
            <a-progress
                style="width:92%"
                :status="p.progress.process_status===4 ? 'active': undefined"
                :percent="p.progress.process_status===2 ? 100: parseInt(p.progress.percentage,10)"
                size="small"
                :show-info="false"
            />
            <div>{{ formatSpeedMiB(p.progress.speed) }}</div>
          </div>
        </div>
      </div>
    </a-card>

  </div>

  <Modal v-model:open="showForm" max-height="80" :footer="false" >
    <YtdlpForm />
  </Modal>
</template>

<style lang="less" scoped>
.custom-tooltip {
  color: black; /* 自定义字体颜色 */
  font-size: 12px; /* 自定义字体大小 */
}
.icon-container {
  margin-left: 20px;
  display: flex; /* 横向排列图标 */
  gap: 10px; /* 图标之间的间距 */
}

.icon-wrapper {
  font-size:16px;
  display: inline-block; /* 确保图标适合内容 */
  cursor: pointer; /* 鼠标悬浮时的光标样式 */
  &:hover {
    color: #008316;
  }

}
.title-grid {
  display: block;
  align-items: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: bold;
}
.title-list {
  //flex: 1; /* 让左侧内容占据剩余的所有空间 */
  //flex: 0 0 70%;
  margin-right: 8px;
  max-width: 80%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: bold;
}
.a-card{
  margin: 5px 0;
  background-color: rgb(243,243,243);
  border: 1px solid  rgba(232,232,232,0.99);
  transition:
      box-shadow 0.4s,
      background-color 0.4s;
  &:hover {
    background-color: var(--card-hover-bg);
    box-shadow: 0 8px 8px rgba(0, 0, 0, 0.06);
  }
}
.card-list {
  width: 100%; /* 子元素宽度为父元素的 100% */
  //background-color: lightblue; /* 方便观察 */
  justify-content: center; /* 水平居中 */
  //padding: 7px 0;
  height: 80px;
  display: flex;
}
.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
.circle-grid {
  width: 130px; /* 圆的直径 */
  height: 130px; /* 圆的直径 */
  border-radius: 50%; /* 将正方形遮罩层变为圆形 */
  background-color: rgba(0, 0, 0, 0.7); /* 半透明黑色背景 */
}
.circle-list {
  //width: 60px; /* 圆的直径 */
  //height: 60px; /* 圆的直径 */
  //border-radius: 50%; /* 将正方形遮罩层变为圆形 */
  //background-color: rgba(0, 0, 0, 0.7); /* 半透明黑色背景 */
  width: 100%;
  height: 110%;
  background-color: rgba(0, 0, 0, 0.7); /* 半透明黑色背景 */
  border-radius: 8px; /* 圆形边框 */
  z-index: 1; /* 确保遮光罩在内容之上 */
}
.icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 48px;
}
.unselectable-image {
  border-radius: 8px; /* 设置圆角 */
}
</style>
