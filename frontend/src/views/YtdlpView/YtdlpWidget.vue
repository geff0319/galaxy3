<script setup lang="ts">
import {formattedFileSize, useYtdlpStore} from '@/stores'
import {onBeforeUnmount, onMounted, ref} from 'vue';
import {Clipboard, Events, HttpHead, Window} from "@/bridge";
import {message} from "ant-design-vue";
import {AppGetBelowWinPos} from "@/bridge/utils";


const ytdlpStore = useYtdlpStore()
const selectedOption = ref('');
const definitionOptions = ref<{ label: string; value: string }[]>([{ value:"",label:""}])
const fileSize = ref('')
const disableDownload = ref(false)

const winId = ref('')

onMounted(async ()=>{
  winId.value = await Window.Name()
  Events.On('get'+ winId.value + 'option',(event:any)=>{
    selectedOption.value = event.data[0]
    const [selectQuality, selectCodecs] = event.data[0].split('|')
    ytdlpStore.resProcess.biliMeta.SelectedVideoQuality = selectQuality
    handleChange(event.data[1])
  })
})
onBeforeUnmount(async ()=>{
  Events.Off('get'+ winId.value + 'option')
})

const closeYtWidget = async ()=> {
  await Window.Get(winId.value + 'option').Close()
  await Window.Close()
}
const hideYtWidget = async ()=> {
  await Window.Minimise()
}
const optionsShow = async ()=>{
  const win = Window.Get(winId.value + 'option')
  Events.Emit({name:winId.value + 'option',data:[definitionOptions.value]})
  const {dstX,dstY} = await AppGetBelowWinPos(await win.Width(),await win.Height())
  console.log(dstX + '-----' + dstY)
  await win.SetRelativePosition(dstX,dstY)
  await win.SetAlwaysOnTop(true)
  await win.Show()
  await win.Focus()
}
const init =  async () => {
  try {
    const clipboardUrl = await Clipboard.Text()
    ytdlpStore.videoUrl = clipboardUrl
    ytdlpStore.determineUrl(clipboardUrl)
    if(ytdlpStore.downloadUrl !=='' ) {
      await ytdlpStore.getVideoMeta()
    }
    generateOptions()
  }catch (error:any){
    ytdlpStore.resProcess.output.Path = ''
    ytdlpStore.resProcess.info.filesize_approx = 0
    ytdlpStore.resProcess.info.filename = ''
    disableDownload.value = true
    message.error('解析链接出错：'+ error,3)
  }
}
init()
const filterCodecString = (codecString:string) =>{
  // 使用正则表达式匹配以字母开头，后面跟随数字和点的模式
  const match = codecString.match(/^([a-zA-Z0-9]+)\./);
  return match ? match[1] : codecString;
}
const generateOptions= ()=>{

  console.log(ytdlpStore.formatResolution(ytdlpStore.resProcess.info.resolution))
  const options: { value: string; label: string }[] = [];
  if(ytdlpStore.videoType==='BILIBILI'){
    ytdlpStore.resProcess.biliMeta.Vir.data.support_formats.flatMap(item => {
      item.codecs.map(codec => (
          options.push({ value: item.quality + '-' + codec, label: item.new_description + '|' + filterCodecString(codec)})
      ))
    })
  }else if (ytdlpStore.videoType==='YOUTUBE') {
    options.push({ value: '', label: ytdlpStore.formatResolution(ytdlpStore.resProcess.info.resolution) })
  }else {
    options.push({ value: '', label: '' })
  }
  // return options
  definitionOptions.value = options
  selectedOption.value = options[0].label
  handleChange(options[0].value)
  const [selectQuality, selectCodecs] = options[0].label.split('|')
  ytdlpStore.resProcess.biliMeta.SelectedVideoQuality = selectQuality
  ytdlpStore.resProcess.biliMeta.SelectedVideoCodecs = selectCodecs
  // console.log(winId.value + 'option')
  // Events.Emit({name:winId.value + 'option',data:[options]})
}
const download = async () =>{
  try {
    await ytdlpStore.downloadYoutube(true,false)
    Events.Emit({name:'notify', data:[true,"info", "【" + ytdlpStore.resProcess.info.filename + "】开始下载"]})
    await closeYtWidget()
  }catch (error:any){
    Events.Emit({name:'notify', data:[true,"error", "【" + ytdlpStore.resProcess.info.filename + "】下载失败：" + error]})
    await closeYtWidget()
  }
}
const getClipboard =async ()=>{
  ytdlpStore.videoUrl = await Clipboard.Text()
}
const refresh= async ()=>{
  try {
    ytdlpStore.determineUrl(ytdlpStore.videoUrl)
    if(ytdlpStore.downloadUrl !=='' ) {
      await ytdlpStore.getVideoMeta()
    }
    generateOptions()
    disableDownload.value = false
  }catch (error:any){
    ytdlpStore.resProcess.output.Path = ''
    ytdlpStore.resProcess.info.filesize_approx = 0
    ytdlpStore.resProcess.info.filename = ''
    ytdlpStore.downloadUrl = ''
    ytdlpStore.videoType = ''
    disableDownload.value = true
    message.error('解析链接出错：'+ error,3)
  }
}
const handleChange = async (val: string)=>{
  if(ytdlpStore.videoType==='BILIBILI'){
    const [selectQuality, selectCodecs] = val.split('-')
    const video = ytdlpStore.resProcess.biliMeta.Vir.data.dash.video.find(video => {
      return video.id == selectQuality && video.codecs === selectCodecs
    } );
    if(video){
      const h = {'Referer':'https://www.bilibili.com/'}
      try{
        const { header:header1 } = await HttpHead(video.backupUrl[0],h)
        fileSize.value = ytdlpStore.formattedFileSize( Number(header1['Content-Length']))
        ytdlpStore.resProcess.info.resolution = video.height.toString()
        ytdlpStore.resProcess.biliMeta.SelectedVideoStreamUrl = video.base_url
        disableDownload.value = false
      }catch (err){
        fileSize.value = ytdlpStore.formattedFileSize( 0)
        disableDownload.value = true
      }
    }else {
      fileSize.value = ytdlpStore.formattedFileSize( -1)
      disableDownload.value = true
    }
  }else{
    fileSize.value = ytdlpStore.formattedFileSize(ytdlpStore.resProcess.info.filesize_approx)
  }
}
</script>

<template>
  <div class="ytcontainer">
    <div style="--wails-draggable: drag;height: 40px" class="titlebar">
      <img class="logo" draggable="false" src="@/assets/logo.png" />
      <div class="appname">下载任务</div>
      <div class="menus"></div>
      <div class="action">
        <Button @click.stop="hideYtWidget" type="text">
          <Icon icon="minimize" />
        </Button>
        <Button
            @click.stop="closeYtWidget"
            :class="{ 'hover-red': true }"
            type="text"
        >
          <Icon icon="close" />
        </Button>
      </div>
    </div>
    <a-flex gap="small" class="spaced-div" vertical>
      <div style="position: relative;">
        <input
            type="text"
            class="custom-input"
            v-model="ytdlpStore.videoUrl"
            placeholder="下载链接"
        />
        <Icon class="div-icon" icon="clipboard" @click="getClipboard"/>
      </div>
      <div style="display:flex; gap: 4px">
        <div  style="flex: 1;position: relative;">
          <input
              class="custom-input"
              type="text"
              v-model="ytdlpStore.resProcess.output.Path"
              placeholder="保存路径"
              style="cursor:default"
              :title="ytdlpStore.resProcess.output.Path"
              readonly
          />
          <Icon class="div-icon" icon="folder"/>
        </div>
<!--        <custom-select style="width: 37%"  @change="handleChange" :options="definitionOptions" :maxHeight="97" :disable="!(ytdlpStore.videoType==='BILIBILI')"/>-->
        <div :style="{'width':'37%',position: 'relative', display: 'inline-block','pointer-events': ytdlpStore.videoType==='BILIBILI' ? 'auto':'none'}" @click.stop="optionsShow" >
          <input
              style="cursor:pointer"
              type="text"
              class="custom-input"
              :value="selectedOption"
              placeholder="请选择"
              readonly
          />
          <Icon class="div-icon" icon="foldDown"/>
        </div>

      </div>
      <div style="display:flex;align-items: center; gap: 4px">
        <input
            type="text"
            style="flex: 1;cursor:default"
            v-model="ytdlpStore.resProcess.info.filename"
            :title="ytdlpStore.resProcess.info.filename"
            class="custom-input"
            placeholder="文件名"
            readonly
        />
        <div style="width:37%;display:flex;align-items: center; gap: 4px">
          <div v-if="ytdlpStore.parseing" class="load">
            <Icon icon="load" />
          </div>
          <Icon v-if="!ytdlpStore.parseing" style="margin-left: 5px" icon="file" />
          <div v-if="!ytdlpStore.parseing && !(ytdlpStore.downloadUrl.length===0)" style="font-size: 12px;font-weight:400;">{{ fileSize }}</div>
          <Icon v-if="!ytdlpStore.parseing && ytdlpStore.downloadUrl.length===0" icon="question" />
          <div style="flex: 1;"></div>
          <div class="refresh" :style="{cursor:ytdlpStore.parseing?'not-allowed':'default'}">
            <Icon icon="refresh" @click="refresh"/>
          </div>
        </div>
      </div>
      <div style="display: flex;justify-content: space-between;margin-top: 10px">
        <Button class="cancel-button" size="large" @click="closeYtWidget">取  消</Button>
        <Button class="check-button" :style="{'pointer-events': ytdlpStore.downloadUrl.length===0 || disableDownload ?'none':'auto'}" @click="download" size="large">下  载</Button>
      </div>
    </a-flex>
  </div>

</template>

<style lang="less" scoped>
.div-icon{
  padding: 0 10px;
  height: 95%;
  position: absolute;
  right: 1px;
  top: 50%;
  transform: translateY(-50%);
  border-radius: 0 4px 4px 0; /* 圆角 */
  cursor: pointer;
  background-color: #dcdcdc;
  transition: background-color 0.3s; /* 背景颜色过渡效果 */
  &:hover {
    background-color: #c8c8c8; // 悬浮时改变背景颜色
  }
}
.option{
  height: 100px;
  line-height: 20px;
  overflow-y: auto;
  display: block;
}
.load{
  margin-left: 10px;
  display: flex; /* 使用 flex 布局 */
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  animation: spin 1s linear infinite; /* 应用旋转动画 */
}
@keyframes spin {
  0% {
    transform: rotate(0deg); /* 起始角度 */
  }
  100% {
    transform: rotate(360deg); /* 结束角度 */
  }
}
.refresh{
  padding: 6px;
  background-color: #dcdcdc;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 8px;
  border: 1px solid  #ccc;
  cursor: pointer;
  &:hover {
    background-color: #c8c8c8; // 悬浮时改变背景颜色
  }
}
.ytcontainer {
  height: 100vh; /* 设置父容器高度为视口的高度 */
}
.spaced-div {
  //width: 95%;
  margin-left: 28px;  /* 左间隔 */
  margin-right: 28px; /* 右间隔 */
}
.menus {
  flex: 1;
  height: 100%;
}
.action {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  //justify-content: center;
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
.hover-red:hover {
  background: rgba(255, 0, 0, 0.6);
}

.titlebar {
  user-select: none;
  display: flex;
  padding: 4px 12px;
  align-items: center;
  margin-bottom: 6px;
}
.logo {
  width: 24px;
  height: 24px;
  user-select: none;
}
.appname {
  font-size: 14px;
  margin-left: 8px;
  //margin-top: 2px;
  font-weight: bold;
}

.custom-input {
  //cursor: not-allowed; /* 禁用点击时的光标样式 */
  background-color: #dcdcdc;
  width: 100%;
  padding: 7px;
  border: 1px solid  #ccc;
  border-radius: 4px;
  box-sizing: border-box;
  //cursor: pointer;
  outline: none; /* 去除默认聚焦边框 */
}

.check-button{
  //height: 30px;
  width: 80px;
  border-radius: 10px;
  background-image: linear-gradient(to right, #add8e6, #d8bfd8);
  color: black; /* 按钮文本颜色 */
  font-size: 12px;
  font-weight:400;
}

.check-button:hover{
  background-image: linear-gradient(to right, #a3cedc, #d8bfd8);
  color: black; /* 按钮文本颜色 */
}
.cancel-button{
  width: 80px;
  border-radius: 10px;
  background-color: rgb(210, 210, 210);
  font-size: 12px;
  font-weight:400;
}
.cancel-button:hover{
  background-color: rgb(200, 200, 200);
  color: black; /* 按钮文本颜色 */
}

.custom-option {
  height: 20px;
}
</style>