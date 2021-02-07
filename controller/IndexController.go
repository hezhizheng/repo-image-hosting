package controller

import (
	"gitee-image-hosting/services"
	"gitee-image-hosting/services/flag_handle"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"current_dir": flag_handle.OWNER+"/"+flag_handle.REPO+"/"+flag_handle.PATH,
		"owner":flag_handle.OWNER,
		"repo":flag_handle.REPO,
	})
}

func Images(c *gin.Context)  {
	url := "https://gitee.com/api/v5/repos/"+flag_handle.OWNER+"/"+flag_handle.REPO+"/contents/"+flag_handle.PATH+"?access_token="+flag_handle.TOKEN
	images := services.GetGiteeFiles(url)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"data": map[string]interface{}{
			"images": images,
		},
	})
}