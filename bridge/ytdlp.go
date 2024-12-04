package bridge

import (
	"encoding/json"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"github.com/geff0319/galaxy3/bridge/website"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
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
			//YdpConfig.Mdb.Set(&process)
			err := process.Insert()
			if err != nil {
				return FlagResult{true, "下载视频失败:" + err.Error()}
			}
			YdpConfig.Mq.PublishByTopic("process:downloading", &process)
			return FlagResult{true, strconv.FormatInt(process.Id, 10)}
		}
	} else {
		if process.Info.Id == "" {
			gefflog.Info("重试解析: " + strconv.FormatInt(process.Id, 10))
			YdpConfig.Mq.Publish(&process)
		} else {
			gefflog.Info("重试下载: " + strconv.FormatInt(process.Id, 10))
			process.SetPending()
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
	//id := YdpConfig.Mdb.Set(p)
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
	//res := YdpConfig.Mdb.All()
	res, err := SqliteS.Select(ProcessAll)
	if err != nil {
		gefflog.Err("SelectProcess err:" + err.Error())
	}
	var p Process
	ps := []Process{}
	for _, r := range res {
		p.Unmarshal(r)
		ps = append(ps, p)
	}
	return FlagResultWithData{
		Flag: true,
		Msg:  "操作成功",
		Data: ps,
	}
}

func (a *App) UpdateYtDlpConfig() FlagResult {
	b, err := os.ReadFile(Env.BasePath + "/data/ytdlp.yaml")
	if err != nil {
		gefflog.Err("更新配置失败: " + err.Error())
		return FlagResult{false, "更新配置失败"}
	}
	if err := yaml.Unmarshal(b, &YdpConfig); err != nil {
		gefflog.Err("更新配置失败: " + err.Error())
		return FlagResult{false, "更新配置失败"}
	}
	return FlagResult{true, "更新配置成功"}
}

func (a *App) Delete(id int64) FlagResult {
	var p Process

	//gefflog.Info("删除任务 id:" + strconv.FormatInt(id, 10))
	//p, err := YdpConfig.Mdb.Get(id)
	//if err != nil {
	//	gefflog.Err(fmt.Sprintf("删除任务失败 id: %s, err: %s", id, err.Error()))
	//	return FlagResult{false, "删除任务失败:"}
	//}
	//if p == nil {
	//	return FlagResult{false, "删除任务失败: nil process"}
	//}

	//p.Id, _ = strconv.ParseInt(id, 10, 64)
	p.Id = id
	p.FindById()
	if p.Id == 0 {
		gefflog.Err("删除任务失败,任务不存在")
		return FlagResult{false, "删除任务失败,任务不存在"}
	}
	if p.Progress.Status == StatusPending || p.Progress.Status == StatusDownloading {
		if err := p.Kill(); err != nil {
			gefflog.Err(fmt.Sprintf("删除任务失败 id: %d, err: %s", p.Id, err.Error()))
			return FlagResult{false, "删除任务失败"}
		}
		//YdpConfig.Mdb.Delete(p.Id)

	} else {
		//YdpConfig.Mdb.Delete(p.Id)
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
