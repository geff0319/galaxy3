package bridge

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed Appicon.png
var Appicon []byte

func InitSystray() {
	systemTray := MainApp.NewSystemTray()
	myMenu := MainApp.NewMenu()
	//设置icon会报错
	systemTray.SetIcon(Appicon)
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
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
