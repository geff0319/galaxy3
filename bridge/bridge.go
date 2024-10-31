package bridge

import (
	_ "embed"
	"galaxy3/bridge/ytdlp"
	"github.com/klauspost/cpuid/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
	"io/fs"

	//"github.com/wailsapp/wails/v2/pkg/options"
	//R "github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed image/icon.ico
var Icon []byte

var MainWin *application.WebviewWindow
var WidgetsWin *application.WebviewWindow
var MainApp *application.App

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

//	func (a *App) OnSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
//		R.WindowUnminimise(a.Ctx)
//		R.Show(a.Ctx)
//		go R.EventsEmit(a.Ctx, "launchArgs", secondInstanceData.Args)
//	}

var Env = &EnvResult{
	BasePath:    "",
	AppName:     "",
	OS:          runtime.GOOS,
	ARCH:        runtime.GOARCH,
	X64Level:    cpuid.CPU.X64Level(),
	FromTaskSch: false,
}

var Config = &AppConfig{}
var AppTitle = "Galaxy"

func InitBridge(assets fs.FS) {
	// step1: Set Env
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	for _, v := range os.Args {
		if v == "tasksch" {
			Env.FromTaskSch = true
			break
		}
	}
	Env.BasePath = filepath.Dir(exePath)
	basePath := Env.BasePath + "/data/user.yaml"
	log.Printf(basePath)
	Env.AppName = filepath.Base(exePath)

	// step2: Read Config
	b, err := os.ReadFile(Env.BasePath + "/data/user.yaml")
	if err == nil {
		yaml.Unmarshal(b, &Config)
	}

	ytdlp.InitYtDlpConfig(Env.BasePath)
	MainApp = application.New(application.Options{
		Name:        "galaxy3",
		Description: "galaxy3",
		Icon:        Icon,
		Services: []application.Service{
			application.NewService(NewApp()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
}

func (a *App) RestartApp() FlagResult {
	exePath := Env.BasePath + "/" + Env.AppName

	cmd := exec.Command(exePath)
	HideExecWindow(cmd)

	err := cmd.Start()
	if err != nil {
		return FlagResult{false, err.Error()}
	}

	a.ExitApp()

	return FlagResult{true, "Success"}
}
