package spider

import (
	_ "embed"
)

var (
	ResultJs  []Link
	ResultUrl []Link

	EndUrl   []string
	Jsinurl  map[string]string
	Jstourl  map[string]string
	Urltourl map[string]string
	Redirect map[string]bool
)
