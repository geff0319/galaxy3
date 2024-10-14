// import * as App from '@wails/go/bridge/App'
import * as App from '@/bindings/galaxy3/bridge/app'

export const TencentTextTranslate = async (sourceText: string, sourceLang: string ,targetLang: string) => {
    const { flag, data } = await App.TencentTextTranslate(sourceText, sourceLang,targetLang)
    if (!flag) {
        throw data
    }
    return data
}