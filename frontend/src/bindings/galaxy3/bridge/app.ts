// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * App struct
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

export function AbsolutePath(path: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2839251837, path) as any;
    return $resultPromise;
}

export function AddScheduledTask(spec: string, event: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1181475892, spec, event) as any;
    return $resultPromise;
}

export function All(): Promise<$models.FlagResultWithData> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2269700366) as any;
    return $resultPromise;
}

export function ChangeLog(level: number, path: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(803743903, level, path) as any;
    return $resultPromise;
}

export function CheckBiliLogin(): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4153684584) as any;
    return $resultPromise;
}

export function ConnectWs(domain: string, id: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1127748881, domain, id) as any;
    return $resultPromise;
}

export function Delete(id: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1281394712, id) as any;
    return $resultPromise;
}

export function DisConnectWs(): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(926751891) as any;
    return $resultPromise;
}

export function Download(url: string, path: string, headers: { [_: string]: string } | null, event: string, proxy: string): Promise<$models.HTTPResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4073986687, url, path, headers, event, proxy) as any;
    return $resultPromise;
}

export function DownloadYoutube(url: string, params: string[] | null): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3031220380, url, params) as any;
    return $resultPromise;
}

export function DownloadYoutubeByKey(p: string, retry: boolean): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3200070448, p, retry) as any;
    return $resultPromise;
}

export function Exec(path: string, args: string[] | null, options: $models.ExecOptions): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4281359550, path, args, options) as any;
    return $resultPromise;
}

export function ExecBackground(path: string, args: string[] | null, outEvent: string, endEvent: string, options: $models.ExecOptions): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2365543414, path, args, outEvent, endEvent, options) as any;
    return $resultPromise;
}

export function ExitApp(): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2741059230) as any;
    return $resultPromise;
}

export function ExitKey(): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1644529186) as any;
    return $resultPromise;
}

export function FileExists(path: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4027796637, path) as any;
    return $resultPromise;
}

/**
 * 返回option窗口展示位置，展示在鼠标正下方
 */
export function GetBelowWinPos(winWidth: number, winHeight: number): Promise<$models.FlagResultWithData> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4247551374, winWidth, winHeight) as any;
    return $resultPromise;
}

export function GetEnv(): Promise<$models.EnvResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2450364276) as any;
    return $resultPromise;
}

export function GetInterfaces(): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2943945273) as any;
    return $resultPromise;
}

/**
 * 获取视频清晰度和名称
 */
export function GetVideoMeta(url: string): Promise<$models.FlagResultWithData> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1044285561, url) as any;
    return $resultPromise;
}

export function HideToolWindow(): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2502583885) as any;
    return $resultPromise;
}

export function HttpDelete(url: string, header: { [_: string]: string } | null, proxy: string): Promise<$models.HTTPResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(119491128, url, header, proxy) as any;
    return $resultPromise;
}

export function HttpGet(url: string, header: { [_: string]: string } | null, proxy: string): Promise<$models.HTTPResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(166310981, url, header, proxy) as any;
    return $resultPromise;
}

export function HttpHead(url: string, headers: { [_: string]: string } | null): Promise<$models.HTTPResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1565176561, url, headers) as any;
    return $resultPromise;
}

export function HttpPost(url: string, header: { [_: string]: string } | null, body: string, proxy: string): Promise<$models.HTTPResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1558268637, url, header, body, proxy) as any;
    return $resultPromise;
}

export function HttpPut(url: string, header: { [_: string]: string } | null, body: string, proxy: string): Promise<$models.HTTPResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1963332568, url, header, body, proxy) as any;
    return $resultPromise;
}

export function KillProcess(pid: number): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1904003264, pid) as any;
    return $resultPromise;
}

export function Makedir(path: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4186104616, path) as any;
    return $resultPromise;
}

export function Movefile(source: string, target: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1276361108, source, target) as any;
    return $resultPromise;
}

export function Notify(title: string, message: string, icon: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3554731110, title, message, icon) as any;
    return $resultPromise;
}

export function OpenDirectoryDialog(): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4106968830) as any;
    return $resultPromise;
}

/**
 * 缓存的数据持久化到文件
 */
export function Persist(): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4260410081) as any;
    return $resultPromise;
}

export function Ping(domain: string, id: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4224225119, domain, id) as any;
    return $resultPromise;
}

export function ProcessInfo(pid: number): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(15724828, pid) as any;
    return $resultPromise;
}

export function Readfile(path: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1838894513, path) as any;
    return $resultPromise;
}

export function RemoveScheduledTask(id: number): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3113710027, id) as any;
    return $resultPromise;
}

export function Removefile(path: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3020309033, path) as any;
    return $resultPromise;
}

export function RestartApp(): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4044707745) as any;
    return $resultPromise;
}

export function ShowToolWindow(): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1074210372) as any;
    return $resultPromise;
}

export function TencentTextTranslate(sourceText: string, sourceLang: string, targetLang: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(2691435463, sourceText, sourceLang, targetLang) as any;
    return $resultPromise;
}

export function UnzipGZFile(path: string, output: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(990009444, path, output) as any;
    return $resultPromise;
}

export function UnzipZIPFile(path: string, output: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(82368564, path, output) as any;
    return $resultPromise;
}

export function UpdateTray(tray: $models.TrayContent): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3846716102, tray) as any;
    return $resultPromise;
}

export function UpdateTrayMenus(menus: $models.MenuItem[] | null): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4180745790, menus) as any;
    return $resultPromise;
}

export function UpdateYtDlpConfig(): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1152808639) as any;
    return $resultPromise;
}

export function Upload(url: string, path: string, headers: { [_: string]: string } | null, event: string, proxy: string): Promise<$models.HTTPResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1016883966, url, path, headers, event, proxy) as any;
    return $resultPromise;
}

export function Writefile(path: string, content: string): Promise<$models.FlagResult> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4078625304, path, content) as any;
    return $resultPromise;
}
