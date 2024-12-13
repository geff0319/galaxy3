package bridge

import (
	"context"
	"github.com/ge-fei-fan/gefflog"
	"gopkg.in/yaml.v3"
	"net/http"
)

// App struct
type App struct {
	Ctx context.Context
}

type EnvResult struct {
	FromTaskSch bool   `json:"-"`
	AppName     string `json:"appName"`
	BasePath    string `json:"basePath"`
	OS          string `json:"os"`
	ARCH        string `json:"arch"`
	X64Level    int    `json:"x64Level"`
}

type ExecOptions struct {
	Convert bool              `json:"convert"`
	Env     map[string]string `json:"env"`
}

type FlagResult struct {
	Flag bool   `json:"flag"`
	Data string `json:"data"`
}

type FlagResultWithData struct {
	Flag bool   `json:"flag"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type HTTPResult struct {
	Flag   bool        `json:"flag"`
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
}

type AppConfig struct {
	WindowStartState int       `yaml:"windowStartState"`
	UserAgent        string    `yaml:"userAgent"`
	ExitOnClose      bool      `yaml:"exitOnClose"`
	Translate        Translate `yaml:"translate"`
	LogPath          string    `yaml:"logPath"`
}

type MenuItem struct {
	Type     string     `json:"type"` // Menu Type: item / separator
	Text     string     `json:"text"`
	Tooltip  string     `json:"tooltip"`
	Event    string     `json:"event"`
	Children []MenuItem `json:"children"`
	Hidden   bool       `json:"hidden"`
	Checked  bool       `json:"checked"`
}

type TrayContent struct {
	Icon    string `json:"icon"`
	Title   string `json:"title"`
	Tooltip string `json:"tooltip"`
}

type WriteTracker struct {
	Total          int64
	Progress       int64
	ProgressChange string
	App            *App
}

type Translate struct {
	TencentTanslateSecretId  string `yaml:"tencentTanslateSecretId"`
	TencentTanslateSecretKey string `yaml:"tencentTanslateSecretKey"`
}

// 获取视频信息
type FormatsMeta struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type YoutubeVideoMeta struct {
	Title            string        `json:"title"`
	Formats          []FormatsMeta `json:"formats"`
	Filename         string        `json:"filename"`
	RequestedFormats []FormatsMeta `json:"requested_formats"`
}

func (ac *AppConfig) Unmarshal() {
	b, err := SqliteS.Select(ConfigUser)
	if err != nil || len(b) == 0 || b == nil {
		ac.WindowStartState = 0
		ac.ExitOnClose = true
		ac.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
		ac.LogPath = Env.BasePath + "/logs"
		return
	} else {
		if value, ok := b[0]["config_value"]; ok {
			t, ok := value.(string)
			if ok {
				err := yaml.Unmarshal([]byte(t), ac)
				if err != nil {
					gefflog.Err("AppConfig Unmarshal err:" + err.Error())
					ac.WindowStartState = 0
					ac.ExitOnClose = true
					ac.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
					ac.LogPath = Env.BasePath + "/logs"
					return
				}
				if ac.LogPath == "" {
					ac.LogPath = Env.BasePath + "/logs"
				}
			} else {
				ac.WindowStartState = 0
				ac.ExitOnClose = true
				ac.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
				ac.LogPath = Env.BasePath + "/logs"
			}
		}
	}
}
