package connector

import (
	"hezhiheng/repo-image-hosting/services"
	"hezhiheng/repo-image-hosting/services/flag_handle"
	"hezhiheng/repo-image-hosting/services/github"
)

// 定义 serve 的映射关系
var serveMap = map[string]services.RepoInterface{
	"gitee":   &services.GiteeServe{},
	"github": &github.GithubServe{},
}

func RepoCreate() services.RepoInterface {
	return serveMap[flag_handle.PLATFORM]
}