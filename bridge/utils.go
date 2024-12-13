package bridge

import (
	"github.com/ge-fei-fan/gefflog"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func GetPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	path = filepath.Join(Env.BasePath, path)
	path = filepath.Clean(path)
	return path
}

func GetProxy(_proxy string) func(*http.Request) (*url.URL, error) {
	proxy := http.ProxyFromEnvironment

	if _proxy != "" {
		if !strings.HasPrefix(_proxy, "http") {
			_proxy = "http://" + _proxy
		}
		proxyUrl, err := url.Parse(_proxy)
		if err == nil {
			proxy = http.ProxyURL(proxyUrl)
		}
	}

	return proxy
}

func ConvertByte2String(byte []byte) string {
	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
	return string(decodeBytes)
}

func (a *App) AbsolutePath(path string) FlagResult {
	log.Printf("AbsolutePath: %s", path)

	if filepath.IsAbs(path) {
		return FlagResult{true, path}
	}

	path = GetPath(path)

	return FlagResult{true, path}
}

func (a *App) HideToolWindow() {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(AppTitle))
	// 获取当前窗口的扩展风格
	currentExStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	// 添加 WS_EX_TOOLWINDOW 样式
	newExStyle := (currentExStyle | win.WS_EX_TOOLWINDOW) &^ win.WS_EX_APPWINDOW
	//win.SetWindowLong(hwnd, win.GWL_EXSTYLE, win.GetWindowLong(hwnd, win.GWL_EXSTYLE)|win.WS_EX_LAYERED)
	// 应用新的扩展风格
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, newExStyle)
}
func (a *App) ShowToolWindow() {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(AppTitle))
	// 获取当前窗口的扩展风格
	currentExStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	// 添加 WS_EX_TOOLWINDOW 样式
	newExStyle := currentExStyle &^ win.WS_EX_TOOLWINDOW
	// 应用新的扩展风格
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, newExStyle)
}

// Function to generate a random string of a specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

// 返回option窗口展示位置，展示在鼠标正下方
func (a *App) GetBelowWinPos(winWidth, winHeight int) FlagResultWithData {
	var dstX, dstY int
	// 获取屏幕的宽度和高度
	screenWidth, screenHeight := robotgo.GetScreenSize()
	// 获取当前鼠标的位置
	x, y := robotgo.Location()

	distanceBottom := screenHeight - y
	// 鼠标下方展示不下

	if (winHeight + 10) > distanceBottom {
		dstY = y - winHeight - 10
	} else {
		dstY = y + 10
	}

	distanceRight := screenWidth - x
	halfWinWidth := winWidth / 2
	//鼠标右边展示不够

	if halfWinWidth > distanceRight {
		dstX = x - winWidth + distanceRight
	} else if halfWinWidth > x {
		dstX = 0
	} else {
		dstX = x - halfWinWidth
	}

	return FlagResultWithData{true, "", map[string]int{"dstX": dstX, "dstY": dstY}}
}
func (a *App) ChangeLog(level byte, path string) {
	Config.LogPath = path
	gefflog.ChangeLogger(level, path)
	gefflog.Info("修改应用日志：" + path)
}
func IsYtDlpExist() bool {
	_, err := os.Stat(YdpConfig.YtDlpPath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
