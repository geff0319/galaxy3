package bridge

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed image/Appicon.png
var Appicon []byte

//go:embed image/exit.png
var exit []byte

func InitSystray() {
	systemTray := MainApp.NewSystemTray()
	myMenu := MainApp.NewMenu()
	//设置icon会报错
	systemTray.SetIcon(Appicon)
	myMenu.Add("Exit").SetBitmap(exit).OnClick(func(ctx *application.Context) {
		MainWin.Close()
	})

	systemTray.SetMenu(myMenu)
	systemTray.OnClick(func() {
		if !MainWin.IsVisible() {
			MainWin.Show()
		}
		if !MainWin.IsFocused() {
			MainWin.Focus()
		}
	})
}
