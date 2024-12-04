// import * as App from '@/bindings/galaxy3/bridge/app'
import * as App from '@/bindings/github.com/geff0319/galaxy3/bridge/app'
import type {MqttInfoType} from "@/stores/mqttClient";

export const appConnectMqtt = async (info:MqttInfoType) => {
    const { flag, data } = await App.ConnectMqtt(JSON.stringify(info))
    if (!flag) {
        throw data
    }
    return data
}

export const appDisconnectMqtt = async () => {
    const { flag, data } = await App.DisconnectMqtt()
    if (!flag) {
        throw data
    }
    return data
}

export const appStatusMqtt = async () => {
    const { flag, data,msg } = await App.StatusMqtt()
    if (!flag) {
        throw msg
    }
    return data
}