import {defineStore} from "pinia";
import {ref, watch} from "vue";
import {
    appAll,
    // appDbPersist,
    appDownloadYoutube,
    appDownloadYoutubeByKey,
    appGetVideoMeta, deleteProcess, getProcessByPage,
    UpdateYtDlpConfig
} from "@/bridge/ytdlp";
import {message} from "ant-design-vue";
import {Readfile, Writefile} from "@/bridge";
import {parse, stringify} from "yaml";
import {debounce} from "@/utils";
import {Execute, Select} from "@/bindings/github.com/geff0319/galaxy3/bridge/sqliteservice";


export type ProcessType = {
    id:string,
    pid:string
    url:string
    progress :{
        process_status:number
        percentage:string
        speed:number
        eta:number
    }
    info     :{
        url:string
        title :string
        thumbnail :string
        resolution :string
        filesize_approx:number
        size        :number
        vCodec      :string
        aCodec      :string
        extension   :string
        originalURL :string
        filename    :string
        created_at   :string
    }
    output   :{
        Path          :string
        filename      :string
        savedFilePath :string

    }
    params   : string[]
    biliMeta:{
        Vir:{
            data:{
                accept_description:string[]
                accept_quality:string[]
                support_formats:{
                    quality:string
                    new_description:string
                    codecs:string[]
                }[]
                dash:{
                    video:{
                        id:string
                        base_url:string
                        codecs:string
                        backupUrl:string[],
                        height: string
                    }[]
                }
            }
        }
        SelectedVideoStreamUrl:string
        SelectedVideoQuality:string
        SelectedVideoCodecs:string
    }
}

export type AllProcessType = {
    totalSize :number
    totalPage :number
    pageNum   :number
    pageSize :number
    processes : ProcessType[]
}

export function formatSize(bytes: number): string {
    const threshold = 1024
    const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    let i = 0
    while (bytes >= threshold) {
        bytes /= threshold
        i = i + 1
    }

    return `${bytes.toFixed(i === 0 ? 0 : 2)} ${units[i]}`;
}
export const formatSpeedMiB = (val: number) =>
    `${(val / 1_048_576).toFixed(2)}MiB/s`

export  const formattedFileSize = (val: number) => {
    if (val === -1) return '请登录';
    if (val === 0) return '0 Bytes';
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(val) / Math.log(1024));
    return `${(val / Math.pow(1024, i)).toFixed(2)} ${sizes[i]}`;
};

export const formatResolution=(val: string):string=>{
   if(val.includes('4320')){
       return '8K'
   }else if (val.includes('2160')||val.includes('3840')||val.includes('2048')){
       return '4K'
   }else if (val.includes('1440')){
       return '2K'
   }else if (val.includes('1080')||val.includes('960')){
       return '1080P'
   }else if (val.includes('720')){
       return '720P'
   }else if (val.includes('480')){
       return '480P'
   }else {
       return '未知'
   }
}

export const useYtdlpStore = defineStore('ytdlp', () => {
    const process = ref<ProcessType[]>([])
    const allProcess = ref<AllProcessType>({
        totalSize: 0,
        totalPage: 0,
        pageNum: 0,
        pageSize: 0,
        processes: [],
    })
    const resProcess = ref<ProcessType>({
        id: "",
        pid:"",
        url:"",
        output: {filename: "", Path: "", savedFilePath: ""},
        params: [],
        progress: {eta: 0, percentage: "", process_status: 0, speed: 0},
        info: {
            aCodec: "",
            created_at: "",
            extension: "",
            filename: "",
            originalURL: "",
            resolution: "",
            size: 0,
            filesize_approx: 0,
            thumbnail: "",
            title: "",
            url: "",
            vCodec: ""
        },
        biliMeta:{
            Vir:{
                data:{
                    accept_description:[],
                    accept_quality:[],
                    support_formats:[{
                        quality:'',
                        new_description:'',
                        codecs:[]
                    }],
                    dash:{
                        video:[{
                            id:'',
                            base_url:'',
                            codecs:'',
                            backupUrl:[],
                            height:''
                        }]
                    }
                }
            },
            SelectedVideoStreamUrl:'',
            SelectedVideoQuality:'',
            SelectedVideoCodecs:''
        }
    })
    const menuShow = ref(true)
    const videoUrl = ref<string>("") //复制的视频链接
    const downloadUrl = ref<string>("") //过滤后可下载视频链接
    const videoType= ref<string>('') //链接类型
    const youtubeRegex = /^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$/;
    const bilibiliRegex = /^https:\/\/www\.bilibili\.com\/video\/BV/;
    const parseing = ref<boolean>(false)
    const loading = ref<boolean>(false)
    const determineUrl = (url:string) =>{
        if(youtubeRegex.test(url)){
            getBaseUrl(url)
            videoType.value = 'YOUTUBE'
            return
        }else if (bilibiliRegex.test(url)){
            downloadUrl.value=url;
            videoType.value = 'BILIBILI'
            return;
        }
        resProcess.value.output.Path = ''
        resProcess.value.info.filesize_approx = 0
        resProcess.value.info.filename = ''
        downloadUrl.value = ''
        videoType.value = ''
    }

    const getBaseUrl = (url:string) =>{
        const index = url.indexOf('?list');
        if(index !== -1){
            downloadUrl.value=url.substring(0, index);
        }else{
            downloadUrl.value=url;
        }
    }

    const getVideoMeta =async ()=>{
        parseing.value=true
        if(downloadUrl.value.length===0){
            resProcess.value.output.Path = ""
            resProcess.value.info.filesize_approx = 0
            resProcess.value.info.filename = ''
            return
        }
        try {
            resProcess.value = await appGetVideoMeta(downloadUrl.value)
            console.log(resProcess.value)
        }catch (error:any){
            message.error(error,3)
        }
        parseing.value=false
    }
    const getAllVideoInfo = async () => {

        process.value = await appAll()
        console.log(process.value)
    }
    const getProcess= async (current: number, pageSize: number)=> {

        allProcess.value = await getProcessByPage(current, pageSize)
        console.log(allProcess.value)
    }

    const downloadYoutube =async (isKey:boolean,retry:boolean)=>{
        try {

            if(isKey){
                await appDownloadYoutubeByKey(resProcess.value,retry)
            }else {
                await appDownloadYoutube(downloadUrl.value,[])
            }
            // await getAllVideoInfo()

        }catch (error:any){
            throw error
        }
    }

    // const dbPersist = async ()=>{
    //     try {
    //         await appDbPersist()
    //     }catch (error:any){
    //         throw error
    //     }
    // }


    return {
        videoUrl,
        downloadUrl,
        videoType,
        parseing,
        process,
        resProcess,
        allProcess,
        menuShow,
        loading,
        getProcess,
        determineUrl,
        getVideoMeta,
        downloadYoutube,
        // dbPersist,
        getAllVideoInfo,
        formatSize,
        formatSpeedMiB,
        formatResolution,
        formattedFileSize
    }
})

type YtDlpCookie ={
    bilibili?:string
}
type YtdlpSetting = {
    downloadPath: string
    queueSize:string
    cookies:YtDlpCookie
}
export const useYtdlpSettingsStore = defineStore('ytdlp-settings', () =>{
    let latestYtdlpConfig = ''
    const ytdlpConfig = ref<YtdlpSetting>({cookies: {bilibili:""}, queueSize: "", downloadPath: ""})
    let firstOpen = true

    const setupYtdlpSettings = async () => {
        try {
            const res = await Select("SELECT config_value FROM config WHERE config_name = 'ytdlp' limit 1;")
            if(res !== null && res !== undefined && res.length !==0){
                if(res[0]?.["config_value"] ==='') {
                    await saveYtdlpSettings()
                } else {
                    ytdlpConfig.value = Object.assign(ytdlpConfig.value, parse(res[0]?.["config_value"]))
                }
            }else {
                //插入配置
                Execute("INSERT into config (config_name,config_value) VALUES(?,?)","ytdlp",stringify(ytdlpConfig.value))
            }
        } catch (error) {
            firstOpen = false
            console.log(error)
        }

    }
    const saveYtdlpSettings = debounce(async (config: string) => {
        console.log('save ytdlp settings')
        try {
            // await Writefile('data/ytdlp.yaml', config)
            Execute("UPDATE config SET  config_value= ? WHERE config_name = 'ytdlp';",config)

        }catch (error:any){
            message.error(error)
        }

    }, 1500)

    watch(
        ytdlpConfig,
        (settings) => {
            if (!firstOpen) {
                UpdateYtDlpConfig()
                const lastModifiedConfig = stringify(settings)
                if (latestYtdlpConfig !== lastModifiedConfig) {
                    saveYtdlpSettings(lastModifiedConfig).then(() => {
                        latestYtdlpConfig = lastModifiedConfig
                    })
                } else {
                    saveYtdlpSettings.cancel()
                }
            }
            firstOpen = false
        },
        { deep: true }
    )
    return {
        ytdlpConfig,
        setupYtdlpSettings,
    }
})