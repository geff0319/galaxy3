package bridge

import (
	"context"
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
	Flag bool        `json:"flag"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
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
