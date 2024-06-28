package spider

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// 蜘蛛抓取页面内容
func Spider(u string, num int) {
	is := true
	defer func() {
		Wg.Done()
		if is {
			<-Ch
		}

	}()
	//Mux.Lock()
	//fmt.Printf("\rStart %d Spider...", Progress)
	//Progress++
	//
	//Mux.Unlock()
	//标记完成

	u, _ = url.QueryUnescape(u)

	if GetEndUrl(u) {
		return
	}
	for _, v := range Risks {
		if strings.Contains(u, v) {
			return
		}
	}
	AppendEndUrl(u)
	request, err := http.NewRequest("GET", u, bytes.NewReader([]byte{}))
	if err != nil {
		return
	}

	request.Header.Set("Accept-Encoding", "gzip") //使用gzip压缩传输数据让访问更快
	request.Header.Set("User-Agent", GetUserAgent())
	request.Header.Set("Accept", "*/*")

	cookie := &http.Cookie{
		Name:  "rememberMe",
		Value: "me",
	}
	request.AddCookie(cookie)

	//加载yaml配置(headers)
	client = &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	result := ""
	//解压
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(response.Body) // gzip解压缩
		if err != nil {
			return
		}
		defer reader.Close()
		con, err := io.ReadAll(reader)
		if err != nil {
			return
		}
		result = string(con)
	} else {
		//提取url用于拼接其他url或js
		dataBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return
		}
		//字节数组 转换成 字符串
		result = string(dataBytes)
	}
	path := response.Request.URL.Path
	host := response.Request.URL.Host
	scheme := response.Request.URL.Scheme
	//处理base标签
	re := regexp.MustCompile("base.{1,5}href.{1,5}(http.+?//[^\\s]+?)[\"'‘“]")
	base := re.FindAllStringSubmatch(result, -1)
	if len(base) > 0 {
		host = regexp.MustCompile("http.*?//([^/]+)").FindAllStringSubmatch(base[0][1], -1)[0][1]
		scheme = regexp.MustCompile("(http.*?)://").FindAllStringSubmatch(base[0][1], -1)[0][1]
		paths := regexp.MustCompile("http.*?//.*?(/.*)").FindAllStringSubmatch(base[0][1], -1)
		if len(paths) > 0 {
			path = paths[0][1]
		} else {
			path = "/"
		}
	}
	is = false
	<-Ch
	//提取js
	jsFind(result, host, scheme, path, u, num)
	//提取url

	urlFind(result, request, num)
	//提取信息

}
