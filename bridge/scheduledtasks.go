package bridge

import (
	"github.com/ge-fei-fan/gefflog"
	"github.com/geff0319/galaxy3/bridge/website"
	"github.com/robfig/cron/v3"
	"strconv"
	"strings"
)

var tasks cron.Cron
var favDownloadId cron.EntryID = -1001

func InitScheduledTasks() {
	tasks = *cron.New(cron.WithSeconds())
	tasks.Start()
	mediaId := strings.Split(YdpConfig.Other.BilibiliFav, ":")[0]
	if len(mediaId) != 0 {
		gefflog.Err("收藏夹：" + mediaId)
		AddFavDownloadTask()
	}
}

func (a *App) AddScheduledTask(spec string, event string) FlagResult {
	id, err := tasks.AddFunc(spec, func() {
		// log.Println("ScheduledTask: ", event)
		gefflog.Info("ScheduledTask: ", event)
		//runtime.EventsEmit(a.Ctx, event)
		MainWin.EmitEvent(event)
	})
	if err != nil {
		return FlagResult{false, err.Error()}
	}
	return FlagResult{true, strconv.Itoa(int(id))}
}

func (a *App) RemoveScheduledTask(id int) {
	gefflog.Info("RemoveScheduledTask:" + strconv.Itoa(id))
	tasks.Remove(cron.EntryID(id))
}

// b站下载定时任务
func AddFavDownloadTask() {
	DelFavDownloadTask()
	id, err := tasks.AddFunc("0 */10 * * * *", func() {
		mediaId := strings.Split(YdpConfig.Other.BilibiliFav, ":")[0]
		frr, err := website.GetFavResource(YdpConfig.Cookies.Bilibili, mediaId)
		if err != nil {
			gefflog.Err("GetFavResource err:" + err.Error())
		}
		for _, media := range frr.Data.Medias {
			gefflog.Info("开始下载：" + media.Title)
			p := &Process{
				Url:    "https://www.bilibili.com/video/" + media.Bvid,
				Params: []string{},
				Output: DownloadOutput{
					Path: YdpConfig.DownloadPath,
				},
				BiliMeta: website.BiliMetadata{
					SelectedVideoQuality: "",
				},
			}
			err = p.Insert()
			if err != nil {
				gefflog.Err("收藏夹下载视频<<" + media.Title + ">>失败:" + err.Error())
				continue
			}
			YdpConfig.Mq.Publish(p)

			//删除收藏夹文件
			cr, err := website.DelFavResource(YdpConfig.Cookies.Bilibili, media, mediaId, YdpConfig.Cookies.Bilijct)
			if err != nil {
				gefflog.Err("收藏夹删除视频<<" + media.Title + ">>失败:" + err.Error())
				continue
			}
			if cr.Code != 0 {
				gefflog.Err("收藏夹删除视频<<" + media.Title + ">>失败:" + cr.Message)
				continue
			}
			gefflog.Info("收藏夹删除视频<<" + media.Title + ">>成功")
		}
	})
	if err != nil {
		gefflog.Err("添加B站收藏夹定时任务失败：" + err.Error())
	}
	gefflog.Info("添加定时任务成功：" + strconv.Itoa(int(id)))
	favDownloadId = id
}

func DelFavDownloadTask() {
	defer func() {
		favDownloadId = -1001
	}()
	if favDownloadId != -1001 {
		tasks.Remove(favDownloadId)
		gefflog.Info("删除定时任务成功" + strconv.Itoa(int(favDownloadId)))
	}
}
