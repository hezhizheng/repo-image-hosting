package main

import (
	"flag"
	"os/exec"
	"repo-image-hosting/routes"
	"repo-image-hosting/services/flag_handle"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
)

func init()  {
	port := flag.String("port", "2047", "本地监听的端口")
	platform := flag.String("platform", "github", "平台名称，支持gitee/github")
	token := flag.String("token", "", "Gitee/Github 的用户授权码")
	owner := flag.String("owner", "hezhizheng", "仓库所属空间地址(企业、组织或个人的地址path)")
	repo := flag.String("repo", "static-image-hosting", "仓库路径(path)")
	path := flag.String("path", "image-hosting", "文件的路径")
	flag.Parse()
	flag_handle.PORT = *port
	flag_handle.OWNER = *owner
	flag_handle.REPO = *repo
	flag_handle.PATH = *path
	flag_handle.TOKEN = *token
	flag_handle.PLATFORM = *platform

	if flag_handle.TOKEN == "" {
		panic("token 必须！")
	}

}

func main() {
	// 关闭debug模式
	gin.SetMode(gin.ReleaseMode)

	port := flag_handle.PORT
	router := routes.InitRoute()

	link := "http://127.0.0.1:"+port
	log.Println("监听端口", link ," 请不要关闭终端")
	// 调用浏览器打开网页
	if runtime.GOOS == "windows" {
		exec.Command("cmd","/c","start", link).CombinedOutput()
	} else {
		// 没有测试环境，暂时不实现
		//exec.Command("open", link).CombinedOutput() // macos
		//exec.Command("x-www-browser", link).CombinedOutput() // linux
	}

	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}