<script setup lang="ts">
import {computed, ref} from "vue";
import {useI18n} from "vue-i18n";
import {getKernelFileName, getYtDlpFileName, KernelWorkDirectory, YtDlpWorkDirectory} from "@/constant";
import {Download, Exec, GetEnv, HttpGet, Makedir, Movefile, Removefile, UnzipGZFile, UnzipZIPFile} from "@/bridge";
import { useMessage } from '@/hooks'
import {useEnvStore} from "@/stores";
import {ignoredError} from "@/utils";
const { t } = useI18n()
const { message } = useMessage()

const releaseUrl = 'https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest'
const localVersion = ref('')
const remoteVersion = ref('')
const localVersionLoading = ref(false)
const remoteVersionLoading = ref(false)
const downloadLoading = ref(false)
const downloadSuccessful = ref(false)

const needUpdate = computed(() => remoteVersion.value && localVersion.value !== remoteVersion.value)
const updateLocalVersion = async (showTips = false) => {
  localVersion.value = await getLocalVersion(showTips)
}
const getLocalVersion = async (showTips = false) => {
  localVersionLoading.value = true
  try {
    const fileName = await getYtDlpFileName()
    const kernelFilePath = YtDlpWorkDirectory + '/' + fileName
    const res = await Exec(kernelFilePath, ['--version'])
    return res.trim()
  } catch (error: any) {
    console.log(error)
    showTips && message.error(error)
  } finally {
    localVersionLoading.value = false
  }
  return ''
}
const updateRemoteVersion = async () => {
  remoteVersion.value = await getRemoteVersion()
}
const getRemoteVersion = async () => {
  remoteVersionLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(releaseUrl)
    const { tag_name } = body
    return tag_name as string
  } catch (error: any) {
    console.log(error)
    message.error(error)
  } finally {
    remoteVersionLoading.value = false
  }
  return ''
}
const downloadCore = async () => {
  downloadLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(releaseUrl)
    const { os, arch } = await GetEnv()

    const { assets, message: msg } = body
    if (msg) throw msg

    const envStore = useEnvStore()
    // const amd64Compatible = arch === 'amd64' && envStore.env.x64Level < 3 ? '-compatible' : ''
    // const suffix = { windows: '.zip', linux: '.gz', darwin: '.gz' }[os]
    const assetName = `yt-dlp.exe`

    const asset = assets.find((v: any) => v.name === assetName)
    if (!asset) throw 'Asset Not Found:' + assetName

    const tmp = `data/yt-dlp.exe` // data/core-alpha.zip or data/core-alpha.gz

    await Makedir('data/yt-dlp')

    const { id } = message.info('Downloading...', 10 * 60 * 1_000)

    await Download(asset.browser_download_url, tmp, undefined, (event) => {
      message.update(id, 'Downloading...' + ((event.data[0] / event.data[1]) * 100).toFixed(2) + '%')
    }).catch((err) => {
      message.destroy(id)
      throw err
    })

    message.destroy(id)

    const fileName = await getYtDlpFileName() // yt-dlp.exe

    const ytdlpFilePath = YtDlpWorkDirectory + '/' + fileName

    await ignoredError(Movefile, ytdlpFilePath, ytdlpFilePath + '.bak')

    // if (suffix === '.zip') {
    //   await UnzipZIPFile(tmp, 'data')
    // } else {
    //   await UnzipGZFile(tmp, alphaKernelFilePath)
    // }

    await Movefile(tmp, ytdlpFilePath)

    await Removefile(tmp)

    downloadSuccessful.value = true

    message.success('Download Successful')
  } catch (error: any) {
    console.log(error)
    message.error(error)
    downloadSuccessful.value = false
  }

  downloadLoading.value = false

  updateLocalVersion()
}

const ffUrl = 'https://api.github.com/repos/yt-dlp/FFmpeg-Builds/releases/tags/latest'
const localFfVersion = ref('')
const localFfVersionLoading = ref(false)
const downloadFfLoading = ref(false)
const downloadFfSuccessful = ref(false)

// const needFfUpdate = computed(() => remoteFfVersion.value && localFfVersion.value !== remoteFfVersion.value)
const updateFfLocalVersion = async (showTips = false) => {
  localFfVersion.value = await getLocalFfVersion(showTips)
}
const getLocalFfVersion = async (showTips = false) => {
  localFfVersionLoading.value = true
  try {
    const fileName = "ffmpeg.exe"
    const kernelFilePath = YtDlpWorkDirectory + '/' + fileName
    const res = await Exec(kernelFilePath, ['-version'])
    return res.trim().match(/version\s([^\s]+)\sCopyright/)?.[1].trim() || ''
  } catch (error: any) {
    console.log(error)
    showTips && message.error(error)
  } finally {
    localFfVersionLoading.value = false
  }
  return ''
}
const downloadFfCore = async () => {
  downloadFfLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(ffUrl)
    const { os, arch } = await GetEnv()

    const { assets, message: msg } = body
    if (msg) throw msg

    const envStore = useEnvStore()
    const amd64Compatible = arch === 'amd64' && envStore.env.x64Level < 3 ? '-compatible' : ''
    const suffix = { windows: '.zip', linux: '.tar.xz', darwin: '.gz' }[os]
    const assetName = `ffmpeg-master-latest-win64-gpl${suffix}`
    const asset = assets.find((v: any) => v.name === assetName)
    if (!asset) throw 'Asset Not Found:' + assetName

    const tmp = `data/core-ffmpeg${suffix}` // data/core-alpha.zip or data/core-alpha.gz

    await Makedir('data/yt-dlp')

    const { id } = message.info('Downloading...', 10 * 60 * 1_000)

    await Download(asset.browser_download_url, tmp, undefined, (event) => {
      message.update(id, 'Downloading...' + ((event.data[0] / event.data[1]) * 100).toFixed(2) + '%')
    }).catch((err) => {
      message.destroy(id)
      throw err
    })

    message.destroy(id)

    // const fileName = await getKernelFileName() // mihomo-windows-amd64.exe
    // const alphaFileName = await getKernelFileName(true) // mihomo-windows-amd64-alpha.exe

    const ffmpegFilePath = YtDlpWorkDirectory + '/ffmpeg.exe'
    const ffplayFilePath = YtDlpWorkDirectory + '/ffplay.exe'
    const ffprobeFilePath = YtDlpWorkDirectory + '/ffprobe.exe'

    await ignoredError(Movefile, ffmpegFilePath, ffmpegFilePath + '.bak')
    await ignoredError(Movefile, ffplayFilePath, ffplayFilePath + '.bak')
    await ignoredError(Movefile, ffprobeFilePath, ffprobeFilePath + '.bak')

    if (suffix === '.zip') {
      await UnzipZIPFile(tmp, 'data')
    } else {
      await UnzipGZFile(tmp, 'data')
    }
    //
    await Movefile('data/ffmpeg-master-latest-win64-gpl/bin/ffmpeg.exe' , ffmpegFilePath)
    await Movefile('data/ffmpeg-master-latest-win64-gpl/bin/ffplay.exe' , ffplayFilePath)
    await Movefile('data/ffmpeg-master-latest-win64-gpl/bin/ffprobe.exe' , ffprobeFilePath)

    await Removefile(tmp)
    await Removefile('data/ffmpeg-master-latest-win64-gpl')

    downloadFfSuccessful.value = true

    message.success('Download Successful')
  } catch (error: any) {
    console.log(error)
    message.error(error)
    downloadFfSuccessful.value = false
  }

  downloadFfLoading.value = false

  updateFfLocalVersion(true)
}
const initVersion = async () => {
  getLocalVersion(false)
      .then((v) => {
        localVersion.value = v
      })
      .catch((error: any) => {
        console.log(error)
      })

  getRemoteVersion()
      .then((versions) => {
        remoteVersion.value = versions
      })
      .catch((error: any) => {
        console.log(error)
      })

  getLocalFfVersion(false)
      .then((v) => {
        localFfVersion.value = v
      })
      .catch((error: any) => {
        console.log(error)
      })
}

initVersion()
</script>

<template>
<!--  <h3>YT-DLP</h3>-->
<!--  <div style="display: flex;align-items: center;">-->
<!--    <Tag @click="updateLocalVersion(true)" style="cursor: pointer">-->
<!--      {{ t('kernel.local') }}-->
<!--      :-->
<!--      {{ localVersionLoading ? 'Loading' : localVersion || t('kernel.notFound') }}-->
<!--    </Tag>-->
<!--    <Tag @click="updateRemoteVersion()" style="cursor: pointer">-->
<!--      {{ t('kernel.remote') }}-->
<!--      :-->
<!--      {{ remoteVersionLoading ? 'Loading' : remoteVersion }}-->
<!--    </Tag>-->
<!--    <Button-->
<!--        v-show="!localVersionLoading && !remoteVersionLoading && needUpdate"-->
<!--        @click="downloadCore"-->
<!--        :loading="downloadLoading"-->
<!--        size="small"-->
<!--        type="primary"-->
<!--    >-->
<!--      {{ t('kernel.update') }} : {{ remoteVersion }}-->
<!--    </Button>-->
<!--  </div>-->
<!--  <h3>FFmpeg(YT-DLP优化版)</h3>-->
<!--  <div style="display: flex;align-items: center;">-->
<!--    <Tag @click="updateFfLocalVersion(true)" style="cursor: pointer">-->
<!--      {{ t('kernel.local') }}-->
<!--      :-->
<!--      {{ localFfVersionLoading ? 'Loading' : localFfVersion || t('kernel.notFound') }}-->
<!--    </Tag>-->
<!--    <Button-->
<!--        v-show="!localFfVersionLoading"-->
<!--        @click="downloadFfCore"-->
<!--        :loading="downloadFfLoading"-->
<!--        size="small"-->
<!--        type="primary"-->
<!--    >-->
<!--      {{ t('kernel.update') }}-->
<!--    </Button>-->
<!--  </div>-->



  <div class="settings">
    <div class="settings-item">
      <div class="title">YT-DLP</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div>版本管理</div>
          <div style="display: flex">
            <button class="button" @click="updateLocalVersion(true)">
              {{ t('kernel.local') }} : {{ localVersionLoading ? 'Loading' : localVersion || t('kernel.notFound') }}
            </button>
            <a @click="downloadCore" v-if="!localVersionLoading && !remoteVersionLoading && needUpdate">
              <a-badge :count="1" title="点击升级">
                <button class="button" @click="updateRemoteVersion()">
                  {{ t('kernel.remote') }} : {{ remoteVersionLoading ? 'Loading' : remoteVersion }}
                </button>
              </a-badge>
            </a>
            <button v-else class="button" @click="updateRemoteVersion()">
              {{ t('kernel.remote') }} : {{ remoteVersionLoading ? 'Loading' : remoteVersion }}
            </button>
          </div>
        </div>
      </a-card>
    </div>

    <div class="settings-item">
      <div class="title">FFMPEG(YT-DLP优化版)</div>
      <a-card class="card" size="small">
        <div class="card-item">
          <div>版本管理</div>
          <a @click="downloadFfCore" v-if="!localFfVersionLoading">
            <a-badge :count="1" title="点击升级">
              <button class="button" @click="updateFfLocalVersion(true)">
                {{ t('kernel.local') }} : {{ localFfVersionLoading ? 'Loading' : localFfVersion || t('kernel.notFound') }}
              </button>
            </a-badge>
          </a>
          <button v-else class="button" @click="updateFfLocalVersion(true)">
            {{ t('kernel.local') }} : {{ localFfVersionLoading ? 'Loading' : localFfVersion || t('kernel.notFound') }}
          </button>
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
.button {
  cursor:pointer;
  border: 1px solid  #ccc;
  border-radius: 4px;
  box-sizing: border-box;
  padding: 6px 20px;
  margin-left: 20px;
  background-color: rgb(250, 250, 250);
  font-family: "幼圆", "Yu Yuan", sans-serif;
}
</style>