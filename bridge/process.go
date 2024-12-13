package bridge

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"github.com/geff0319/galaxy3/bridge/website"
	"github.com/geff0319/galaxy3/bridge/ytdlp"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const template = `download:
{
	"eta":%(progress.eta)s, 
	"percentage":"%(progress._percent_str)s",
	"speed":%(progress.speed)s
}`

const (
	StatusPending = iota
	StatusDownloading
	StatusCompleted
	StatusErrored
	StatusMergeing
)

// Process descriptor
type Process struct {
	Id       int64            `json:"id"`
	Pid      string           `json:"pid"`
	Url      string           `json:"url"`
	Params   []string         `json:"params"`
	Info     DownloadInfo     `json:"info"`
	Progress DownloadProgress `json:"progress"`
	Output   DownloadOutput   `json:"output"`
	proc     *os.Process
	BiliMeta website.BiliMetadata `json:"biliMeta"`
	//mux      sync.RWMutex
	//Logger   *slog.Logger
}

// Starts spawns/forks a new ytdlp process and parse its stdout.
// The process is spawned to outputting a custom progress text that
// Resembles a JSON Object in order to Unmarshal it later.
// This approach is anyhow not perfect: quotes are not escaped properly.
// Each process is not identified by its PID but by a UUIDv4
func (p *Process) Start() {
	YdpConfig.Mdb.Set(p)
	// escape bash variable escaping and command piping, you'll never know
	// what they might come with...
	p.Params = slices.DeleteFunc(p.Params, func(e string) bool {
		match, _ := regexp.MatchString(`(\$\{)|(\&\&)`, e)
		return match
	})

	p.Params = slices.DeleteFunc(p.Params, func(e string) bool {
		return e == ""
	})

	//out := DownloadOutput{
	//	Path:     filepath.Join(bridge.Env.BasePath, "video"),
	//	Filename: "%(title)s.%(ext)s",
	//}
	//
	//if p.Output.Path != "" {
	//	out.Path = p.Output.Path
	//}
	//
	//if p.Output.Filename != "" {
	//	out.Filename = p.Output.Filename
	//}
	//
	//buildFilename(&p.Output)

	//TODO: it spawn another one ytdlp process, too slow.
	//go p.GetFileName(&out)
	//p.Output.SavedFilePath = filepath.Join(YdpConfig.DownloadPath, sanitizeFileName(p.Info.FileName))
	p.Output.SavedFilePath = filepath.Join(p.Output.Path, ytdlp.SanitizeFileName(p.Info.FileName))
	// bilibii下载
	if strings.Contains(p.Url, "bilibili") {
		// 视频下载链接会过期，每次都重新获取，根据SelectedVideoQuality来查找历史的
		//tmpQuality := ""
		//if p.BiliMeta.SelectedVideoQuality != "" {
		//	tmpQuality = p.BiliMeta.SelectedVideoQuality
		//}
		err := p.SetMetadata()
		if err != nil {
			gefflog.Err(fmt.Sprintf("failed to Download bilibili process: err=%s", err.Error()))
			MainWin.EmitEvent("notify", false, "error", "下载bilibili视频失败"+err.Error())
			p.UpdateProgress(StatusErrored)
			return
		}
		p.BiliMeta.WriteFn = func(percentage string, speed float32) {
			if percentage == "100" {
				p.Progress = DownloadProgress{
					Status:     StatusMergeing,
					Percentage: percentage,
					Speed:      speed,
					ETA:        0,
				}
			} else {
				p.Progress = DownloadProgress{
					Status:     StatusDownloading,
					Percentage: percentage,
					Speed:      speed,
					ETA:        0,
				}
			}
			//p.UpdateProgress()
		}
		p.UpdateProgress(StatusDownloading)
		err = p.BiliMeta.Download(YdpConfig.BasePath)
		if err != nil {
			gefflog.Err(fmt.Sprintf("failed to Download bilibili process: err=%s", err.Error()))
			MainWin.EmitEvent("notify", false, "error", "下载bilibili视频失败"+err.Error())
			p.UpdateProgress(StatusErrored)
			p.BiliMeta.DoneChan <- struct{}{}
			close(p.BiliMeta.DoneChan)
			return
		}
		p.BiliMeta.DoneChan <- struct{}{}
		close(p.BiliMeta.DoneChan)
		p.Complete()
		return
	}

	baseParams := []string{
		strings.Split(p.Url, "?list")[0], //no playlist
		"--newline",
		"--no-colors",
		"--no-playlist",
		"--progress-template",
		strings.NewReplacer("\n", "", "\t", "", " ", "").Replace(template),
	}

	// if user asked to manually override the output path...
	if !(slices.Contains(p.Params, "-P") || slices.Contains(p.Params, "--paths")) {
		p.Params = append(p.Params, "-o")
		p.Params = append(p.Params, p.Output.SavedFilePath)
	}
	if strings.Contains(p.Url, "x.com") {
		if !ytdlp.IsFileExist(YdpConfig.BasePath + "/data/yt-dlp/cookies.txt") {
			gefflog.Err("下载X视频,cookies.txt不存在")
		}
		baseParams = append(baseParams, "--cookies", YdpConfig.BasePath+"/data/yt-dlp/cookies.txt")
	}

	params := append(baseParams, p.Params...)
	gefflog.Info(fmt.Sprintf("执行参数%s", params))

	// ----------------- main block ----------------- //
	if !IsYtDlpExist() {
		gefflog.Err("failed to start ytdlp process: err=ytdlp程序不存在")
		MainWin.EmitEvent("notify", false, "error", "启动任务失败:ytdlp程序不存在,请下载")
		p.UpdateProgress(StatusErrored)
		return
	}
	cmd := exec.Command(YdpConfig.YtDlpPath, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	r, err := cmd.StdoutPipe()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to connect to stdout: err=%s", err.Error()))
		MainWin.EmitEvent("notify", false, "error", "启动任务失败")
		p.UpdateProgress(StatusErrored)
		return
	}

	if err := cmd.Start(); err != nil {
		gefflog.Err(fmt.Sprintf("failed to start ytdlp process: err=%s", err.Error()))
		MainWin.EmitEvent("notify", false, "error", "启动任务失败")
		p.UpdateProgress(StatusErrored)
		return
	}

	p.proc = cmd.Process
	p.UpdateProgress(StatusDownloading)
	gefflog.Info(fmt.Sprintf("%s开始启动 pid：%d", p.Info.FileName, p.proc.Pid))
	// --------------- progress block --------------- //
	var (
		sourceChan = make(chan []byte)
		doneChan   = make(chan struct{})
	)

	// spawn a goroutine that does the dirty job of parsing the stdout
	// filling the channel with as many stdout line as ytdlp produces (producer)
	go func() {
		scan := bufio.NewScanner(r)

		defer func() {
			r.Close()
			p.Complete()

			doneChan <- struct{}{}

			close(sourceChan)
			close(doneChan)
		}()

		for scan.Scan() {
			sourceChan <- scan.Bytes()
		}
	}()

	// Slows down the unmarshal operation to every 500ms
	go func() {
		ytdlp.Sample(time.Millisecond*500, sourceChan, doneChan, func(event []byte) {
			var progress ProgressTemplate

			if err := json.Unmarshal(event, &progress); err != nil {
				return
			}
			p.Progress = DownloadProgress{
				Status:     StatusDownloading,
				Percentage: progress.Percentage,
				Speed:      progress.Speed,
				ETA:        progress.Eta,
			}
			//p.UpdateProgress()
		})
	}()

	// ------------- end progress block ------------- //
	cmd.Wait()
}

// Keep process in the memoryDB but marks it as complete
// Convention: All completed processes has progress -1
// and speed 0 bps.
func (p *Process) Complete() {
	p.Progress = DownloadProgress{
		Status:     StatusCompleted,
		Percentage: "-1",
		Speed:      0,
		ETA:        0,
	}
	p.UpdateProgress()
	if MqttC.Client != nil && MqttC.Client.IsConnectionOpen() {
		MqttC.Client.Publish(DOWNLOAD_RESULT_TOPIC, 0, false,
			fmt.Sprintf("[%s]%s: 下载完成", MqttC.opt.ClientID, p.Info.FileName))
	}
	gefflog.Info(fmt.Sprintf("finished: id=%s, url=%s", p.GetShortId(), p.Url))
	YdpConfig.Mdb.Delete(p.Pid)
}

// Kill a process and remove it from the memory
func (p *Process) Kill() error {

	defer func() {
		p.UpdateProgress(StatusCompleted)
	}()
	if strings.Contains(p.Url, "bilibili") {

		if p.BiliMeta.Cl != nil {
			p.BiliMeta.Cl()
		}
		return nil
	}
	// ytdlp uses multiple child process the parent process
	// has been spawned with setPgid = true. To properly kill
	// all subprocesses a SIGTERM need to be sent to the correct
	// process group
	if p.proc == nil {
		return errors.New("*os.Process not set")
	}
	p.proc.Kill()

	return nil
}

// Returns the available format for this URL
// TODO: Move out from process.go
func (p *Process) GetFormatsSync() (DownloadFormats, error) {
	cmd := exec.Command(YdpConfig.YtDlpPath, p.Url, "-J")

	stdout, err := cmd.Output()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to retrieve metadata: err=%s", err.Error()))
		//p.Logger.Error("failed to retrieve metadata", slog.String("err", err.Error()))
		return DownloadFormats{}, err
	}

	info := DownloadFormats{URL: p.Url}
	best := Format{}

	var (
		wg            sync.WaitGroup
		decodingError error
	)

	wg.Add(2)

	log.Println(
		ytdlp.BgRed, "Metadata", ytdlp.Reset,
		ytdlp.BgBlue, "Formats", ytdlp.Reset,
		p.Url,
	)
	gefflog.Info(fmt.Sprintf("retrieving metadata: caller=%s, url=%s", "getFormats", p.Url))
	//p.Logger.Info(
	//	"retrieving metadata",
	//	slog.String("caller", "getFormats"),
	//	slog.String("url", p.Url),
	//)

	go func() {
		decodingError = json.Unmarshal(stdout, &info)
		wg.Done()
	}()

	go func() {
		decodingError = json.Unmarshal(stdout, &best)
		wg.Done()
	}()

	wg.Wait()

	if decodingError != nil {
		return DownloadFormats{}, err
	}

	info.Best = best

	return info, nil
}

func (p *Process) GetFileName(o *DownloadOutput) error {
	cmd := exec.Command(
		YdpConfig.YtDlpPath,
		"--print", "filename",
		"-o", fmt.Sprintf("%s/%s", o.Path, o.Filename),
		p.Url,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	p.Output.SavedFilePath = strings.Trim(string(out), "\n")
	return nil
}

func (p *Process) SetPending() {
	// Since video's title isn't available yet, fill in with the URL.
	p.Info = DownloadInfo{
		URL:       p.Url,
		Title:     p.Url,
		CreatedAt: time.Now(),
	}
	p.UpdateProgress(StatusPending)
}

func (p *Process) SetMetadata() error {
	var result []string
	if strings.Contains(p.Url, "bilibili") {
		metadata, err := website.GetBilibiliInfo(p.Url, YdpConfig.Cookies.Bilibili)
		if err != nil {
			return err
		}

		info := DownloadInfo{
			Id:          strconv.FormatInt(metadata.Cr.Data.Cid, 10),
			URL:         p.Url,
			Title:       metadata.Cr.Data.Title,
			Thumbnail:   metadata.Cr.Data.Pic,
			Resolution:  strconv.Itoa(metadata.Vir.Data.Dash.Video[0].Height),
			Size:        int64(0),
			VCodec:      "",
			ACodec:      "",
			OriginalURL: p.Url,
			FileName:    metadata.Cr.Data.Title + ".mp4",
			CreatedAt:   time.Now(),
		}
		p.Info = info
		p.SetThumbnail()
		p.Pid = info.Id
		// 查询清晰度
		metadata.SelectedVideoQuality = p.BiliMeta.SelectedVideoQuality
		metadata.SelectedVideoCodecs = p.BiliMeta.SelectedVideoCodecs
		metadata.GetDefaultVideoStreamUrl()
		p.BiliMeta = metadata
		p.Progress.Status = StatusPending
		p.BiliMeta.SavedFilePath = p.Output.SavedFilePath
		p.Update()
		return nil
	}
	//检查ytdlp程序是否存在
	if !IsYtDlpExist() {
		return errors.New("tydlp程序不存在,请下载")
	}
	//cmd := exec.Command(YdpConfig.YtDlpPath, p.Url, "-J")
	baseParams := []string{
		strings.Split(p.Url, "?list")[0],
		"--dump-json",
		"--no-warnings",
	}
	if strings.Contains(p.Url, "x.com") {
		if !ytdlp.IsFileExist(YdpConfig.BasePath + "/data/yt-dlp/cookies.txt") {
			return errors.New("下载X视频,cookies.txt不存在")
		}
		baseParams = append(baseParams, "--cookies", YdpConfig.BasePath+"/data/yt-dlp/cookies.txt")
	}
	for _, str := range p.Params {
		if strings.TrimSpace(str) != "" {
			result = append(result, str)
		}
	}
	params := append(baseParams, result...)
	gefflog.Info(fmt.Sprintf("执行参数%s", params))
	cmd := exec.Command(YdpConfig.YtDlpPath, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to connect to stdout: id=%s, url=%s, err=%s", p.GetShortId(), p.Url, err.Error()))
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		gefflog.Err(fmt.Sprintf("failed to connect to stderr: id=%s, url=%s, err=%s", p.GetShortId(), p.Url, err.Error()))
		return err
	}

	info := DownloadInfo{
		URL:       p.Url,
		CreatedAt: time.Now(),
	}

	if err := cmd.Start(); err != nil {
		gefflog.Err(fmt.Sprintf("failed to start cmd: id=%s, url=%s, err=%s", p.GetShortId(), p.Url, err.Error()))
		return err
	}

	var bufferedStderr bytes.Buffer

	go func() {
		io.Copy(&bufferedStderr, stderr)
	}()
	gefflog.Info(fmt.Sprintf("retrieving metadata: id=%s, url=%s", p.GetShortId(), p.Url))

	if err := json.NewDecoder(stdout).Decode(&info); err != nil {
		gefflog.Err(fmt.Sprintf("failed to Decode json : id=%s, url=%s, err=%s", p.GetShortId(), p.Url, err.Error()))
		gefflog.Err(bufferedStderr.String())
		return errors.New(bufferedStderr.String())
	}

	if err := cmd.Wait(); err != nil {
		gefflog.Err(fmt.Sprintf("failed to wait cmd: id=%s, url=%s, err=%s", p.GetShortId(), p.Url, err.Error()))
		return errors.New(bufferedStderr.String())
	}
	gefflog.Info(info)
	p.Info = info
	p.SetThumbnail()
	p.Progress.Status = StatusPending
	p.Pid = p.Info.Id
	p.Update()
	return nil
}

func (p *Process) SetThumbnail() {
	// 下载图片
	if p.Info.Thumbnail == "" {
		return
	}
	resp, err := http.Get(p.Info.Thumbnail)
	if err != nil {
		gefflog.Err("SetThumbnail err:" + err.Error())
		p.Info.Thumbnail = ""
		return
	}
	defer resp.Body.Close()

	// 读取图片内容
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		gefflog.Err("SetThumbnail err:" + err.Error())
		p.Info.Thumbnail = ""
		return
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		gefflog.Err("SetThumbnail err:URL does not point to an image" + p.Info.Thumbnail)
	}

	// 将图片内容转换为 Base64 编码
	base64Str := base64.StdEncoding.EncodeToString(imgData)
	p.Info.Thumbnail = fmt.Sprintf("data:%s;base64,%s", contentType, base64Str)

}

func (p *Process) GetShortId() string {
	//return strings.Split(p.Pid, "-")[0]
	return strconv.FormatInt(p.Id, 10)
}

func buildFilename(o *DownloadOutput) {
	if o.Filename != "" && strings.Contains(o.Filename, ".%(ext)s") {
		o.Filename += ".%(ext)s"
	}

	o.Filename = strings.Replace(
		o.Filename,
		".%(ext)s.%(ext)s",
		".%(ext)s",
		1,
	)
}

func (p *Process) Insert() error {
	//p.mux.Lock()
	//defer p.mux.Unlock()
	pInfo, _ := json.Marshal(p.Info)
	pProgress, _ := json.Marshal(p.Progress)
	pOutput, _ := json.Marshal(p.Output)
	pBiliMeta, _ := json.Marshal(p.BiliMeta)
	res, err := SqliteS.Execute(ProcessInsert, p.Url, strings.Join(p.Params, ","), string(pInfo), string(pProgress), string(pOutput), string(pBiliMeta))
	if err != nil {
		gefflog.Err("Process Insert error: " + err.Error())
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		gefflog.Err("Process Insert error: " + err.Error())
		return err
	}
	p.Id = id
	return nil
}
func (p *Process) Update() {
	//p.mux.Lock()
	//defer p.mux.Unlock()
	gefflog.Info("更新数据:" + strconv.FormatInt(p.Id, 10))
	pInfo, _ := json.Marshal(p.Info)
	pProgress, _ := json.Marshal(p.Progress)
	pOutput, _ := json.Marshal(p.Output)
	pBiliMeta, _ := json.Marshal(p.BiliMeta)
	_, err := SqliteS.Execute(ProcessUpdate, p.Info.Id, string(pInfo), string(pProgress), string(pOutput), string(pBiliMeta), p.Id)
	if err != nil {
		gefflog.Err("Process Update error: " + err.Error())
	}
}
func (p *Process) UpdateProgress(args ...int) {
	//p.mux.Lock()
	//defer p.mux.Unlock()
	if len(args) != 0 {
		p.Progress.Status = args[0]
	}
	pProgress, _ := json.Marshal(p.Progress)
	_, err := SqliteS.Execute(ProcessUpdateProgress, string(pProgress), p.Id)
	if err != nil {
		gefflog.Err("Process UpdateProgress error: " + err.Error())
	}
}
func (p *Process) Delete() {
	//p.mux.Lock()
	//defer p.mux.Unlock()
	_, err := SqliteS.Execute(ProcessDelete, 1, p.Id)
	if err != nil {
		gefflog.Err("Process Delete error: " + err.Error())
	}
	YdpConfig.Mdb.Delete(p.Pid)
}
func (p *Process) IsExist() int {
	gefflog.Info(p.Pid)
	res, err := SqliteS.Select("SELECT id FROM process WHERE pid = ? AND is_delete = 0;", p.Pid)
	gefflog.Info(fmt.Sprintf("文件个数：%d", len(res)))
	if err != nil {
		gefflog.Err("IsExist err:" + err.Error())
		return -1
	}
	return len(res)
}
func (p *Process) FindById() {
	res, err := SqliteS.Select("SELECT pid,url, params,info,progress,output FROM process where is_delete = 0 and id = ?", p.Id)
	if err != nil {
		p.Id = 0
		gefflog.Err("FindById err:" + err.Error())
		return
	}
	if len(res) == 0 {
		p.Id = 0
		return
	}
	p.Unmarshal(res[0])
}
func (p *Process) Unmarshal(dst map[string]any) {
	if value, ok := dst["id"]; ok {
		t, ok := value.(int64)
		if ok {
			p.Id = t
		} else {
			p.Id = 0
		}
	}
	if value, ok := dst["pid"]; ok {
		t, ok := value.(string)
		if ok {
			p.Pid = t
		} else {
			p.Pid = ""
		}
	}
	if value, ok := dst["url"]; ok {
		t, ok := value.(string)
		if ok {
			p.Url = t
		} else {
			p.Url = ""
		}
	}
	if value, ok := dst["progress"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &p.Progress); err != nil {
			gefflog.Err("ProcessResponse Unmarshal progress err:" + err.Error())
		}
	}
	if value, ok := dst["info"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &p.Info); err != nil {
			gefflog.Err("ProcessResponse Unmarshal info err:" + err.Error())
		}
	}
	if value, ok := dst["output"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &p.Output); err != nil {
			gefflog.Err("ProcessResponse Unmarshal output err:" + err.Error())
		}
	}
	if value, ok := dst["params"]; ok {
		t, ok := value.(string)
		//if ok {
		//	p.Params = strings.Split(t, ",")
		//} else {
		//	p.Params = []string{}
		//}
		if ok {
			p.Params = strings.Split(t, ",")
		}
	}
	if value, ok := dst["biliMeta"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &p.BiliMeta); err != nil {
			gefflog.Err("ProcessResponse Unmarshal biliMeta err:" + err.Error())
		}
	}
}
