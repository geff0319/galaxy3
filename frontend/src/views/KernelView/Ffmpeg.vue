<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { computed, ref } from 'vue'

import { useMessage } from '@/hooks'
import { ignoredError } from '@/utils'
import {KernelWorkDirectory, getKernelFileName, getFfmpegFileName, FfmpegWorkDirectory} from '@/constant'
import { useAppSettingsStore, useEnvStore, useKernelApiStore } from '@/stores'
import {
  Download,
  UnzipZIPFile,
  HttpGet,
  Exec,
  Movefile,
  Removefile,
  GetEnv,
  Makedir,
  UnzipGZFile
} from '@/bridge'

const releaseUrl = 'https://api.github.com/repos/geff0319/ffmpeg/releases/latest'
const localVersion = ref('')
const remoteVersion = ref('')
const versionDetail = ref('')
const localVersionLoading = ref(false)
const remoteVersionLoading = ref(false)
const downloadLoading = ref(false)
const downloadSuccessful = ref(false)

const needRestart = computed(() => {
  const { running, branch } = appSettings.app.kernel
  if (!running) return false
  return localVersion.value && downloadSuccessful.value && branch === 'main'
})

const needUpdate = computed(() => remoteVersion.value && localVersion.value !== remoteVersion.value)

const { t } = useI18n()
const { message } = useMessage()
const appSettings = useAppSettingsStore()
const kernelApiStore = useKernelApiStore()

const updateLocalVersion = async (showTips = false) => {
  localVersion.value = await getLocalVersion(showTips)
}

const updateRemoteVersion = async (showTips = false) => {
  remoteVersion.value = await getRemoteVersion(showTips)
}

const downloadCore = async () => {
  downloadLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(releaseUrl)
    const { os, arch } = await GetEnv()

    const { assets, tag_name, message: msg } = body
    if (msg) throw msg

    const envStore = useEnvStore()
    const amd64Compatible = arch === 'amd64' && envStore.env.x64Level < 3 ? '-compatible' : ''
    const suffix = { windows: '.zip', linux: '.gz', darwin: '.gz' }[os]
    const assetName = `ffmpeg.exe`

    const asset = assets.find((v: any) => v.name === assetName)
    if (!asset) throw 'Asset Not Found:' + assetName

    const tmp = `data/core` // data/core.zip or data/core.gz

    await Makedir('data/ffmpeg')

    const { id } = message.info('Downloading...', 10 * 60 * 1_000)

    await Download(asset.browser_download_url, tmp, undefined, (event) => {
      // message.update(id, 'Downloading...' + ((progress / total) * 100).toFixed(2) + '%')
      message.update(id, 'Downloading...' + ((event.data[0] / event.data[1]) * 100).toFixed(2) + '%')
    }).catch((err) => {
      message.destroy(id)
      throw err
    })

    message.destroy(id)

    const fileName = await getFfmpegFileName()

    const kernelFilePath = FfmpegWorkDirectory + '/' + fileName

    await ignoredError(Movefile, kernelFilePath, kernelFilePath + localVersion.value + '.bak')

    await ignoredError(Movefile,tmp,kernelFilePath)
    // if (suffix === '.zip') {
    //   await UnzipZIPFile(tmp, FfmpegWorkDirectory)
    // } else {
    //   await UnzipGZFile(tmp, kernelFilePath)
    // }

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

const getLocalVersion = async (showTips = false) => {
  localVersionLoading.value = true
  try {
    const fileName = await getFfmpegFileName()
    const FfmpegFilePath = FfmpegWorkDirectory + '/' + fileName
    const res = await Exec(FfmpegFilePath, ['-version'])
    // versionDetail.value = res.trim().match(/version\s([^\s]+)\sCopyright/)?.[1].trim() || ''
    versionDetail.value = ''
    return res.trim().match(/version\s([^\s]+)\sCopyright/)?.[1].trim() || ''
  } catch (error: any) {
    console.log(error)
    showTips && message.error(error)
  } finally {
    localVersionLoading.value = false
  }
  return ''
}

const getRemoteVersion = async (showTips = false) => {
  remoteVersionLoading.value = true
  try {
    const { body } = await HttpGet<Record<string, any>>(releaseUrl)
    const { tag_name } = body
    return tag_name as string
  } catch (error: any) {
    console.log(error)
    showTips && message.error(error)
  } finally {
    remoteVersionLoading.value = false
  }
  return ''
}

const handleRestartKernel = async () => {
  if (!appSettings.app.kernel.running) return

  try {
    await kernelApiStore.restartKernel()

    downloadSuccessful.value = false

    message.success('common.success')
  } catch (error: any) {
    message.error(error)
  }
}

const initVersion = async () => {
  getLocalVersion()
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
}

initVersion()
</script>

<template>
  <h3>ffmpeg</h3>
  <Tag @click="updateLocalVersion(true)" style="cursor: pointer">
    {{ t('kernel.local') }}
    :
    {{ localVersionLoading ? 'Loading' : localVersion || t('kernel.notFound') }}
  </Tag>
  <Tag @click="updateRemoteVersion(true)" style="cursor: pointer">
    {{ t('kernel.remote') }}
    :
    {{ remoteVersionLoading ? 'Loading' : remoteVersion }}
  </Tag>
  <Button
    v-show="!localVersionLoading && !remoteVersionLoading && needUpdate"
    @click="downloadCore"
    :loading="downloadLoading"
    size="small"
    type="primary"
  >
    {{ t('kernel.update') }}
  </Button>
  <Button
    v-show="!localVersionLoading && !remoteVersionLoading && needRestart"
    @click="handleRestartKernel"
    :loading="kernelApiStore.loading"
    size="small"
    type="primary"
  >
    {{ t('kernel.restart') }}
  </Button>
  <div class="detail">
    {{ versionDetail }}
  </div>
</template>

<style lang="less" scoped>
.detail {
  font-size: 12px;
  padding: 8px 4px;
  word-wrap: break-word;
  word-break: break-all;
}
</style>
