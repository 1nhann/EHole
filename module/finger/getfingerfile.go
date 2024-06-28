package finger

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

type Packjson struct {
	Fingerprint []Fingerprint `json:"fingerprint"`
}

type Fingerprint struct {
	Cms      string   `json:"cms"`
	Method   string   `json:"method"`
	Location string   `json:"location"`
	Keyword  []string `json:"keyword"`
}

var (
	Webfingerprint *Packjson
)

func removeComments(data []byte) []byte {
	re := regexp.MustCompile(`(?m)^\s*//.*$`)
	return re.ReplaceAll(data, nil)
}
func LoadWebfingerprint(path string) error {
	data, err := ioutil.ReadFile(path)
	data = removeComments(data)
	if err != nil {
		return err
	}

	var config Packjson
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	Webfingerprint = &config
	return nil
}

func GetWebfingerprint() *Packjson {
	return Webfingerprint
}
