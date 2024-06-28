package spider

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var client *http.Client

func start(u string) {
	fmt.Println("Target URL: " + u)
	Wg.Add(1)
	Ch <- 1
	go Spider(u, 1)
	Wg.Wait()
	Progress = 1
	//fmt.Printf("\r\nSpider OK \n")
	ResultUrl = RemoveRepeatElement(ResultUrl)
	AddSource()
}

func AppendJs(ur string, urltjs string) int {
	Lock.Lock()
	defer Lock.Unlock()

	_, err := url.Parse(ur)
	if err != nil {
		return 2
	}
	for _, eachItem := range ResultJs {
		if eachItem.Url == ur {
			return 0
		}
	}
	ResultJs = append(ResultJs, Link{Url: ur})
	if strings.HasSuffix(urltjs, ".js") {
		Jsinurl[ur] = Jsinurl[urltjs]
	} else {
		re := regexp.MustCompile("[a-zA-z]+://[^\\s]*/|[a-zA-z]+://[^\\s]*")
		u := re.FindAllStringSubmatch(urltjs, -1)
		Jsinurl[ur] = u[0][0]
	}
	Jstourl[ur] = urltjs
	return 0

}

func AppendUrl(ur string, req *http.Request) int {

	Lock.Lock()
	defer Lock.Unlock()

	_, err := url.Parse(ur)
	if err != nil {
		return 2
	}
	for _, eachItem := range ResultUrl {
		if eachItem.Url == ur {
			return 0
		}
	}
	url.Parse(ur)
	ResultUrl = append(ResultUrl, Link{Url: ur})
	Urltourl[ur] = req.URL.String()
	return 0
}

func AppendEndUrl(url string) {
	Lock.Lock()
	defer Lock.Unlock()
	for _, eachItem := range EndUrl {
		if eachItem == url {
			return
		}
	}
	EndUrl = append(EndUrl, url)

}

func GetEndUrl(url string) bool {
	Lock.Lock()
	defer Lock.Unlock()
	for _, eachItem := range EndUrl {
		if eachItem == url {
			return true
		}
	}
	return false

}

func AddRedirect(url string) {
	Lock.Lock()
	defer Lock.Unlock()
	Redirect[url] = true
}

func AddSource() {
	for i := range ResultUrl {
		ResultUrl[i].Source = Urltourl[ResultUrl[i].Url]
	}

}

func Initialization() {
	ResultUrl = []Link{}

	EndUrl = []string{}
	Jsinurl = make(map[string]string)
	Jstourl = make(map[string]string)
	Urltourl = make(map[string]string)
	Redirect = make(map[string]bool)
	//SmartFilterrr = NewSmartFilter(NewSimpleFilter())
}
