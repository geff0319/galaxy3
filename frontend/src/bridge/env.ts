// import * as App from '@wails/go/bridge/App'
import * as App from '@/bindings/galaxy3/bridge/app'
export const GetEnv = App.GetEnv

export const GetInterfaces = async () => {
  const { flag, data } = await App.GetInterfaces()
  if (!flag) {
    throw data
  }
  return data.split('|')
}
