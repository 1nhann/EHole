package cmd

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

var fingerOpenCmd = &cobra.Command{
	Use: "open",
	Run: func(cmd *cobra.Command, args []string) {
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
		openFile(fingerJson)
	},
}

func init() {
	fingerCmd.AddCommand(fingerOpenCmd)
}
func openFile(filename string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filename)
	case "darwin":
		cmd = exec.Command("open", filename)
	case "linux":
		cmd = exec.Command("xdg-open", filename)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
