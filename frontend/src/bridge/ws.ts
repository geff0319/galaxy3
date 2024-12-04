// import * as App from "@wails/go/bridge/App";
import * as App from '@/bindings/github.com/geff0319/galaxy3/bridge/app'

export const appConnectWs = async (domain:string, id:string) => {
    const { flag, data } = await App.ConnectWs(domain,id)
    if (!flag) {
        throw data
    }
    return data
}

export const appDisConnectWs = async () => {
    const { flag, data } = await App.DisConnectWs()
    if (!flag) {
        throw data
    }
    return data
}

export const appPing = async (domain:string, id:string) => {
    const { flag, data } = await App.Ping(domain,id)
    if (!flag) {
        throw data
    }
    return data
}
