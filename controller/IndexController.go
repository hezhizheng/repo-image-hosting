package controller

import (
	"repo-image-hosting/services/connector"
	"repo-image-hosting/services/flag_handle"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"current_dir": flag_handle.OWNER+"/"+flag_handle.REPO+"/"+flag_handle.PATH,
		"owner":flag_handle.OWNER,
		"repo":flag_handle.REPO,
		"platform":flag_handle.PLATFORM,
	})
}

func Images(c *gin.Context)  {
	images := connector.RepoCreate().GetFiles()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"data": map[string]interface{}{
			"images": images,
		},
	})
}