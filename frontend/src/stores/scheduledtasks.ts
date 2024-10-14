import { defineStore } from 'pinia'
import { stringify, parse } from 'yaml'
import { computed, ref, watch } from 'vue'

import i18n from '@/lang'
import {FileExists, Notify} from '@/bridge'
import { debounce } from '@/utils'
import { ScheduledTasksFilePath, ScheduledTasksType } from '@/constant'
import { useSubscribesStore, useRulesetsStore, usePluginsStore, useLogsStore } from '@/stores'
import {
  Readfile,
  Writefile,
  AddScheduledTask,
  RemoveScheduledTask,
  Events
} from '@/bridge'
import {appDbPersist} from "@/bridge/ytdlp";

export type ScheduledTaskType = {
  id: string
  name: string
  type: ScheduledTasksType
  subscriptions: string[]
  rulesets: string[]
  plugins: string[]
  script: string
  cron: string
  notification: boolean
  disabled: boolean
  lastTime: string
}

export const useScheduledTasksStore = defineStore('scheduledtasks', () => {
  const scheduledtasks = ref<ScheduledTaskType[]>([])
  const ScheduledTasksEvents: string[] = []
  const ScheduledTasksIDs: number[] = []

  const setupScheduledTasks = async () => {
    const stat = await FileExists(ScheduledTasksFilePath)
    if (!stat){
      let ytdlptasks = {
        id: "ID_YTDLP",
        name: "ytdlp定时保存",
        type: ScheduledTasksType.RunScript,
        subscriptions: [],
        rulesets: [],
        plugins: [],
        script: '',
        cron: '0 */5 * * * *',
        notification: false,
        disabled: false,
        lastTime: ''
      }
      console.log("添加ytdlp的定时保存任务")
      await addScheduledTask(ytdlptasks)
      return
    }
    const data = await Readfile(ScheduledTasksFilePath)
    scheduledtasks.value = parse(data)
    // 添加ytdlp的定时保存任务
    const exists = scheduledtasks.value.some(task => task.id === "ID_YTDLP");
    if(!exists){
      let ytdlptasks = {
        id: "ID_YTDLP",
        name: "ytdlp定时保存",
        type: ScheduledTasksType.RunScript,
        subscriptions: [],
        rulesets: [],
        plugins: [],
        script: '',
        cron: '0 */5 * * * *',
        notification: false,
        disabled: false,
        lastTime: ''
      }
      console.log("添加ytdlp的定时保存任务")
      await addScheduledTask(ytdlptasks)
    }
  }

  const initScheduledTasks = async () => {
    removeScheduledTasks()

    const { t } = i18n.global
    const logsStore = useLogsStore()

    for (const {disabled, cron, id} of scheduledtasks.value) {
      if (disabled) continue;
      const taskID = await AddScheduledTask(cron, id)
      ScheduledTasksEvents.push(id)
      ScheduledTasksIDs.push(taskID)
      if(id === "ID_YTDLP"){
        Events.On(id, async () => {
          const task = getScheduledTaskById(id)
          if (!task) return

          task.lastTime = new Date().toLocaleString()
          editScheduledTask(id, task)
          const startTime = Date.now()
          const result= await appDbPersist()
          const resultArray = [result]

          task.notification && Notify(task.name, result)

          logsStore.recordScheduledTasksLog({
            name: task.name,
            startTime,
            endTime: Date.now(),
            result:resultArray
          })
        })
      }else {
        Events.On(id, async () => {
          const task = getScheduledTaskById(id)
          if (!task) return

          task.lastTime = new Date().toLocaleString()
          editScheduledTask(id, task)

          const startTime = Date.now()
          const result = await getTaskFn(task)()

          task.notification && Notify(task.name, result.join('\n'))

          logsStore.recordScheduledTasksLog({
            name: task.name,
            startTime,
            endTime: Date.now(),
            result
          })
        })
      }

    }
  }

  const removeScheduledTasks = () => {
    ScheduledTasksEvents.forEach((event) => Events.Off(event))
    ScheduledTasksIDs.forEach((id) => RemoveScheduledTask(id))
    ScheduledTasksEvents.splice(0)
    ScheduledTasksIDs.splice(0)
  }

  const withOutput = (list: string[], fn: (id: string) => Promise<string>) => {
    return async () => {
      const output = []
      for (const id of list) {
        try {
          const res = await fn(id)
          output.push(res)
        } catch (error: any) {
          output.push(error.message || error)
        }
      }
      return output
    }
  }

  const getTaskFn = (task: ScheduledTaskType) => {
    switch (task.type) {
      case ScheduledTasksType.UpdateSubscription: {
        const subscribesStore = useSubscribesStore()
        return withOutput(task.subscriptions, subscribesStore.updateSubscribe)
      }
      case ScheduledTasksType.UpdateRuleset: {
        const rulesetsStore = useRulesetsStore()
        return withOutput(task.rulesets, rulesetsStore.updateRuleset)
      }
      case ScheduledTasksType.UpdatePlugin: {
        const pluginsStores = usePluginsStore()
        return withOutput(task.plugins, pluginsStores.updatePlugin)
      }
      case ScheduledTasksType.RunPlugin: {
        const pluginsStores = usePluginsStore()
        return withOutput(task.plugins, async (id: string) =>
          pluginsStores.manualTrigger(id, 'onTask' as any)
        )
      }
      case ScheduledTasksType.RunScript: {
        const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor
        return withOutput([task.script], (script: string) => new AsyncFunction(script)())
      }
    }
  }

  const saveScheduledTasks = debounce(async () => {
    await Writefile(ScheduledTasksFilePath, stringify(scheduledtasks.value))
  }, 500)

  const addScheduledTask = async (s: ScheduledTaskType) => {
    scheduledtasks.value.push(s)
    try {
      await saveScheduledTasks()
    } catch (error) {
      scheduledtasks.value.pop()
      throw error
    }
  }

  const deleteScheduledTask = async (id: string) => {
    const idx = scheduledtasks.value.findIndex((v) => v.id === id)
    if (idx === -1) return
    const backup = scheduledtasks.value.splice(idx, 1)[0]
    try {
      await saveScheduledTasks()
    } catch (error) {
      scheduledtasks.value.splice(idx, 0, backup)
      throw error
    }
  }

  const editScheduledTask = async (id: string, s: ScheduledTaskType) => {
    const idx = scheduledtasks.value.findIndex((v) => v.id === id)
    if (idx === -1) return
    const backup = scheduledtasks.value.splice(idx, 1, s)[0]
    try {
      await saveScheduledTasks()
    } catch (error) {
      scheduledtasks.value.splice(idx, 1, backup)
      throw error
    }
  }

  const getScheduledTaskById = (id: string) => scheduledtasks.value.find((v) => v.id === id)

  const _watchCron = computed(() =>
    scheduledtasks.value
      .map((v) => v.cron)
      .sort()
      .join()
  )

  const _watchDisabled = computed(() =>
    scheduledtasks.value
      .map((v) => v.disabled)
      .sort()
      .join()
  )

  watch([_watchCron, _watchDisabled], () => {
    initScheduledTasks()
  })

  return {
    scheduledtasks,
    setupScheduledTasks,
    saveScheduledTasks,
    addScheduledTask,
    editScheduledTask,
    deleteScheduledTask,
    getScheduledTaskById,
    getTaskFn,
    removeScheduledTasks
  }
})
