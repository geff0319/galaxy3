package bridge

import (
	"fmt"
	"galaxy3/bridge/website"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ge-fei-fan/gefflog"
)

const (
	BILIBIL_DOWNLOAD_TOPIC = "video/download/bilibili/%s"
	YOUTUBE_DOWNLOAD_TOPIC = "video/download/youtube/%s"
	TWITTER_DOWNLOAD_TOPIC = "video/download/twitter/%s"
	DOWNLOAD_RESULT_TOPIC  = "video/download/result"
)

var YoutubePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	gefflog.Info(fmt.Sprintf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))
	url, ok := website.NewYoutube(string(msg.Payload())).AppCompile()
	if ok {
		p := &Process{
			Url:    url,
			Params: []string{},
			Output: DownloadOutput{
				Path: YdpConfig.DownloadPath,
			},
			BiliMeta: &website.BiliMetadata{
				SelectedVideoQuality: "",
			},
		}
		YdpConfig.Mdb.Set(p)
		YdpConfig.Mq.Publish(p)
	}
}
var BilibiliPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	gefflog.Info(fmt.Sprintf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))
	url, ok := website.NewBlibili(string(msg.Payload())).AppCompile()
	if ok {
		p := &Process{
			Url:    url,
			Params: []string{},
			Output: DownloadOutput{
				Path: YdpConfig.DownloadPath,
			},
			BiliMeta: &website.BiliMetadata{
				SelectedVideoQuality: "",
			},
		}
		YdpConfig.Mdb.Set(p)
		YdpConfig.Mq.Publish(p)
	}
}

var TwitterPubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	gefflog.Info(fmt.Sprintf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))
	url, ok := website.NewTwitter(string(msg.Payload())).AppCompile()
	if ok {
		p := &Process{
			Url:    url,
			Params: []string{},
			Output: DownloadOutput{
				Path: YdpConfig.DownloadPath,
			},
			BiliMeta: &website.BiliMetadata{
				SelectedVideoQuality: "",
			},
		}
		YdpConfig.Mdb.Set(p)
		YdpConfig.Mq.Publish(p)
	}
}
