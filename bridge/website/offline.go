package website

import (
	"regexp"
)

type Offline struct {
	sourceUrl string
}

func NewOffline(url string) *Offline {
	return &Offline{
		sourceUrl: url,
	}
}
func (t *Offline) Compile() (string, bool) {
	return "", true
}

func (t *Offline) AppCompile() (string, bool) {
	// 匹配磁力链接的正则表达式（以 magnet:? 开头）
	magnetURLPattern := `^magnet:\?.*`

	// 编译正则表达式
	magnetRe := regexp.MustCompile(magnetURLPattern)
	magnetLinks := magnetRe.FindAllString(t.sourceUrl, -1)
	if len(magnetLinks) != 0 {
		return magnetLinks[0], true // 返回第一个磁力链接
	}
	return "", false
}
