import {defineStore} from "pinia";
import {ref, watch} from "vue";
import {Readfile, Writefile} from "@/bridge";
import {parse, stringify} from "yaml";
import {debounce} from "@/utils";
import {message} from "ant-design-vue";
import {appConnectMqtt, appDisconnectMqtt, appStatusMqtt} from "@/bridge/mqtt";
import {message as antmessage} from "ant-design-vue/es/components";


export type MqttInfoType = {
    broker:string
    port:string
    clientID:string
    userName:string
    password:  string
    autoConnect:boolean
}

function generateRandomId(length: number): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    for (let i = 0; i < length; i++) {
        const randomIndex = Math.floor(Math.random() * chars.length);
        result += chars[randomIndex];
    }
    return result;
}

export const useMqttClientStore = defineStore('mqttClient', () => {
    let firstOpen = true
    let latestWsConfig = ''
    const isConnected = ref(false)
    const mqttInfo = ref<MqttInfoType>({
        broker:'',
        port:'10000',
        clientID:generateRandomId(8),
        userName:'',
        password:  '',
        autoConnect:false
    })
    const connectMqtt =async ()=>{
        try{
            const res = await appConnectMqtt(mqttInfo.value)
            isConnected.value = true
            message.info(res)
        }catch (e:any){
            isConnected.value = false
            message.error(e)
        }
    }
    const disConnectMqtt =async ()=>{
        try{
            const res = await appDisconnectMqtt()
            isConnected.value = false
            message.info(res)
        }catch (e:any){
            message.error(e)
        }
    }
    const connectionStatus =async ()=>{

        try{
            isConnected.value = await appStatusMqtt()
        }catch (e:any){
            message.error(e)
        }
    }

    const setupMqttSettings = async () => {
        try {
            const b = await Readfile('data/mqtt.yaml')
            mqttInfo.value = Object.assign(mqttInfo.value, parse(b))
            if(mqttInfo.value.autoConnect){
                console.log("mqtt init connect")
                await connectMqtt()
            }
        } catch (error) {
            firstOpen = false
            console.log(error)
        }
    }
    const saveMqttSettings = debounce((config: string) => {
        console.log('save mqtt settings')
        Writefile('data/mqtt.yaml', config)

    }, 1500)

    watch(
        mqttInfo,
        (settings) => {
            if (!firstOpen) {
                const lastModifiedConfig = stringify(settings)
                if (latestWsConfig !== lastModifiedConfig) {
                    saveMqttSettings(lastModifiedConfig).then(() => {
                        latestWsConfig = lastModifiedConfig
                    })
                } else {
                    saveMqttSettings.cancel()
                }
            }

            firstOpen = false
        },
        { deep: true }
    )
    return {
        connectMqtt,
        disConnectMqtt,
        connectionStatus,
        setupMqttSettings,
        mqttInfo,
        isConnected
    }
})