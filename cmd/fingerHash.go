package cmd

import (
	"crypto/tls"
	"ehole/module/finger"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var fingerHash = &cobra.Command{
	Use: "hash",
	Run: func(cmd *cobra.Command, args []string) {
		if u != "" {
			f, _ := fav([]string{u}, proxy)
			if f == "" {

			} else {
				fmt.Println(f)
			}
		} else {
			cmd.Usage()
		}
	},
}
var (
	u string
)

func init() {
	fingerCmd.AddCommand(fingerHash)
	fingerHash.Flags().StringVarP(&u, "url", "u", "", "url")
}

func fav(url1 []string, proxy string) (string, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		proxys := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           proxys,
		}
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	req, err := http.NewRequest("GET", url1[0], nil)
	if err != nil {
		return "", err
	}
	//cookie := &http.Cookie{
	//	Name:  "rememberMe",
	//	Value: "me",
	//}
	//req.AddCookie(cookie)
	req.Header.Set("Accept", "*/*;q=0.8")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", finger.Rndua())
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	httpbody := string(result)
	httpbody = finger.ToUtf8(httpbody, contentType)

	favhash := finger.Getfavicon(httpbody, url1[0])
	return favhash, nil
}

//func xegexpjs(reg string, resp string) (reslut1 [][]string) {
//	reg1 := regexp.MustCompile(reg)
//	if reg1 == nil {
//		log.Println("regexp err")
//		return nil
//	}
//	result1 := reg1.FindAllStringSubmatch(resp, -1)
//	return result1
//}

//func getfavicon(httpbody string, turl string) string {
//	faviconpaths := xegexpjs(`href="(.*?favicon....)"`, httpbody)
//	var faviconpath string
//	u, err := url.Parse(turl)
//	if err != nil {
//		panic(err)
//	}
//	turl = u.Scheme + "://" + u.Host
//	if len(faviconpaths) > 0 {
//		fav := faviconpaths[0][1]
//		if fav[:2] == "//" {
//			faviconpath = "http:" + fav
//		} else {
//			if fav[:4] == "http" {
//				faviconpath = fav
//			} else {
//				faviconpath = turl + "/" + fav
//			}
//
//		}
//	} else {
//		faviconpath = turl + "/favicon.ico"
//	}
//	return favicohash(faviconpath)
//}
//func favicohash(host string) string {
//	timeout := time.Duration(8 * time.Second)
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := http.Client{
//		Timeout:   timeout,
//		Transport: tr,
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return http.ErrUseLastResponse /* 不进入重定向 */
//		},
//	}
//	resp, err := client.Get(host)
//	if err != nil {
//		//log.Println("favicon client error:", err)
//		return "0"
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode == 200 {
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			//log.Println("favicon file read error: ", err)
//			return "0"
//		}
//		return finger.Mmh3Hash32(finger.StandBase64(body))
//	} else {
//		return "0"
//	}
//}

//func rndua() string {
//	ua := []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
//		"Mozilla/5.0 (X11; Linux x86_64; rv:96.0) Gecko/20100101 Firefox/96.0",
//		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",
//		"Mozilla/5.0 (X11; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0",
//		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36"}
//	n := rand.Intn(13) + 1
//	return ua[n]
//}
//func toUtf8(content string, contentType string) string {
//	var htmlEncode string
//	var htmlEncode2 string
//	var htmlEncode3 string
//	htmlEncode = "gb18030"
//	if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
//		htmlEncode = "gb18030"
//	} else if strings.Contains(contentType, "big5") {
//		htmlEncode = "big5"
//	} else if strings.Contains(contentType, "utf-8") {
//		//实际上，这里获取的编码未必是正确的，在下面还要做比对
//		htmlEncode = "utf-8"
//	}
//
//	reg := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`)
//	match := reg.FindStringSubmatch(content)
//	if len(match) > 1 {
//		contentType = strings.ToLower(match[1])
//		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
//			htmlEncode2 = "gb18030"
//		} else if strings.Contains(contentType, "big5") {
//			htmlEncode2 = "big5"
//		} else if strings.Contains(contentType, "utf-8") {
//			htmlEncode2 = "utf-8"
//		}
//	}
//
//	reg = regexp.MustCompile(`(?is)<title[^>]*>(.*?)<\/title>`)
//	match = reg.FindStringSubmatch(content)
//	if len(match) > 1 {
//		aa := match[1]
//		_, contentType, _ = charset.DetermineEncoding([]byte(aa), "")
//		contentType = strings.ToLower(contentType)
//		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
//			htmlEncode3 = "gb18030"
//		} else if strings.Contains(contentType, "big5") {
//			htmlEncode3 = "big5"
//		} else if strings.Contains(contentType, "utf-8") {
//			htmlEncode3 = "utf-8"
//		}
//	}
//
//	if htmlEncode != "" && htmlEncode2 != "" && htmlEncode != htmlEncode2 {
//		htmlEncode = htmlEncode2
//	}
//	if htmlEncode == "utf-8" && htmlEncode != htmlEncode3 {
//		htmlEncode = htmlEncode3
//	}
//
//	if htmlEncode != "" && htmlEncode != "utf-8" {
//		content = Convert(content, htmlEncode, "utf-8")
//	}
//
//	return content
//}
//func Convert(src string, srcCode string, tagCode string) string {
//	if srcCode == tagCode {
//		return src
//	}
//	srcCoder := mahonia.NewDecoder(srcCode)
//	srcResult := srcCoder.ConvertString(src)
//	tagCoder := mahonia.NewDecoder(tagCode)
//	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
//	result := string(cdata)
//	return result
//}
