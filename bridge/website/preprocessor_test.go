package website

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestBi(t *testing.T) {

	b := Bilibili{sourceUrl: "【小小牛仔裤，拿下！-哔哩哔哩】 https://b23.tv/FJbWCkw"}
	s, ok := b.AppCompile()
	fmt.Println(ok)
	fmt.Println(s)
	//assert.Equal(t, "https://b23.tv/FJbWCkw", s)
}

func TestBV(t *testing.T) {
	rawURL := "https://www.bilibili.com/video/BV11z421a7eP?-Arouter=story&buvid=Z649F066A39DF7174D09877BE99E4EE21BE6&from_spmid=main.ugc-video-detail-vertical.0.0&is_story_h5=true&mid=SNGSACb7aZGM61g%2FweOIAH8FTQ%2FSZMtL1rElX6M3iMo%3D&p=1&plat_id=143&share_from=ugc&share_medium=iphone&share_plat=ios&share_session_id=4B041299-1085-4F5E-BB01-55B9F73C3AEC&share_source=COPY&share_tag=s_i&spmid=main.ugc-video-detail-verticalspace.0.0&timestamp=1724821911&unique_k=4Ouf3Pw&up_id=484810608"
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	// Extract the path and split by '/'
	segments := strings.Split(u.Path, "/")
	if len(segments) > 2 {
		fmt.Println("BV Code:", segments[2])
	} else {
		fmt.Println("BV Code not found")
	}
}

func TestD(t *testing.T) {
	//var client = &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyFromEnvironment,
	//		//禁止复用连接，防止同一个连接长时间大流量被限速
	//		DisableKeepAlives: true,
	//	},
	//}
	//u := "https://upos-hz-mirrorakam.akamaized.net/upgcxcode/72/49/25709184972/25709184972-1-30120.m4s?e=ig8euxZM2rNcNbdlhoNvNC8BqJIzNbfqXBvEqxTEto8BTrNvN0GvT90W5JZMkX_YN0MvXg8gNEV4NC8xNEV4N03eN0B5tZlqNxTEto8BTrNvNeZVuJ10Kj_g2UB02J0mN0B5tZlqNCNEto8BTrNvNC7MTX502C8f2jmMQJ6mqF2fka1mqx6gqj0eN0B599M=&uipk=5&nbs=1&deadline=1725887681&gen=playurlv2&os=akam&oi=2491906636&trid=79cb1ed371ca4c2aad5d939c43892048u&mid=34801693&platform=pc&og=cos&upsig=9406cf4dc45df3a702ea721582a353e5&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,mid,platform,og&hdnts=exp=1725887681~hmac=0b84d12a7515e1e0a6de57327c0fe09398c603d9362207add558788f77f917bf&bvc=vod&nettype=0&orderid=0,1&buvid=F90DF026-1264-3D20-794D-F215C9F274AF98340infoc&build=0&f=u_0_0&agrr=0&bw=1191639&logo=80000000"
	//// 发送 GET 请求以实际下载文件
	//req, _ := http.NewRequest(http.MethodGet, u, nil)
	//req.Header.Set("Referer", "https://www.bilibili.com")
	//
	//getResp, err := client.Do(req)
	//if err != nil {
	//	t.Fatalf("Failed to send GET request: %v", err)
	//}
	//defer getResp.Body.Close()
	//
	//// 验证内容大小
	//fmt.Println(getResp.ContentLength)
	//metadata, err := GetBilibiliInfo("https://www.bilibili.com/video/BV11z421a7eP", "")
	//if err != nil {
	//	t.Log(err)
	//}
	//fmt.Printf("metadata:%v", metadata)
	name := "Love one 유솜 Yousom 러브원 Drama Queen 축하공연 팬사인회 아세아항공직업전문학교 [ZfVEbHx_4K0].webm]"
	name1 := strings.ReplaceAll(name, " ", "")
	fmt.Println(name1)
}

func TestThumbnail(t *testing.T) {
	//resp, err := http.Get("http://i1.hdslb.com/bfs/archive/117b3665d363a4bcb58c49d1bd429089fc1094ea.jpg")
	//defer resp.Body.Close()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//path := "C:\\Users\\geff\\Desktop\\code\\Galaxy\\build\\bin\\data\\yt-dlp-download\\Thumbnail"
	//
	//if !DirExists(path) {
	//	err := os.MkdirAll(path, os.ModePerm)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//}

	//file, err := os.OpenFile(filepath.Join(path, "111.jpg"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	//defer file.Close()
	//_, err = io.Copy(file, resp.Body)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//url := "https://api.bilibili.com/x/space/wbi/acc/info?mid=34801693"
	urlstr := "https://api.bilibili.com/x/web-interface/nav"
	urlstr, err := sign(urlstr)
	fmt.Println(urlstr)
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.AddCookie(&http.Cookie{Name: "SESSDATA", Value: "11098c71%2C1742182391%2Cd6cfe%2A92CjDjkGIBU1BONTU4L7Rqaqh29uUON0p7_AmXLUY879U7kO7NM9ztoSnDVrBocwThv-kSVkpIRHFTS28zUkMzdFZwYnE2S0x0a2t4V0ZrM3J2bU9IOHpicGg1RElYSUE5OEF2RmVWVzU2ZXNjRl9HR3JtU1Fhc0VpTFVfMkV1YXV0QjFtMWxJalFRIIEC"})
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %s", err)
		return
	}
	defer response.Body.Close()
	res, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Request failed: %s", err)
		return
	}
	fmt.Println(string(res))
}
func DirExists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if f.IsDir() {
		return true
	}
	return false
}
