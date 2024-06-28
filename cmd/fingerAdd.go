package cmd

import (
	"bytes"
	"ehole/module/finger"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"
)

var fingerAddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		if productAdd != "" && len(patterns) > 0 {
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
					color.RGBStyleFromString("237,64,35").Println("[Error] finger json file error!!! 打开 " + fingerJson + " 错误")
					cmd.Usage()
					os.Exit(1)
				}
				if err = f.Close(); err != nil {
					log.Fatal(err)
				}

			}
			//for _, f := range fingers {
			//	if regexMatch(product, f.Cms) {
			//		//result = append(result, f)
			//		b, _ := json.Marshal(f)
			//		fmt.Println(string(b))
			//	}
			//}
			var method string
			var location string
			var keyword []string
			f := finger.Fingerprint{
				Cms: productAdd,
			}
			if fingerType == "word" {
				method = "keyword"
				if part != "body" && part != "header" {
					cmd.Usage()
					return
				}
				keyword = patterns
				location = part
			} else if fingerType == "hash" {
				method = "faviconhash"
				keyword = patterns
				location = "body"
			} else if fingerType == "title" {
				method = "keyword"
				if len(patterns) != 1 {
					cmd.Usage()
					return
				}
				location = "title"
				title := patterns[0]
				keyword = []string{title}
			}
			f.Method = method
			f.Keyword = keyword
			f.Location = location

			err := finger.LoadWebfingerprint(fingerJson)
			if err != nil {
				color.RGBStyleFromString("237,64,35").Println("[error] fingerprint file error!!! parse json error !!!")
				os.Exit(1)
			}
			packJson := finger.GetWebfingerprint()
			fingers := packJson.Fingerprint
			fingers = append(fingers, f)

			packJson.Fingerprint = fingers

			bf := bytes.NewBuffer([]byte{})
			jsonEncoder := json.NewEncoder(bf)
			jsonEncoder.SetEscapeHTML(false)
			jsonEncoder.SetIndent("", "    ")
			err = jsonEncoder.Encode(packJson)
			if err != nil {
				fmt.Println("JSON 编码错误:", err)
				return
			}
			file, err := os.Create(fingerJson)
			if err != nil {
				fmt.Println("无法创建文件:", err)
				return
			}
			defer file.Close()

			// 写入字符串
			_, err = file.WriteString(string(bytes.TrimSpace(bf.Bytes())))
			if err != nil {
				fmt.Println("无法写入文件:", err)
				return
			}

			bf2 := bytes.NewBuffer([]byte{})
			jsonEncoder = json.NewEncoder(bf2)
			jsonEncoder.SetEscapeHTML(false)
			jsonEncoder.SetIndent("", "    ")
			err = jsonEncoder.Encode(f)
			fmt.Println(string(bytes.TrimSpace(bf2.Bytes())))
		} else {
			cmd.Usage()
		}
	},
}
var (
	fingerType string
	productAdd string
	part       string
	patterns   []string
)

func init() {
	fingerCmd.AddCommand(fingerAddCmd)
	fingerAddCmd.Flags().StringVarP(&fingerType, "type", "t", "word", "finger type,可以是word,hash,title")
	fingerAddCmd.Flags().StringVarP(&productAdd, "product", "p", "", "产品名称或关键字")
	fingerAddCmd.Flags().StringVarP(&part, "part", "P", "body", "finger part,要匹配的位置,可以是body,header")
	fingerAddCmd.Flags().StringSliceVarP(&patterns, "patterns", "k", []string{}, "finger patterns,是一个数组,表示and关系")
}
