package bridge

import (
	"encoding/json"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"github.com/geff0319/galaxy3/bridge/website"
	"math"
	"strconv"
	"time"
)

// 获取视频清晰度和名称
func (a *App) GetVideoMeta(url string) FlagResultWithData {
	p := &Process{
		Url:    url,
		Params: []string{},
		Output: DownloadOutput{
			Path: YdpConfig.DownloadPath,
		},
		BiliMeta: website.BiliMetadata{
			SelectedVideoQuality: "",
		},
	}
	p.SetPending()
	if err := p.SetMetadata(); err != nil {
		gefflog.Err("解析视频失败: " + err.Error())
		return FlagResultWithData{false, fmt.Sprintf("解析视频失败"), nil}
	}
	return FlagResultWithData{true, "解析视频成功", p}
}

func (a *App) DownloadYoutubeByKey(p string, retry bool) FlagResult {
	var process Process
	if err := json.Unmarshal([]byte(p), &process); err != nil {
		return FlagResult{false, "下载视频失败,解析任务信息出错"}
	}
	if !retry {
		if process.Info.Id == "" {
			return FlagResult{false, "下载视频失败,任务信息为空"}
		}
		if process.IsExist() >= 1 {
			return FlagResult{false, "下载视频失败,任务已存在"}
		} else {
			err := process.Insert()
			if err != nil {
				return FlagResult{true, "下载视频失败:" + err.Error()}
			}
			gefflog.Info(fmt.Sprintf("开始下载：%s", process.Pid))
			YdpConfig.Mq.PublishByTopic("process:downloading", &process)
			return FlagResult{true, strconv.FormatInt(process.Id, 10)}
		}
	} else {
		if process.Info.Id == "" {
			gefflog.Info(fmt.Sprintf("重试解析：%+v", process))
			YdpConfig.Mq.Publish(&process)
		} else {
			gefflog.Info(fmt.Sprintf("重试下载：%s", process.Pid))
			process.Progress.Status = StatusPending
			process.Update()
			YdpConfig.Mq.PublishByTopic("process:downloading", &process)
		}
		return FlagResult{true, strconv.FormatInt(process.Id, 10)}
	}

}

func (a *App) DownloadYoutube(url string, params []string) FlagResult {
	p := &Process{
		Url:    url,
		Params: params,
		Output: DownloadOutput{
			Path: YdpConfig.DownloadPath,
		},
		BiliMeta: website.BiliMetadata{
			SelectedVideoQuality: "",
		},
	}
	err := p.Insert()
	if err != nil {
		return FlagResult{Flag: false, Data: "添加失败:" + err.Error()}
	}
	YdpConfig.Mq.Publish(p)
	return FlagResult{Flag: true, Data: strconv.FormatInt(p.Id, 10)}
}

// 缓存的数据持久化到文件
//func (a *App) Persist() FlagResult {
//	err := YdpConfig.Mdb.Persist(Env.BasePath)
//	if err != nil {
//		gefflog.Err("视频信息保存失败：" + err.Error())
//		return FlagResult{Flag: false, Data: err.Error()}
//	}
//	return FlagResult{Flag: true, Data: "视频信息保存成功"}
//}

func (a *App) All() FlagResultWithData {
	res, err := SqliteS.Select(ProcessAll)
	if err != nil {
		gefflog.Err("SelectProcess err:" + err.Error())
	}

	//ps := []Process{}
	ps := make([]Process, 0, len(res))
	for _, r := range res {
		var p Process
		p.Unmarshal(r)
		if p.Progress.Status != StatusDownloading {
			//gefflog.Info(p.Info.FileName + ": 是下载完成的")
			ps = append(ps, p)
			continue
		}
		//gefflog.Info(p.Info.FileName + ": 是正在下载的")
		dstP, _ := YdpConfig.Mdb.Get(p.Id)
		if dstP != nil {
			ps = append(ps, *dstP)
		} else {
			ps = append(ps, p)
		}
	}
	return FlagResultWithData{
		Flag: true,
		Msg:  "操作成功",
		Data: ps,
	}
}

func AllFinish() AllProcess {
	//res, err := SqliteS.Select(`SELECT id,pid,url, params,info,progress,output,biliMeta FROM process where is_delete = 0 AND json_extract(progress, '$.process_status') = 2 ORDER BY create_time DESC`)
	//if err != nil {
	//	gefflog.Err("SelectProcess err:" + err.Error())
	//}
	currentTime := time.Now()
	gefflog.Err("自定义格式0:", currentTime.Format("2006年01月02日 15:04:05"))
	res, err := SqliteS.Select(`SELECT id,pid,url, params,info,progress,output,biliMeta FROM process where is_delete = 0 AND json_extract(progress, '$.process_status') = 2 ORDER BY create_time DESC LIMIT 10`)
	if err != nil {
		gefflog.Err("SelectProcess err:" + err.Error())
	}

	currentTime = time.Now()
	gefflog.Err("自定义格式1:", currentTime.Format("2006年01月02日 15:04:05"))

	//ps := []Process{}
	ps := make([]Process, 0, len(res))
	for _, r := range res {
		var p Process
		p.Unmarshal(r)
		ps = append(ps, p)

	}
	currentTime = time.Now()
	gefflog.Err("自定义格式2:n", currentTime.Format("2006年01月02日 15:04:05"))

	item, err := SqliteS.Select(`SELECT count(id) as num FROM process where is_delete = 0 AND json_extract(progress, '$.process_status') = 2`)
	if err != nil {
		gefflog.Err("SelectProcess err:" + err.Error())
	}
	var totalSize int64
	var totalPage int64

	currentTime = time.Now()
	gefflog.Err("自定义格式3:", currentTime.Format("2006年01月02日 15:04:05"))

	if num, ok := item[0]["num"]; ok {
		if t, ok := num.(int64); ok {
			totalSize = t
			totalPage = int64(math.Ceil(float64(totalSize) / float64(10)))
		}
	}

	currentTime = time.Now()
	gefflog.Info("自定义格式4: %s\n", currentTime.Format("2006年01月02日 15:04:05"))
	allProcess := AllProcess{
		TotalSize: totalSize,
		TotalPage: totalPage,
		PageNum:   1,
		PageSize:  10,
		Processes: ps,
	}
	return allProcess
}
func (a *App) GetProcessByPage(current int64, pageSize int64) FlagResultWithData {
	offset := (current - 1) * pageSize
	res, err := SqliteS.Select(`SELECT id,pid,url, params,info,progress,output,biliMeta FROM process where is_delete = 0 AND json_extract(progress, '$.process_status') = 2 ORDER BY create_time DESC LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		gefflog.Err("SelectProcess err:" + err.Error())
	}

	//ps := []Process{}
	ps := make([]Process, 0, len(res))
	for _, r := range res {
		var p Process
		p.Unmarshal(r)
		ps = append(ps, p)

	}

	item, err := SqliteS.Select(`SELECT count(id) as num FROM process where is_delete = 0 AND json_extract(progress, '$.process_status') = 2`)
	if err != nil {
		gefflog.Err("SelectProcess err:" + err.Error())
	}
	var totalSize int64
	var totalPage int64

	if num, ok := item[0]["num"]; ok {
		if t, ok := num.(int64); ok {
			totalSize = t
			totalPage = int64(math.Ceil(float64(totalSize) / float64(10)))
		}
	}
	allProcess := AllProcess{
		TotalSize: totalSize,
		TotalPage: totalPage,
		PageNum:   current,
		PageSize:  pageSize,
		Processes: ps,
	}
	return FlagResultWithData{
		Flag: true,
		Msg:  "操作成功",
		Data: allProcess,
	}
}

func downloading() {

}
func (a *App) UpdateYtDlpConfig() FlagResult {
	//b, err := os.ReadFile(Env.BasePath + "/data/ytdlp.yaml")
	//if err != nil {
	//	gefflog.Err("更新配置失败: " + err.Error())
	//	return FlagResult{false, "更新配置失败"}
	//}
	//if err := yaml.Unmarshal(b, &YdpConfig); err != nil {
	//	gefflog.Err("更新配置失败: " + err.Error())
	//	return FlagResult{false, "更新配置失败"}
	//}
	err := YdpConfig.Unmarshal()
	if err != nil {
		gefflog.Err("更新配置失败: " + err.Error())
		return FlagResult{false, "更新配置失败"}
	}
	if len(YdpConfig.Other.BilibiliFav) == 0 {
		DelFavDownloadTask()
	} else {
		AddFavDownloadTask()
	}
	return FlagResult{true, "更新配置成功"}
}

func (a *App) Delete(id int64) FlagResult {
	p := &Process{
		Id: id,
	}
	p.FindById()
	if p.Id == 0 {
		gefflog.Err("删除任务失败,任务不存在")
		return FlagResult{false, "删除任务失败,任务不存在"}
	}
	res, _ := YdpConfig.Mdb.Get(p.Id)
	if res != nil {
		p = res
	}
	if p.Progress.Status == StatusPending || p.Progress.Status == StatusDownloading {
		if err := p.Kill(); err != nil {
			gefflog.Err(fmt.Sprintf("删除任务失败 id: %d, err: %s", p.Id, err.Error()))
			return FlagResult{false, "删除任务失败"}
		}
	}
	p.Delete()
	return FlagResult{true, "删除成功"}
}

func MqNotifyConsumer() {
	err := YdpConfig.Mq.SetConsumer("notify", func(level, data string) {
		MainWin.EmitEvent("notify", false, level, data)
	})
	if err != nil {
		gefflog.Err(err.Error())
	}
}

func (a *App) CheckBiliLogin() FlagResult {
	ok, err := website.CheckLogin(YdpConfig.Cookies.Bilibili)
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	if !ok {
		return FlagResult{false, "账号未登录"}
	}
	return FlagResult{true, "账号已登录"}

}

func (a *App) GetFavList() FlagResultWithData {
	data, err := website.GetFavList(YdpConfig.Cookies.Bilibili)
	if err != nil {
		return FlagResultWithData{false, err.Error(), nil}
	}
	ps := make([]struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	}, 0, len(data.Data.List))
	for _, l := range data.Data.List {
		var p struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
		}
		p.Title = l.Title
		p.Id = l.Id
		ps = append(ps, p)

	}

	return FlagResultWithData{true, "操作成功", ps}

}
