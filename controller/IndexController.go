package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"repo-image-hosting/services/flag_handle"
)

func Index(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"current_dir": flag_handle.OWNER + "/" + flag_handle.REPO + "/" + flag_handle.PATH,
		"owner":       flag_handle.OWNER,
		"repo":        flag_handle.REPO,
		"platform":    flag_handle.PLATFORM,
	})
}
