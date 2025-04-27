package bridge

import (
	"github.com/ge-fei-fan/clouddrive2api"
	"github.com/ge-fei-fan/gefflog"
)

var Cd2Client *clouddrive2api.Client

func InitCd2Client() {
	Cd2Client = clouddrive2api.NewClient(
		"192.168.2.80:19798",
		"754277710@qq.com",
		"987lxgff.",
		"/115open/云下载",
		"/115open/tg")

	err := Cd2Client.Login()
	if err != nil {
		gefflog.Err("初始化clouddrive2失败:" + err.Error())
		Cd2Client = nil
	}
}
