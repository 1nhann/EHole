/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"ehole/module/finger"
	"ehole/module/finger/source"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"
)

// fingerCmd represents the finger command
var fingerCmd = &cobra.Command{
	Use:   "finger",
	Short: "ehole的指纹识别模块",
	Long:  `从fofa或者本地文件获取资产进行指纹识别，支持单条url识别。`,
	Run: func(cmd *cobra.Command, args []string) {
		color.RGBStyleFromString("105,187,92").Println("\n     ______    __         ______                 \n" +
			"    / ____/___/ /___ ____/_  __/__  ____ _____ ___ \n" +
			"   / __/ / __  / __ `/ _ \\/ / / _ \\/ __ `/ __ `__ \\\n" +
			"  / /___/ /_/ / /_/ /  __/ / /  __/ /_/ / / / / / /\n" +
			" /_____/\\__,_/\\__, /\\___/_/  \\___/\\__,_/_/ /_/ /_/ \n" +
			"			 /____/ https://forum.ywhack.com  By:shihuang\n")
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

		if localfile != "" {
			urls := removeRepeatedElement(source.LocalFile(localfile))
			s := finger.NewScan(urls, thread, output, proxy, fingerJson, spider)
			s.StartScan(spider)
			os.Exit(1)
		}

		if urla != "" {
			s := finger.NewScan([]string{urla}, thread, output, proxy, fingerJson, spider)
			s.StartScan(spider)
			os.Exit(1)
		}

	},
}

var (
	localfile  string
	urla       string
	thread     int
	output     string
	proxy      string
	fingerJson string
	spider     bool
)

func init() {
	rootCmd.AddCommand(fingerCmd)
	fingerCmd.Flags().StringVarP(&localfile, "local", "l", "", "从本地文件读取资产，进行指纹识别，支持无协议，列如：192.168.1.1:9090 | http://192.168.1.1:9090")
	fingerCmd.Flags().StringVarP(&urla, "url", "u", "", "识别单个目标。")
	fingerCmd.Flags().StringVarP(&output, "output", "o", "", "输出所有结果，当前仅支持json和xlsx后缀的文件。")
	fingerCmd.Flags().IntVarP(&thread, "thread", "t", 100, "指纹识别线程大小。")
	fingerCmd.Flags().StringVarP(&proxy, "proxy", "p", "", "指定访问目标时的代理，支持http代理和socks5，例如：http://127.0.0.1:8080、socks5://127.0.0.1:8080")
	fingerCmd.Flags().StringVarP(&fingerJson, "finger-json", "j", "", "默认是$HOME/.finger.json")
	fingerCmd.Flags().BoolVarP(&spider, "spider", "s", false, "默认是false，默认关闭爬虫")
}

func removeRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
