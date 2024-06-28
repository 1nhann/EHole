package spider

import (
	"regexp"
)

//var SmartFilterrr *SmartFilter
//var Wgspider sizedwaitgroup.SizedWaitGroup

func Spide(url string) []Link {
	Initialization()
	host := GetHost(url)

	start(url)
	ResultUrlHost, _ := UrlDispose(ResultUrl, host, GetHost(url))
	return ResultUrlHost
}

func GetHostFromUrl(url string) string {
	//host:port
	var host string
	re := regexp.MustCompile("([a-z0-9\\-]+\\.)*([a-z0-9\\-]+\\.[a-z0-9\\-]+)(:[0-9]+)?")
	hosts := re.FindAllString(url, 1)
	if len(hosts) == 0 {
		host = url
	} else {
		host = hosts[0]
	}
	return host
}
