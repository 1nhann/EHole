package cmd

import (
	"bufio"
	"ehole/module/fofa"
	_ "embed"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

//go:embed fofa-hack_win_amd64.exe
var hackExe []byte

var fofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "fofa get targets",
	Long:  `调用fofa接口，获取url`,
	Run: func(cmd *cobra.Command, args []string) {
		if query == "" {
			fmt.Print("Please enter your query: \n")
			reader := bufio.NewReader(os.Stdin)
			q, err := reader.ReadString('\n')
			query = strings.TrimSpace(q)
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
		}
		if hack {
			tmpFile, err := ioutil.TempFile("", "fofa-hack-*.exe")
			if err != nil {
				fmt.Println("无法创建临时文件:", err)
				return
			}
			defer os.Remove(tmpFile.Name()) // 在退出时删除临时文件
			// 将嵌入的 exe 文件写入临时文件
			if _, err := tmpFile.Write(hackExe); err != nil {
				fmt.Println("无法写入临时文件:", err)
				return
			}
			if err := tmpFile.Close(); err != nil {
				fmt.Println("无法关闭临时文件:", err)
				return
			}
			// 使临时文件可执行
			if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
				fmt.Println("无法设置临时文件权限:", err)
				return
			}
			// 执行临时文件
			cmd := exec.Command(tmpFile.Name(), "--keyword", query, "--endcount", strconv.FormatInt(int64(size), 10), "--output", "txt", "--outputname", strings.Split(out, ".")[0])
			// 获取标准输出管道
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Printf("无法获取标准输出管道: %v\n", err)
				return
			}
			// 获取标准错误输出管道
			stderr, err := cmd.StderrPipe()
			if err != nil {
				fmt.Printf("无法获取标准错误输出管道: %v\n", err)
				return
			}
			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					fmt.Printf("%s\n", scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					fmt.Printf("读取标准输出错误: %v\n", err)
				}
			}()
			// 实时读取并打印标准错误输出
			go func() {
				scanner := bufio.NewScanner(stderr)
				for scanner.Scan() {
					fmt.Printf("%s\n", scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					fmt.Printf("读取标准错误输出错误: %v\n", err)
				}
			}()
			// 启动命令
			if err := cmd.Start(); err != nil {
				fmt.Printf("无法启动命令: %v\n", err)
				return
			}
			// 实时读取并打印标准输出

			// 等待命令执行完成
			if err := cmd.Wait(); err != nil {
				fmt.Printf("命令执行失败: %v\n", err)
				return
			}
			//_, err = cmd.CombinedOutput()
			//if err != nil {
			//	fmt.Println("执行命令失败:", err)
			//	return
			//}

		} else {
			if fofaIni == "" {
				currentUser, err := user.Current()
				if err != nil {
					fmt.Println("无法获取当前用户信息:", err)
					return
				}
				fofaIni = currentUser.HomeDir + "/" + ".fofa.ini"
				//打开文件io流
				f, err := os.Open(fofaIni)
				if err != nil {
					file, err := os.Create(fofaIni)
					if err != nil {
						fmt.Println("无法创建文件:", err)
						return
					}

					// 写入字符串
					content := `Email=xxxxxxxxx@admin.com
Fofa_token=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Fofa_timeout=10`
					_, err = file.WriteString(content)
					if err != nil {
						fmt.Println("无法写入文件:", err)
						return
					}
					file.Close()
					color.RGBStyleFromString("237,64,35").Println("[Error] Fofa configuration file error!!!")
					cmd.Usage()
					os.Exit(1)
				}
				if err = f.Close(); err != nil {
					log.Fatal(err)
				}

			}

			targets := fofa.FofaSearch(query, size, page, fofaIni)
			file, err := os.Create(out)
			if err != nil {
				fmt.Println("无法创建文件:", err)
				return
			}
			defer file.Close()

			// 写入字符串
			content := strings.Join(targets, "\n")
			_, err = file.WriteString(content)
			if err != nil {
				fmt.Println("无法写入文件:", err)
				return
			}
		}
	},
}
var (
	query   string
	page    int
	size    int
	fofaIni string
	out     string
	hack    bool
)

func init() {
	rootCmd.AddCommand(fofaCmd)
	fofaCmd.Flags().StringVarP(&query, "query", "q", "", "fofa query")
	fofaCmd.Flags().IntVarP(&page, "page", "p", 1, "page")
	fofaCmd.Flags().IntVarP(&size, "size", "s", 20, "size per page")
	fofaCmd.Flags().StringVarP(&fofaIni, "fofa-ini", "i", "", "默认是$HOME/.fofa.ini放置email和token")
	fofaCmd.Flags().StringVarP(&out, "out", "o", "targets.txt", "默认是targets.txt")
	fofaCmd.Flags().BoolVarP(&hack, "hack", "x", false, "默认是false，如果是true的话就调用fofa-hack.exe:https://github.com/Cl0udG0d/Fofa-hack")
}
