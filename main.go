package main

import (
	"flag"
	"gitee-image-hosting/routes"
	"gitee-image-hosting/services/flag_handle"
	"log"
)

func init()  {
	port := flag.String("port", "2047", "本地监听的端口")
	token := flag.String("token", "", "Gitee 的用户授权码")
	owner := flag.String("owner", "hezhizheng", "仓库所属空间地址(企业、组织或个人的地址path)")
	repo := flag.String("repo", "pictest", "仓库路径(path)")
	path := flag.String("path", "image-hosting", "文件的路径")
	flag.Parse()
	flag_handle.PORT = *port
	flag_handle.OWNER = *owner
	flag_handle.REPO = *repo
	flag_handle.PATH = *path
	flag_handle.TOKEN = *token

	if flag_handle.TOKEN == "" {
		panic("token 必须！")
	}

}

func main() {
	port := flag_handle.PORT
	router := routes.InitRoute()
	log.Println("监听端口", "http://127.0.0.1:"+port)
	err := router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}