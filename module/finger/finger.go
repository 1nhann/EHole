package finger

import (
	"ehole/module/queue"
	"ehole/module/spider"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"os"
	"sort"
	"strings"
	"sync"
)

type Outrestul struct {
	Url        string `json:"url"`
	Cms        string `json:"cms"`
	Server     string `json:"server"`
	Statuscode int    `json:"statuscode"`
	Length     int    `json:"length"`
	Title      string `json:"title"`
}

type FinScan struct {
	UrlQueue    *queue.Queue
	Ch          chan []string
	Wg          sync.WaitGroup
	Thread      int
	Output      string
	Proxy       string
	AllResult   []Outrestul
	FocusResult []Outrestul
	Finpx       *Packjson
}

func NewScan(urls []string, thread int, output string, proxy string, fingerJson string, sp bool) *FinScan {
	//var wgspider sizedwaitgroup.SizedWaitGroup
	//wgspider = sizedwaitgroup.New(5)
	s := &FinScan{
		UrlQueue:    queue.NewQueue(),
		Ch:          make(chan []string, thread),
		Wg:          sync.WaitGroup{},
		Thread:      thread,
		Output:      output,
		Proxy:       proxy,
		AllResult:   []Outrestul{},
		FocusResult: []Outrestul{},
	}
	err := LoadWebfingerprint(fingerJson)
	if err != nil {
		color.RGBStyleFromString("237,64,35").Println("[error] fingerprint file error!!! parse json error !!")
		os.Exit(1)
	}
	s.Finpx = GetWebfingerprint()
	for _, url := range urls {
		s.UrlQueue.Push([]string{url, "0"})
	}
	//wgspider.Wait()
	return s
}

func (s *FinScan) StartScan(sp bool) {
	//var wgspider sizedwaitgroup.SizedWaitGroup
	//wgspider = sizedwaitgroup.New(5)
	//if sp {
	//	for _, url := range urls {
	//		wgspider.Add()
	//		go func(uuu string) {
	//			defer wgspider.Done()
	//			spidedUrls := spider.Spide(uuu)
	//			for _, spidedUrl := range spidedUrls {
	//				if spidedUrl.Url == uuu {
	//					continue
	//				}
	//				s.UrlQueue.Push([]string{spidedUrl.Url, "3"})
	//			}
	//
	//		}(url)
	//	}
	//}
	//wgspider.Wait()
	for i := 0; i <= s.Thread; i++ {
		s.Wg.Add(1)
		go func() {
			defer s.Wg.Done()
			s.fingerScan(sp)
		}()
	}
	s.Wg.Wait()

	color.RGBStyleFromString("244,211,49").Println("\n重点资产：")

	var AllResult sync.Map    //前缀
	var AllAllResult sync.Map //全部
	for _, aas := range s.FocusResult {
		AllAllResult.Store(aas.Url, aas)
		var urlsWithSameFinger []string
		AllAllResult.Range(func(key, value interface{}) bool {
			uu := key.(string)
			ff := value.(Outrestul).Cms
			t := value.(Outrestul).Title
			if aas.Cms == ff && spider.GetHostFromUrl(uu) == spider.GetHostFromUrl(aas.Url) && aas.Title == t {
				urlsWithSameFinger = append(urlsWithSameFinger, uu)
			}
			return true
		})
		// 计算最长公共前缀
		commonPrefix := findCommonPrefix(urlsWithSameFinger)
		AllResult.Store(commonPrefix, Outrestul{
			Url:        commonPrefix,
			Cms:        aas.Cms,
			Server:     aas.Server,
			Statuscode: aas.Statuscode,
			Title:      aas.Title,
		})
		AllResult.Range(func(key, value interface{}) bool {
			uu := key.(string)
			ff := value.(Outrestul).Cms
			t := value.(Outrestul).Title
			if uu != commonPrefix && ff == aas.Cms && spider.GetHostFromUrl(uu) == spider.GetHostFromUrl(aas.Url) && aas.Title == t {
				AllResult.Delete(uu)
			}
			return true
		})

	}
	var keys []string
	keyValueMap := make(map[string]Outrestul)

	AllResult.Range(func(key, value interface{}) bool {
		u := key.(string)
		aas := value.(Outrestul)
		keyValueMap[u] = aas
		keys = append(keys, u)
		return true
	})
	sort.Strings(keys)
	for _, u := range keys {
		aas := keyValueMap[u]
		fmt.Printf(fmt.Sprintf("[ %s | ", u))
		color.RGBStyleFromString("237,64,35").Printf(fmt.Sprintf("%s", aas.Cms))
		fmt.Printf(fmt.Sprintf(" | %s | %d | %d | %s ]\n", aas.Server, aas.Statuscode, aas.Length, aas.Title))
	}
	if s.Output != "" {
		outfile(s.Output, s.AllResult)
	}
}

func MapToJson(param map[string][]string) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func (s *FinScan) fingerScan(sp bool) {
	for s.UrlQueue.Len() != 0 {
		dataface := s.UrlQueue.Pop()
		switch dataface.(type) {
		case []string:
			url := dataface.([]string)
			var urlType string
			if url[1] == "3" {
				urlType = "url from spider"
			} else if url[1] == "0" && sp {
				spidedUrls := spider.Spide(url[0])
				for _, spidedUrl := range spidedUrls {
					if spidedUrl.Url == url[0] {
						continue
					}
					s.UrlQueue.Push([]string{spidedUrl.Url, "3"})
				}
			}

			var data *resps
			data, err := httprequest(url, s.Proxy)
			if err != nil {
				url[0] = strings.ReplaceAll(url[0], "https://", "http://")
				data, err = httprequest(url, s.Proxy)
				if err != nil {
					continue
				}
			}
			for _, jurl := range data.jsurl {
				if jurl != "" {
					s.UrlQueue.Push([]string{jurl, "1"})
				}
			}
			headers := MapToJson(data.header)
			var cms []string
			for _, finp := range s.Finpx.Fingerprint {
				if finp.Location == "body" {
					if finp.Method == "keyword" {
						if iskeyword(data.body, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "faviconhash" {
						if data.favhash == finp.Keyword[0] {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if isregular(data.body, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
				if finp.Location == "header" {
					if finp.Method == "keyword" {
						if iskeyword(headers, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if isregular(headers, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
				if finp.Location == "title" {
					if finp.Method == "keyword" {
						if iskeyword(data.title, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if isregular(data.title, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
			}
			cms = RemoveDuplicatesAndEmpty(cms)
			sort.Strings(cms)
			cmss := strings.Join(cms, ",")
			if urlType == "url from spider" {
				if cmss == "" || data.statuscode/100 == 4 || data.statuscode/100 == 5 {
					continue
				}
			}

			out := Outrestul{data.url, cmss, data.server, data.statuscode, data.length, data.title}
			s.AllResult = append(s.AllResult, out)
			if len(out.Cms) != 0 || out.Title != "" || out.Server != "" {
				outstr := fmt.Sprintf("[ %s | %s | %s | %d | %d | %s ]", out.Url, out.Cms, out.Server, out.Statuscode, out.Length, out.Title)
				color.RGBStyleFromString("237,64,35").Println(outstr)

				s.FocusResult = append(s.FocusResult, out)
			} else {
				outstr := fmt.Sprintf("[ %s | %s | %s | %d | %d | %s ]", out.Url, out.Cms, out.Server, out.Statuscode, out.Length, out.Title)
				fmt.Println(outstr)
			}
		default:
			continue
		}
	}
}
func findCommonPrefix(urls []string) string {
	if len(urls) == 0 {
		return ""
	}

	// 找到最短的 URL 长度
	minLen := len(urls[0])
	for _, url := range urls {
		if len(url) < minLen {
			minLen = len(url)
		}
	}

	// 找到最长公共前缀
	var commonPrefix strings.Builder
	for i := 0; i < minLen; i++ {
		char := urls[0][i]
		for _, url := range urls {
			if url[i] != char {
				return commonPrefix.String()
			}
		}
		commonPrefix.WriteByte(char)
	}
	return commonPrefix.String()
}
