package bridge

import "C"
import (
	"bytes"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"image"
	"image/color"
	"image/png"
)

func WindowMask(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充透明背景
	clearColor := color.RGBA{255, 255, 255, 13} // 完全透明
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, clearColor)
		}
	}

	// 将图像编码为 PNG 格式
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		gefflog.Err("WindowMask err: " + err.Error())
	}

	// 获取 PNG 数据的 []byte
	return buf.Bytes()
}
func InitMianWin() {
	MainWin = MainApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "galaxy3",
		Name:            "MainWin",
		Width:           1000,
		Height:          680,
		Frameless:       true,
		DevToolsEnabled: true,
		DisableResize:   true,
		Centered:        true,
		//ShouldClose: func(window *application.WebviewWindow) bool {
		//	b, err := os.ReadFile(Env.BasePath + "/data/user.yaml")
		//	if err != nil {
		//		MainApp.Quit()
		//		return true
		//	}
		//	err = yaml.Unmarshal(b, &Config)
		//	if err != nil {
		//		MainApp.Quit()
		//		return true
		//	}
		//	if Config.ExitOnClose {
		//		//WidgetsWin.Close()
		//		MainApp.Quit()
		//		return true
		//	} else {
		//		MainWin.Hide()
		//		return false
		//	}
		//},
		//KeyBindings: map[string]func(window *application.WebviewWindow){
		//	"alt+w": func(window *application.WebviewWindow) {
		//		YtdlpWinShow()
		//	},
		//},
		//BackgroundType: application.BackgroundTypeTranslucent,
		Windows: application.WindowsWindow{
			//BackdropType: application.Acrylic,
		},
		BackgroundColour: application.NewRGB(33, 37, 41),
		URL:              "/",
	})
	MainWin.SetRelativePosition(0, 0)
	MainWin.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		fmt.Println("quit app ")
		MainWin.EmitEvent("beforeClose")
	})
}

func InitWidgetsWin() {
	WidgetsWin = MainApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "WidgetsWin",
		Name:            "WidgetsWin",
		Frameless:       true,
		DisableResize:   true,
		Hidden:          true,
		AlwaysOnTop:     true,
		DevToolsEnabled: true,
		Windows: application.WindowsWindow{
			WindowMask:          WindowMask(380, 70),
			WindowMaskDraggable: true,
			ExStyle:             w32.WS_EX_TOOLWINDOW,
		},
		BackgroundType: application.BackgroundTypeTransparent,
		URL:            "#/widgets",
	})
	contextMenu := MainApp.NewMenu()
	contextMenu.Add("隐藏").OnClick(func(data *application.Context) {
		WidgetsWin.Hide()
	})
	WidgetsWin.RegisterContextMenu("WidgetsWinMenu", contextMenu)
}

func YtdlpWinShow() {
	id := generateRandomString(8)
	win := MainApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "解析链接",
		Name:            id,
		Frameless:       true,
		DisableResize:   true,
		Hidden:          true,
		AlwaysOnTop:     true,
		Centered:        true,
		DevToolsEnabled: true,
		Windows: application.WindowsWindow{
			//BackdropType:        application.Acrylic,
			WindowMask:          WindowMask(600, 230),
			WindowMaskDraggable: true,
			//ExStyle:             w32.WS_EX_TOOLWINDOW,
		},
		//BackgroundType: application.BackgroundTypeTranslucent,
		BackgroundColour: application.NewRGB(33, 37, 41),
		URL:              "#/ytdlpWidgets",
	})
	win.RegisterHook(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
		win.Show()
		fmt.Println(id)
		win.EmitEvent("winId", id)
	})

}
