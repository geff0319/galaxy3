package ytdlp

import (
	"os"
	"regexp"
)

func IsFileExist(path string) bool {
	// 使用 os.Stat 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		//if os.IsNotExist(err) {
		//	// 文件不存在
		//	return false
		//}
		//// 其他错误
		return false
	}

	// 检查文件是否是一个目录
	if fileInfo.IsDir() {
		return false
	}

	// 文件存在且不是目录
	return true
}

func IsDirExists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if f.IsDir() {
		return true
	}
	return false
}

// 过滤 Windows 不支持的文件名字符
func SanitizeFileName(name string) string {
	// 定义一个正则表达式模式来匹配 Windows 不支持的字符
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	// 用下划线替换这些字符
	sanitized := re.ReplaceAllString(name, "_")
	return sanitized
}
