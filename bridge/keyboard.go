package bridge

import (
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	hook "github.com/robotn/gohook"
)

func (a *App) ExitKey() {
	//hook.End()
}

func CreateHook() {
	go func() {
		hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
			fmt.Println("ctrl-shift-q")
			hook.End()
		})
		hook.Register(hook.KeyDown, []string{"alt", "w"}, func(e hook.Event) {
			text, ok := MainApp.Clipboard().Text()
			if !ok {
				gefflog.Err("Failed to get clipboard text")
				return
			}
			if len(text) != 0 {
				YtdlpWinShow()
			}
		})

		s := hook.Start()
		<-hook.Process(s)
	}()
}
