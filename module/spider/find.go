package spider

import (
	"net/http"
	pathlib "path"
	"regexp"
	"strings"
)

// 分析内容中的js
func jsFind(cont, host, scheme, path, source string, num int) {
	var cata string
	care := regexp.MustCompile("/.*/{1}|/")
	catae := care.FindAllString(path, -1)
	if len(catae) == 0 {
		cata = "/"
	} else {
		cata = catae[0]
	}
	//js匹配正则
	host = scheme + "://" + host
	for _, re := range JsFind {
		reg := regexp.MustCompile(re)
		jss := reg.FindAllStringSubmatch(cont, -1)
		//return
		jss = jsFilter(jss)
		//循环提取js放到结果中
		for _, js := range jss {
			if js[0] == "" {
				continue
			}
			if strings.HasPrefix(js[0], "https:") || strings.HasPrefix(js[0], "http:") {
				switch AppendJs(js[0], source) {
				case 0:
					if num <= JsSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(js[0], num+1)
					}
				case 1:
					return
				case 2:
					continue
				}

			} else if strings.HasPrefix(js[0], "//") {
				switch AppendJs(scheme+":"+js[0], source) {
				case 0:
					if num <= JsSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(scheme+":"+js[0], num+1)
					}
				case 1:
					return
				case 2:
					continue
				}

			} else if strings.HasPrefix(js[0], "/") {
				switch AppendJs(host+js[0], source) {
				case 0:
					if num <= JsSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(host+js[0], num+1)
					}
				case 1:
					return
				case 2:
					continue
				}

			} else {
				switch AppendJs(host+cata+js[0], source) {
				case 0:
					if num <= JsSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(host+cata+js[0], num+1)
					}
				case 1:
					return
				case 2:
					continue
				}

			}
		}

	}

}

// 分析内容中的url
func urlFind(cont string, req *http.Request, num int) {
	scheme := req.URL.Scheme
	host := req.URL.Host
	source := req.URL.String()
	host = scheme + "://" + host
	path := req.URL.Path

	//url匹配正则

	for _, re := range UrlFind {
		reg := regexp.MustCompile(re)
		urls := reg.FindAllStringSubmatch(cont, -1)
		//fmt.Println(urls)
		urls = urlFilter(urls)
	label1:
		//循环提取url放到结果中
		for _, url := range urls {
			if url[0] == "" {
				continue
			}
			if strings.HasPrefix(url[0], "https:") || strings.HasPrefix(url[0], "http:") {
				switch AppendUrl(url[0], req) {
				case 0:
					if num <= UrlSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(url[0], num+1)
					}
				case 1:
					return
				case 2:
					continue
				}
			} else if strings.HasPrefix(url[0], "//") {
				switch AppendUrl(scheme+":"+url[0], req) {
				case 0:
					if num <= UrlSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(scheme+":"+url[0], num+1)
					}
				case 1:
					return
				case 2:
					continue
				}

			} else if strings.HasPrefix(url[0], "/") {
				urlz := ""
				urlz = host + url[0]
				switch AppendUrl(urlz, req) {
				case 0:
					if num <= UrlSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(urlz, num+1)
					}
				case 1:
					return
				case 2:
					continue
				}
			} else {
				//如果是js、css当中的url，不以/开头，则直接忽视
				if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
					continue
				}
				var ignore = []string{"text/javascript", "image/x-icon", "text/html"}
				for _, i := range ignore {
					if i == url[0] {
						continue label1
					}
				}

				dir := pathlib.Dir(source)
				if !strings.HasSuffix(dir, "/") {
					dir += "/"
				}
				urlz := host + dir + url[0]
				switch AppendUrl(urlz, req) {
				case 0:
					if num <= UrlSteps {
						Wg.Add(1)
						Ch <- 1
						go Spider(urlz, num+1)
					}
				case 1:
					return
				case 2:
					continue
				}
			}

		}
	}
}
