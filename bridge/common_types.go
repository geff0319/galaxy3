package bridge

import (
	"encoding/json"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"github.com/geff0319/galaxy3/bridge/website"
	"github.com/geff0319/galaxy3/bridge/ytdlp"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Used to unmarshall ytdlp progress
type ProgressTemplate struct {
	Percentage string  `json:"percentage"`
	Speed      float32 `json:"speed"`
	Size       string  `json:"size"`
	Eta        float32 `json:"eta"`
}

// Defines where and how the download needs to be saved
type DownloadOutput struct {
	Path          string
	Filename      string
	SavedFilePath string `json:"savedFilePath"`
}

// Progress for the Running call
type DownloadProgress struct {
	Status     int     `json:"process_status"`
	Percentage string  `json:"percentage"`
	Speed      float32 `json:"speed"`
	ETA        float32 `json:"eta"`
}

// Used to deser the ytdlp -J output
type DownloadInfo struct {
	Id          string    `json:"Id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Thumbnail   string    `json:"thumbnail"`
	Resolution  string    `json:"resolution"`
	Size        int64     `json:"filesize_approx"`
	VCodec      string    `json:"vcodec"`
	ACodec      string    `json:"acodec"`
	Extension   string    `json:"ext"`
	OriginalURL string    `json:"original_url"`
	FileName    string    `json:"filename"`
	CreatedAt   time.Time `json:"created_at"`
	//RequestedDownloads struct {
	//	Filename string `json:"filename"`
	//} `json:"requested_downloads"`
}

// Used to deser the formats in the -J output
type DownloadFormats struct {
	Formats   []Format `json:"formats"`
	Best      Format   `json:"best"`
	Thumbnail string   `json:"thumbnail"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
}

// A skimmed ytdlp format node
type Format struct {
	Format_id   string  `json:"format_id"`
	Format_note string  `json:"format_note"`
	FPS         float32 `json:"fps"`
	Resolution  string  `json:"resolution"`
	VCodec      string  `json:"vcodec"`
	ACodec      string  `json:"acodec"`
	Size        float32 `json:"filesize_approx"`
}

// struct representing the response sent to the client
// as JSON-RPC result field
type ProcessResponse struct {
	Id       int64                `json:"id"`
	Pid      string               `json:"pid"`
	Url      string               `json:"url"`
	Progress DownloadProgress     `json:"progress"`
	Info     DownloadInfo         `json:"info"`
	Output   DownloadOutput       `json:"output"`
	Params   []string             `json:"params"`
	BiliMeta website.BiliMetadata `json:"biliMeta"`
}

func (pr *ProcessResponse) Unmarshal(dst map[string]any) {
	if value, ok := dst["id"]; ok {
		pr.Id = value.(int64)
	}
	if value, ok := dst["pid"]; ok {
		pr.Pid = value.(string)
	}
	if value, ok := dst["url"]; ok {
		pr.Url = value.(string)
	}
	if value, ok := dst["progress"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &pr.Progress); err != nil {
			gefflog.Err("ProcessResponse Unmarshal progress err:" + err.Error())
		}
	}
	if value, ok := dst["info"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &pr.Info); err != nil {
			gefflog.Err("ProcessResponse Unmarshal info err:" + err.Error())
		}
	}
	if value, ok := dst["output"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &pr.Output); err != nil {
			gefflog.Err("ProcessResponse Unmarshal output err:" + err.Error())
		}
	}
	if value, ok := dst["params"]; ok {
		pr.Params = strings.Split(value.(string), ",")
	}
	if value, ok := dst["biliMeta"]; ok {
		if err := json.Unmarshal([]byte(value.(string)), &pr.BiliMeta); err != nil {
			gefflog.Err("ProcessResponse Unmarshal biliMeta err:" + err.Error())
		}
	}
}

// struct representing the current status of the memoryDB
// used for serializaton/persistence reasons
type Session struct {
	Processes []ProcessResponse `json:"processes"`
}

// struct representing the intent to stop a specific process
type AbortRequest struct {
	Id string `json:"id"`
}

// struct representing the intent to start a download
type DownloadRequest struct {
	Id     string
	URL    string   `json:"url"`
	Path   string   `json:"path"`
	Rename string   `json:"rename"`
	Params []string `json:"params"`
}

// struct representing request of creating a netscape cookies file
type SetCookiesRequest struct {
	Cookies string `json:"cookies"`
}

// represents a user defined collection of ytdlp arguments
type CustomTemplate struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type YtDlpCookie struct {
	Bilibili string `yaml:"bilibili"`
}
type YtDlpConfig struct {
	BasePath     string
	DownloadPath string        `yaml:"downloadPath"` //视频保存路径
	YtDlpPath    string        `yaml:"ytDlpPath"`    //下载程序路径
	QueueSize    int           `yaml:"queueSize"`
	Mdb          *MemoryDB     `yaml:"-"`
	Mq           *MessageQueue `yaml:"-"`
	Cookies      YtDlpCookie   `yaml:"cookies"`
}

func (yc *YtDlpConfig) Unmarshal() error {
	b, err := SqliteS.Select(ConfigYtdlp)
	if err != nil {
		gefflog.Err("Select YtDlpConfig err:" + err.Error())
		return err
	}
	if len(b) == 0 {

		marshal, err := yaml.Marshal(yc)
		if err != nil {
			gefflog.Err("Marshal YtDlpConfig err:" + err.Error())
			return err
		}
		fmt.Println(string(marshal))
		_, err = SqliteS.Execute(ConfigInsert, "ytdlp", string(marshal))
		if err != nil {
			gefflog.Err("YtDlpConfig ConfigInsert err:" + err.Error())
		}
		return err
	} else {
		if value, ok := b[0]["config_value"]; ok {
			t, ok := value.(string)
			if ok {
				err := yaml.Unmarshal([]byte(t), &yc)
				if err != nil {
					gefflog.Err("YtDlpConfig Unmarshal err:" + err.Error())
				}
				return err
			} else {
				marshal, err := yaml.Marshal(yc)
				if err != nil {
					gefflog.Err("Marshal YtDlpConfig err:" + err.Error())
					return err
				}
				_, err = SqliteS.Execute(ConfigUpdate, string(marshal), "ytdlp")
				if err != nil {
					gefflog.Err("YtDlpConfig ConfigUpdate err:" + err.Error())
				}
			}
		}
		return err
	}
}

var YdpConfig YtDlpConfig

func InitYtDlpConfig() {
	var mdb MemoryDB
	YdpConfig = YtDlpConfig{
		BasePath:     Env.BasePath,
		DownloadPath: Env.BasePath + "/data/yt-dlp-download/",
		YtDlpPath:    Env.BasePath + "/data/yt-dlp/yt-dlp.exe",
		QueueSize:    8,
		Mdb:          &mdb,
		Cookies:      YtDlpCookie{},
	}
	if !ytdlp.IsDirExists(YdpConfig.DownloadPath) {
		os.MkdirAll(filepath.Dir(YdpConfig.DownloadPath), os.ModePerm)
	}
	//b, err := os.ReadFile(basePath + "/data/ytdlp.yaml")
	//if os.IsNotExist(err) {
	//	os.MkdirAll(filepath.Dir(basePath+"/data/ytdlp.yaml"), os.ModePerm)
	//	content, err := yaml.Marshal(YdpConfig)
	//	if err != nil {
	//		gefflog.Err("生成ytdlp.yaml配置失败")
	//	}
	//	if err = os.WriteFile(basePath+"/data/ytdlp.yaml", content, 0644); err != nil {
	//		gefflog.Err("生成ytdlp.yaml配置失败")
	//	}
	//} else {
	//	if err := yaml.Unmarshal(b, &YdpConfig); err != nil {
	//		gefflog.Err("获取ytdlp.yaml配置失败")
	//	}
	//}

	YdpConfig.Unmarshal()
	fmt.Println(fmt.Sprintf("YdpConfig:%v+", YdpConfig))
	mq, err := NewMessageQueue()
	if err != nil {
		gefflog.Err("初始化下载消息队列失败：" + err.Error())
		return
	}
	YdpConfig.Mq = mq
	YdpConfig.Mq.SetupConsumers()

	//go YdpConfig.Mdb.Restore(basePath, mq)
	//YdpConfig.Mdb.Restore(basePath, mq)
}
