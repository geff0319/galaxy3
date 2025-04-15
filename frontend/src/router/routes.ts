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
import BallMenuView from '@/views/BallView/index.vue'
import SingleBallView from '@/views/BallView/SingleBall.vue'
import FlipClock from "@/views/FlipClockView/FlipClock.vue";



const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
    meta: {
      name: 'router.overview',
      newlayout: false,
      keepAlive:true,
      icon:"home"
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
      keepAlive:true,
      icon:"api"
    }
  },
  {
    path: '/translation',
    name: 'Translation',
    component: TranslationView,
    meta: {
      name: 'router.translation',
      keepAlive:true,
      icon:"translation"
    }
  },
  {
    path: '/ytdlp',
    name: 'Ytdlp',
    component: YtdlpView,
    meta: {
      name: '下载',
      keepAlive:false,
      icon:"download-three"
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
    },
    children: [
      {
        path: 'single-ball',  // 子路由路径，注意这里不用写 `/widgets/`，它会自动附加到父路由
        name: 'SingleBall',
        component: SingleBallView,  // 子路由组件
        meta: {
          name: '单悬浮球',
          newlayout: true,
          keepAlive:true
        }
      },
      {
        path: 'ball-menu',  // 子路由路径，注意这里不用写 `/widgets/`，它会自动附加到父路由
        name: 'BallMenu',
        component: BallMenuView,  // 子路由组件
        meta: {
          name: '悬浮球菜单',
          newlayout: true,
          keepAlive:true
        }
      },
      {
        path: 'clock',  // 子路由路径，注意这里不用写 `/widgets/`，它会自动附加到父路由
        name: 'Clock',
        component: FlipClock,  // 子路由组件
        meta: {
          name: '时钟',
          newlayout: true,
          keepAlive:true
        }
      }
    ]
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
      keepAlive:true,
      icon:"plan"
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SettingsView,
    meta: {
      name: 'router.settings',
      keepAlive:true,
      icon:"setting-two"
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
