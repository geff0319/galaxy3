package bridge

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
	"math"
	"time"
)

func WindowMask(width, height int, transparency uint8) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充透明背景
	clearColor := color.RGBA{A: transparency} // 完全透明
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
func BallMask(width, height int, transparency uint8) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{A: transparency})
		}
	}

	// 设置圆形区域的透明度为 10
	circleRadius := 25          // 半径为 25，圆形区域大小为 50x50
	circleCenterX := width - 30 // 圆心 x 坐标，位于右下角
	circleCenterY := 30         // 圆心 y 坐标，位于右下角

	// 在右下角绘制透明度为 10 的圆形
	for y := 0; y < 55; y++ {
		for x := width - 55; x < width; x++ {
			// 计算当前点到圆心的距离
			distance := math.Sqrt(float64((x-circleCenterX)*(x-circleCenterX) + (y-circleCenterY)*(y-circleCenterY)))
			if distance <= float64(circleRadius) {
				// 如果当前点在圆形内部，设置透明度为 10
				img.Set(x, y, color.RGBA{A: 1})
			}
		}
	}
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
		InitialPosition: application.WindowCentered,
		Hidden:          true,
		//ShouldClose: func(window *application.WebviewWindow) bool {
		//	//直接窗口不关闭，等处理完成app退出
		//	return false
		//},
		Windows: application.WindowsWindow{
			//BackdropType: application.Acrylic,
		},
		BackgroundColour: application.NewRGB(33, 37, 41),
		URL:              "/",
	})
	MainWin.SetRelativePosition(0, 0)
	MainWin.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		fmt.Println("quit app ")
		if Cd2Client != nil {
			Cd2Client.Close()
		}
		MainWin.EmitEvent("beforeClose")
	})
	MainWin.RegisterHook(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
		MainWin.EmitEvent("appInit")
		time.Sleep(200 * time.Millisecond)
		MainWin.Show()

		InitSystray()
		//InitCd2Client()
		InitScheduledTasks()
		CreateHook()

		InitBallWin()
		InitWidgetsWin()
	})
	//YtdlpWinShow()
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
			WindowMask:          WindowMask(380, 70, 13),
			WindowMaskDraggable: true,
			ExStyle:             w32.WS_EX_TOOLWINDOW,
		},
		BackgroundType: application.BackgroundTypeTransparent,
		URL:            "#/widgets/clock",
	})
	contextMenu := MainApp.NewMenu()
	contextMenu.Add("隐藏").OnClick(func(data *application.Context) {
		WidgetsWin.Hide()
	})
	WidgetsWin.RegisterContextMenu("WidgetsWinMenu", contextMenu)
}
func InitBallWin() {
	SingleBallWin := MainApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "SingleBallWin",
		Name:            "SingleBallWin",
		Frameless:       true,
		DisableResize:   true,
		Hidden:          true,
		AlwaysOnTop:     true,
		DevToolsEnabled: true,
		Windows: application.WindowsWindow{
			WindowMask:          BallMask(250, 180, 0),
			WindowMaskDraggable: true,
			ExStyle:             w32.WS_EX_TOOLWINDOW,
		},
		BackgroundType: application.BackgroundTypeTransparent,
		URL:            "#/widgets/single-ball",
	})
	contextMenu := MainApp.NewMenu()
	contextMenu.Add("隐藏").OnClick(func(data *application.Context) {
		SingleBallWin.Hide()
	})

	BallMenuWin := MainApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "BallMenuWin",
		Name:            "BallMenuWin",
		Frameless:       true,
		DisableResize:   true,
		Hidden:          true,
		AlwaysOnTop:     true,
		DevToolsEnabled: true,
		Windows: application.WindowsWindow{
			WindowMask:          BallMask(250, 180, 1),
			WindowMaskDraggable: true,
			ExStyle:             w32.WS_EX_TOOLWINDOW,
		},
		BackgroundType: application.BackgroundTypeTransparent,
		URL:            "#/widgets/ball-menu",
	})
	contextMenu.Add("隐藏").OnClick(func(data *application.Context) {
		BallMenuWin.Hide()
	})

}
func YtdlpWinShow() {
	id := generateRandomString(8)
	win := MainApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:           "解析链接",
		Name:            id,
		Width:           480,
		Height:          220,
		Frameless:       true,
		DisableResize:   true,
		Hidden:          true,
		AlwaysOnTop:     true,
		InitialPosition: application.WindowCentered,
		DevToolsEnabled: true,
		Windows:         application.WindowsWindow{
			//BackdropType:        application.Acrylic,
			//WindowMask:          WindowMask(600, 230),
			//WindowMaskDraggable: true,
			//ExStyle:             w32.WS_EX_TOOLWINDOW,
		},
		//BackgroundType: application.BackgroundTypeTranslucent,
		BackgroundColour: application.NewRGB(33, 37, 41),
		URL:              "#/ytdlpWidgets",
	})

	win.RegisterHook(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
		OptionShow(id)
		time.Sleep(200 * time.Millisecond)
		win.Show()
	})
}

func OptionShow(fatherId string) {
	fmt.Println(fatherId + "option")
	win := MainApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:         "选项",
		Name:          fatherId + "option",
		Width:         150,
		Height:        170,
		Frameless:     true,
		DisableResize: true,
		Hidden:        true,
		AlwaysOnTop:   true,
		//Centered:        true,
		DevToolsEnabled: true,
		Windows: application.WindowsWindow{
			//WindowMask: WindowMask(150, 170),
			ExStyle: w32.WS_EX_TOOLWINDOW,
		},
		BackgroundColour: application.NewRGB(33, 37, 41),
		URL:              "#/option",
	})
	win.RegisterHook(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
		win.Hide()
	})
}
