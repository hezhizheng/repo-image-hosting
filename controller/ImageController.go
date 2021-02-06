package controller

import (
	"gitee-image-hosting/services"
	"gitee-image-hosting/services/flag_handle"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func ImgUpload(c *gin.Context) {
	file, err := c.FormFile("qqfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
			"data":  map[string]interface{}{},
		})
		return
	}

	filepath := "./"

	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	fileExt := path.Ext(filepath + file.Filename)
	file.Filename = services.GetRandomString(16) + fileExt

	filename := filepath + file.Filename

	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
			"data":  map[string]interface{}{},
		})
		return
	}

	Base64 := services.ImagesToBase64(filename)

	PushGiteeUrl := "https://gitee.com/api/v5/repos/" + flag_handle.OWNER + "/" + flag_handle.REPO + "/contents/" + flag_handle.PATH + "/" + file.Filename

	picUrl := services.PushGitee(PushGiteeUrl, Base64)

	//删除临时图片
	os.Remove(filename)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"data": map[string]interface{}{
			"url": picUrl,
		},
	})
}