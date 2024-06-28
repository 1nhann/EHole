package cmd

import (
	"ehole/module/finger"
	"fmt"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

var fingerGetCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		if product != "" {
			if viper.GetString("finger.finger-json") != "" {
				fingerJson = viper.GetString("finger.finger-json")
			}
			if fingerJson == "" {

				currentUser, err := user.Current()
				if err != nil {
					fmt.Println("无法获取当前用户信息:", err)
					return
				}
				fingerJson = currentUser.HomeDir + "/" + ".finger.json"

				//打开文件io流
				f, err := os.Open(fingerJson)
				if err != nil {
					color.RGBStyleFromString("237,64,35").Println("[Error] finger json file error!!! 打开" + fingerJson + " 错误")
					cmd.Usage()
					os.Exit(1)
				}
				if err = f.Close(); err != nil {
					log.Fatal(err)
				}

			}
			err := finger.LoadWebfingerprint(fingerJson)
			if err != nil {
				color.RGBStyleFromString("237,64,35").Println("[Error] finger json file error!!! 打开" + fingerJson + " 错误")
				os.Exit(1)
			}

			fingers := finger.GetWebfingerprint().Fingerprint
			// 创建一个新的表格
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Product", "Fofa Query"})

			for _, f := range fingers {
				if regexMatch(product, f.Cms) {
					t.AppendRow(table.Row{
						f.Cms,
						getFofaFromFingerprint(f),
					})
				}
			}

			t.SetStyle(table.StyleLight)
			t.Render()

		} else {
			cmd.Usage()
		}
	},
}

var (
	product string
)

func getFofaFromFingerprint(f finger.Fingerprint) string {
	if f.Method == "keyword" {
		patterns = []string{}

		if f.Location == "body" {
			for _, k := range f.Keyword {
				patterns = append(patterns, "body="+strconv.Quote(k))
			}
			return strings.Join(patterns, " && ")
		} else if f.Location == "header" {
			for _, k := range f.Keyword {
				patterns = append(patterns, "header="+strconv.Quote(k))
			}
			return strings.Join(patterns, " && ")
		} else if f.Location == "title" {
			return "title=" + strconv.Quote(f.Keyword[0])
		}

	} else if f.Method == "faviconhash" {
		return "icon_hash=" + strconv.Quote(f.Keyword[0])
	}
	return ""
}

func init() {
	fingerCmd.AddCommand(fingerGetCmd)
	fingerGetCmd.Flags().StringVarP(&product, "product", "p", "", "产品名称或关键字")
}

func regexMatch(product, cms string) bool {
	product = strings.ReplaceAll(product, " ", "")
	product = strings.TrimSpace(product)
	product = strings.ToLower(product)

	cms = strings.ReplaceAll(cms, " ", "")
	cms = strings.ToLower(cms)
	reg := regexp.MustCompile(product)
	match := reg.FindStringSubmatch(cms)
	return len(match) > 0
}
