import {ref, watch} from 'vue'
import {defineStore} from 'pinia'
import {parse, stringify} from 'yaml'

import i18n from '@/lang'
import {debounce} from '@/utils'
import {Readfile, Writefile} from '@/bridge'
// import { Readfile, Writefile } from "@/bindings/galaxy3/bridge/app";
import {Color, Colors, Lang, Theme, View, WindowStartState} from '@/constant'
import {AppChangeLog} from "@/bridge/utils";

type AppSettings = {
  lang: Lang
  theme: Theme
  color: Color
  'font-family': string
  profilesView: View
  subscribesView: View
  rulesetsView: View
  pluginsView: View
  scheduledtasksView: View
  ytdlpView:View
  windowStartState: WindowStartState
  exitOnClose: boolean
  closeKernelOnExit: boolean
  autoSetSystemProxy: boolean
  autoStartKernel: boolean
  userAgent: string
  startupDelay: number
  connections: {
    visibility: Record<string, boolean>
    order: string[]
  }
  kernel: {
    branch: 'main' | 'alpha'
    profile: string
    pid: number
    running: boolean
    autoClose: boolean
    unAvailable: boolean
    cardMode: boolean
    sortByDelay: boolean
    testUrl: string
  }
  addPluginToMenu: boolean
  pluginSettings: Record<string, Record<string, any>>
  translate: Record<string, string>
  logPath:string
}

export const useAppSettingsStore = defineStore('app-settings', () => {
  let firstOpen = true
  let latestUserConfig = ''

  const themeMode = ref<Theme.Dark | Theme.Light>(Theme.Light)

  const app = ref<AppSettings>({
    lang: Lang.ZH,
    theme: Theme.Light,
    color: Color.Default,
    'font-family': '"幼圆", "Yu Yuan", sans-serif;',
    profilesView: View.Grid,
    subscribesView: View.Grid,
    rulesetsView: View.Grid,
    pluginsView: View.Grid,
    scheduledtasksView: View.Grid,
    ytdlpView:View.List,
    windowStartState: WindowStartState.Normal,
    exitOnClose: true,
    closeKernelOnExit: true,
    autoSetSystemProxy: false,
    autoStartKernel: false,
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36',
    startupDelay: 30,
    connections: {
      visibility: {
        'metadata.inboundName': true,
        'metadata.type': true,
        'metadata.process': false,
        'metadata.processPath': false,
        'metadata.host': true,
        'metadata.sniffHost': false,
        'metadata.sourceIP': false,
        'metadata.remoteDestination': false,
        rule: true,
        chains: true,
        up: true,
        down: true,
        upload: true,
        download: true,
        start: true
      },
      order: [
        'metadata.inboundName',
        'metadata.type',
        'metadata.process',
        'metadata.processPath',
        'metadata.host',
        'metadata.sniffHost',
        'metadata.sourceIP',
        'metadata.remoteDestination',
        'rule',
        'chains',
        'up',
        'down',
        'upload',
        'download',
        'start'
      ]
    },
    kernel: {
      branch: 'main',
      profile: '',
      pid: 0,
      running: false,
      autoClose: true,
      unAvailable: true,
      cardMode: true,
      sortByDelay: false,
      testUrl: 'https://www.gstatic.com/generate_204'
    },
    addPluginToMenu: false,
    pluginSettings: {},
    translate: {
      tencentTanslateSecretId: '',
      tencentTanslateSecretKey: ''
    },
    logPath:''
  })

  const saveAppSettings = debounce((config: string) => {
    console.log('save app settings')
    Writefile('data/user.yaml', config)

  }, 1500)

  const setupAppSettings = async () => {
    try {
      const b = await Readfile('data/user.yaml')
      app.value = Object.assign(app.value, parse(b))
    } catch (error) {
      firstOpen = false
      console.log(error)
    }

    updateAppSettings(app.value)
  }

  const mediaQueryList = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQueryList.addEventListener('change', ({ matches }) => {
    console.log('onSystemThemeChange')
    if (app.value.theme === Theme.Auto) {
      themeMode.value = matches ? Theme.Dark : Theme.Light
    }
  })

  const setAppTheme = (theme: Theme.Dark | Theme.Light) => {
    if (document.startViewTransition) {
      document.startViewTransition(() => {
        // console.log('set1'+theme)
        // document.body.setAttribute('theme-mode', theme)
      })
    } else {
      console.log('set2'+theme)
      document.body.setAttribute('theme-mode', theme)
    }
    // WindowSetSystemDefaultTheme()
  }

  const updateAppSettings = (settings: AppSettings) => {
    i18n.global.locale.value = settings.lang
    themeMode.value =
      settings.theme === Theme.Auto
        ? mediaQueryList.matches
          ? Theme.Dark
          : Theme.Light
        : settings.theme
    const { primary, secondary } = Colors[settings.color]
    document.documentElement.style.setProperty('--primary-color', primary)
    document.documentElement.style.setProperty('--secondary-color', secondary)
    document.body.style.fontFamily = settings['font-family']
    AppChangeLog(0,settings.logPath)
  }

  watch(
    app,
    (settings) => {
      updateAppSettings(settings)

      if (!firstOpen) {
        const lastModifiedConfig = stringify(settings)
        if (latestUserConfig !== lastModifiedConfig) {
          saveAppSettings(lastModifiedConfig).then(() => {
            latestUserConfig = lastModifiedConfig
          })
        } else {
          saveAppSettings.cancel()
        }
      }

      firstOpen = false
    },
    { deep: true }
  )

  // watch(
  //   [
  //     themeMode,
  //     () => app.value.color,
  //     () => app.value.lang,
  //     () => app.value.addPluginToMenu,
  //     // () => app.value.kernel.running,
  //     () => app.value.kernel.unAvailable,
  //     () => app.value.kernel.sortByDelay
  //   ],
  //     updateTrayMenus
  // )

  watch(themeMode, setAppTheme, { immediate: true })

  return { setupAppSettings, app, themeMode }
})
