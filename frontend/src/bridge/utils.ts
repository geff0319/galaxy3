import * as App from '@/bindings/galaxy3/bridge/app'

export const AppGetBelowWinPos = async (x:number,y:number) => {
    const { flag, msg,data } = await App.GetBelowWinPos(x,y)
    if (!flag) {
        throw msg
    }
    return data
}

export const AppChangeLog = async (level:number,path:string) => {
     await App.ChangeLog(level,path)
}