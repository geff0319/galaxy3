package website

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type progressWriter struct {
	totalBytes   int64
	currentBytes int64
	writer       io.Writer
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	if err != nil {
		return n, err
	}
	pw.currentBytes += int64(n)
	return n, nil
}

type Bilibili struct {
	sourceUrl string
}
type CidResponse struct {
	Data struct {
		Bvid  string `json:"bvid"`
		Cid   int64  `json:"cid"`
		Pic   string `json:"pic"`
		Title string `json:"title"`
	} `json:"data"`
}
type VideoInfoResponse struct {
	Data struct {
		AcceptDescription []string `json:"accept_description"`
		AcceptQuality     []int64  `json:"accept_quality"`
		SupportFormats    []struct {
			Quality        int      `json:"quality"`
			NewDescription string   `json:"new_description"`
			Codecs         []string `json:"codecs"`
		} `json:"support_formats"`
		Dash struct {
			Video []struct {
				Id        int      `json:"id"`
				BaseUrl   string   `json:"base_url"`
				BackupUrl []string `json:"backupUrl"`
				Bandwidth int64    `json:"bandwidth"`
				Codecs    string   `json:"codecs"`
				Width     int      `json:"width"`
				Height    int      `json:"height"`
			} `json:"video"`
			Audio []struct {
				BaseUrl   string `json:"base_url"`
				Bandwidth int64  `json:"bandwidth"`
			} `json:"audio"`
		} `json:"dash"`
	} `json:"data"`
}
type BiliMetadata struct {
	ctx                    context.Context
	Cl                     context.CancelFunc `json:"-"`
	DoneChan               chan struct{}      `json:"-"`
	SavedFilePath          string
	SelectedVideoStreamUrl string
	SelectedVideoQuality   string
	SelectedVideoCodecs    string
	Cr                     CidResponse
	Vir                    VideoInfoResponse
	WriteFn                func(string, float32) `json:"-"`
	pWriter                *progressWriter
}
type BiliMetaMemory struct {
	SavedFilePath          string
	SelectedVideoStreamUrl string
	SelectedVideoQuality   string
	SelectedVideoCodecs    string
	Cr                     CidResponse
	Vir                    VideoInfoResponse
}

type WebInterfaceNav struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		IsLogin bool   `json:"isLogin"`
		Uname   string `json:"uname"`
		Face    string `json:"face"`
		Mid     int    `json:"mid"`
	} `json:"data"`
}

func NewBlibili(url string) *Bilibili {
	return &Bilibili{
		sourceUrl: url,
	}
}

func (b *Bilibili) Compile() (string, bool) {
	return "", false
}

func (b *Bilibili) AppCompile() (string, bool) {
	// 定义正则表达式来匹配 HTTP 链接
	re := regexp.MustCompile(`https?://[^\s]+`)

	// 查找匹配的链接
	Links := re.FindAllString(b.sourceUrl, -1)

	// 输出所有匹配的链接
	if len(Links) != 0 {
		req, err := http.NewRequest(http.MethodGet, Links[0], nil)
		if err != nil {
			gefflog.Err("解析b站分享链接失败：" + err.Error())
			return "", false
		}
		req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
		res, err := client.Do(req)
		if err != nil {
			gefflog.Err("解析b站分享链接失败：" + err.Error())
			return "", false
		}
		defer res.Body.Close()
		if res.Request.URL != nil {
			return res.Request.URL.String(), true
		}
		return "", false
	}
	return "", false
}

var client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		//禁止复用连接，防止同一个连接长时间大流量被限速
		DisableKeepAlives: true,
	},
}

func getCid(bv string) ([]byte, error) {
	url := "https://api.bilibili.com/x/web-interface/view?bvid=" + bv
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//log.Println(string(body))
	return body, nil
}

func getStream(bv, ck string, cid int64) (VideoInfoResponse, error) {
	var vif VideoInfoResponse
	url := fmt.Sprintf("https://api.bilibili.com/x/player/wbi/playurl?fnver=0&fnval=3216&fourk=1&qn=127&bvid=%s&cid=%d", bv, cid)
	fmt.Println(url)
	var err error
	url, err = sign(url)
	if err != nil {
		return vif, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return vif, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	if ck != "" {
		req.AddCookie(&http.Cookie{Name: "SESSDATA", Value: ck})
	}

	res, err := client.Do(req)
	if err != nil {
		return vif, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return vif, err
	}

	err = json.Unmarshal(body, &vif)
	if err != nil {
		return vif, err
	}

	return vif, nil
}

func GetBilibiliInfo(url, ck string) (BiliMetadata, error) {
	var md BiliMetadata
	u, err := url2.Parse(url)
	if err != nil {
		return md, err
	}
	segments := strings.Split(u.Path, "/")
	if len(segments) < 2 {
		return md, errors.New("获取bv失败")
	}
	bvid := segments[2]
	var cr CidResponse
	res, err := getCid(bvid)
	if err != nil {
		return md, errors.New("getCid err: " + err.Error())
	}
	err = json.Unmarshal(res, &cr)
	if err != nil {
		return md, err
	}
	infoResp, err := getStream(cr.Data.Bvid, ck, cr.Data.Cid)
	if err != nil {
		return md, errors.New("getStream err: " + err.Error())
	}
	md.Cr = cr
	md.Vir = infoResp

	return md, nil
}

func CheckLogin(ck string) (bool, error) {
	navUrl := "https://api.bilibili.com/x/web-interface/nav"
	navUrl, err := sign(navUrl)
	req, err := http.NewRequest(http.MethodGet, navUrl, nil)
	if err != nil {
		gefflog.Err("CheckLogin NewRequest err: " + err.Error())
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	if ck != "" {
		req.AddCookie(&http.Cookie{Name: "SESSDATA", Value: ck})
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		gefflog.Err("CheckLogin Request err: %s" + err.Error())
		return false, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		gefflog.Err("CheckLogin ReadAll err: %s" + err.Error())
		return false, err
	}
	var win WebInterfaceNav
	err = json.Unmarshal(body, &win)
	if err != nil {
		gefflog.Err("CheckLogin json Unmarshal err: %s" + err.Error())
		return false, err
	}
	if win.Code != 0 {
		return false, errors.New(win.Message)
	}
	return win.Data.IsLogin, nil
}

func (bmd *BiliMetadata) GetDefaultVideoStreamUrl() {
	//根据清晰度下载视频，此时编码就随机吧
	tmpId := -1
	if bmd.SelectedVideoQuality != "" {
		bmd.SelectedVideoStreamUrl = ""
		for _, formats := range bmd.Vir.Data.SupportFormats {
			if formats.NewDescription == bmd.SelectedVideoQuality {
				tmpId = formats.Quality
				break
			}
		}
		// -1就是没找到清晰度，下去选默认清晰度
		if tmpId != -1 {
			if bmd.SelectedVideoCodecs != "" {
				for _, elem := range bmd.Vir.Data.Dash.Video {
					if elem.Id == tmpId && elem.Codecs == bmd.SelectedVideoCodecs {
						bmd.SelectedVideoStreamUrl = elem.BaseUrl
						break
					}
				}
			} else { //无编码就随机编码吧
				for _, elem := range bmd.Vir.Data.Dash.Video {
					if elem.Id == tmpId {
						bmd.SelectedVideoStreamUrl = elem.BaseUrl
						break
					}
				}
			}

		}
		if bmd.SelectedVideoStreamUrl != "" {
			return
		}
	}

	//获取默认清晰度
	id := bmd.Vir.Data.Dash.Video[0].Id
	var dstUrl string
	for _, elem := range bmd.Vir.Data.Dash.Video {
		if elem.Id == id {
			dstUrl = elem.BaseUrl
			if strings.Contains(elem.Codecs, "avc1") {
				break
			}
		}
	}
	bmd.SelectedVideoStreamUrl = dstUrl
	for _, formats := range bmd.Vir.Data.SupportFormats {
		if formats.Quality == id {
			bmd.SelectedVideoQuality = formats.NewDescription
			break
		}
	}
}

func (bmd *BiliMetadata) Download(basePath string) error {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	bmd.ctx = ctx
	bmd.Cl = cancel
	bmd.DoneChan = make(chan struct{})
	err = bmd.DA()
	if err != nil {
		if errors.Is(err, bmd.ctx.Err()) {
			return nil
		}
		return err
	}
	err = bmd.DV()
	if err != nil {
		if errors.Is(err, bmd.ctx.Err()) {
			return nil
		}
		return err
	}
	err = bmd.Merge(basePath)
	if err != nil {
		return err
	}
	return nil
}

func (bmd *BiliMetadata) speed() {
	var lastBytes int64

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			speed := bmd.pWriter.currentBytes - lastBytes
			lastBytes = bmd.pWriter.currentBytes
			percentage := float64(bmd.pWriter.currentBytes) / float64(bmd.pWriter.totalBytes) * 100
			bmd.WriteFn(strconv.Itoa(int(percentage)), float32(speed))
		case <-bmd.DoneChan:
			ticker.Stop()
			gefflog.Info("speed 退出")
			//bmd.WriteFn("0",0)
			return
		}
	}
}

func (bmd *BiliMetadata) DV() error {
	//bmd.GetDefaultVideoStreamUrl()
	req, err := http.NewRequestWithContext(bmd.ctx, http.MethodGet, bmd.SelectedVideoStreamUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://www.bilibili.com")
	resp, err := client.Do(req)
	if err != nil {
		gefflog.Err("client request err: " + err.Error())
		return err
	}
	var file *os.File
	file, err = os.OpenFile(bmd.SavedFilePath+".video", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		gefflog.Err("OpenFile video err: " + err.Error())
		return err
	}
	defer file.Close()
	_ = file.Truncate(0)
	pw := progressWriter{
		totalBytes:   resp.ContentLength,
		currentBytes: int64(0),
		writer:       file,
	}
	fmt.Println(resp.ContentLength)
	bmd.pWriter = &pw
	go bmd.speed()
	//defer func() {
	//	bmd.DoneChan <- struct{}{}
	//	close(bmd.DoneChan)
	//}()
	_, err = io.Copy(bmd.pWriter, resp.Body)
	if err != nil {
		gefflog.Err("CopyFile video err: " + err.Error())
		return err
	}
	return nil
}

func (bmd *BiliMetadata) DA() error {
	req, err := http.NewRequestWithContext(bmd.ctx, http.MethodGet, bmd.Vir.Data.Dash.Audio[0].BaseUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://www.bilibili.com")
	resp, err := client.Do(req)
	if err != nil {
		gefflog.Err("client request err: " + err.Error())
		return err
	}
	var file *os.File
	file, err = os.OpenFile(bmd.SavedFilePath+".audio", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		gefflog.Err("OpenFile audio err: " + err.Error())
		return err
	}
	defer file.Close()
	_ = file.Truncate(0)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		gefflog.Err("CopyFile audio err: " + err.Error())
		return err
	}
	return nil
}

func (bmd *BiliMetadata) Merge(basePath string) error {

	video := bmd.SavedFilePath + ".audio"
	audio := bmd.SavedFilePath + ".video"
	ffPath := basePath + "/data/yt-dlp/ffmpeg.exe"
	cmd := exec.Command(ffPath, "-y", "-i", video, "-i", audio, "-c", "copy", bmd.SavedFilePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		gefflog.Err("Merge error: " + err.Error())
		return errors.New("Merge error: " + err.Error())
	}
	err = os.Remove(video)
	if err != nil {
		gefflog.Err("Remove video error: " + err.Error())
		return errors.New(video + "Remove video error: " + err.Error())
	}
	err = os.Remove(audio)
	if err != nil {
		gefflog.Err("Remove audio error: " + err.Error())
		return errors.New(audio + "Remove audio error: " + err.Error())
	}
	return nil
}
