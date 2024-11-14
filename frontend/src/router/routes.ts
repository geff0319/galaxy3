import { type RouteRecordRaw } from 'vue-router'

import { isDev } from '@/utils'

import HomeView from '@/views/HomeView/index.vue'
import SubscribesView from '@/views/SubscribesView/index.vue'
import SettingsView from '@/views/SettingsView/index.vue'
import ProfilesView from '@/views/ProfilesView/index.vue'
import RulesetsView from '@/views/RulesetsView/index.vue'
import PluginsView from '@/views/PluginsView/index.vue'
import ScheduledTasksView from '@/views/ScheduledTasksView/index.vue'
import PlaygroundView from '@/views/PlaygroundView/index.vue'
import TranslationView from '@/views/TranslationView/index.vue'
import YtdlpView from '@/views/YtdlpView/index.vue'
import WidgetsView from '@/views/WidgetsView.vue'
import YtdlpWidgetView from '@/views/YtdlpView/YtdlpWidget.vue'
import OptionView from '@/components/Option/index.vue'


const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
    meta: {
      name: 'router.overview',
      newlayout: false,
      keepAlive:true
    }
  },
  {
    path: '/profiles',
    name: 'Profiles',
    component: ProfilesView,
    meta: {
      name: 'router.profiles',
      hidden: !isDev,
      keepAlive:true
    }
  },
  {
    path: '/subscriptions',
    name: 'Subscriptions',
    component: SubscribesView,
    meta: {
      name: 'router.subscriptions',
      hidden: !isDev,
      keepAlive:true
    }
  },
  {
    path: '/rulesets',
    name: 'Rulesets',
    component: RulesetsView,
    meta: {
      name: 'router.rulesets',
      hidden: !isDev,
      keepAlive:true
    }
  },
  {
    path: '/plugins',
    name: 'PluginsView',
    component: PluginsView,
    meta: {
      name: 'router.plugins',
      keepAlive:true
    }
  },
  {
    path: '/translation',
    name: 'Translation',
    component: TranslationView,
    meta: {
      name: 'router.translation',
      keepAlive:true
    }
  },
  {
    path: '/ytdlp',
    name: 'Ytdlp',
    component: YtdlpView,
    meta: {
      name: '下载',
      keepAlive:false
    }
  },
  {
    path: '/widgets',
    name: 'Widgets',
    component: WidgetsView,
    meta: {
      name: '小组件',
      newlayout: true,
      hidden: true,
      keepAlive:true
    }
  },
  {
    path: '/ytdlpWidgets',
    name: 'YtdlpWidgets',
    component: YtdlpWidgetView,
    meta: {
      name: 'ytdlp组件',
      newlayout: true,
      hidden: true,
      keepAlive:true
    }
  },
  {
    path: '/scheduledtasks',
    name: 'ScheduledTasksView',
    component: ScheduledTasksView,
    meta: {
      name: 'router.scheduledtasks',
      keepAlive:true
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SettingsView,
    meta: {
      name: 'router.settings',
      keepAlive:true
    }
  },
  {
    path: '/playground',
    name: 'Playground',
    component: PlaygroundView,
    meta: {
      name: 'Playground',
      hidden: !isDev,
      keepAlive:true
    }
  },
  {
    path: '/option',
    name: 'Option',
    component: OptionView,
    meta: {
      name: 'option组件',
      newlayout: true,
      hidden: true,
      keepAlive:true
    },
  }
]

export default routes
