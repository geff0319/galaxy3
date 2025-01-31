// import * as App from '@wails/go/bridge/App'
import * as App from '@/bindings/github.com/geff0319/galaxy3/bridge/app'
import type { TrayContent } from '@/constant'

// export * from '@wails/runtime/runtime'
export * from "@wailsio/runtime"


export const RestartApp = App.RestartApp

export const ExitApp = App.ExitApp

export const ExitKey = App.ExitKey

// export const UpdateTrayMenus = App.UpdateTrayMenus

export const UpdateTray = async (tray: TrayContent) => {
  const { icon = '', title = '', tooltip = '' } = tray
  await App.UpdateTray({
    icon,
    title,
    tooltip
  })
}

export const AddScheduledTask = async (cron: string, event: string) => {
  const { flag, data } = await App.AddScheduledTask(cron, event)
  if (!flag) {
    throw data
  }
  return Number(data)
}

export const RemoveScheduledTask = async (id: number) => {
  await App.RemoveScheduledTask(id)
}

export const Notify = async (title: string, message: string, icon = '') => {
  icon =
    {
      success: 'frontend/dist/notify_success.png',
      error: 'frontend/dist/notify_error.png'
    }[icon] || 'frontend/dist/favicon.ico'
  await App.Notify(title, message, icon)
}

export const HideToolWindow = async () =>{
  await App.HideToolWindow()
}

export const ShowToolWindow = async () =>{
  await App.ShowToolWindow()
}
