package controller

import (
	"gitee-image-hosting/services"
	"gitee-image-hosting/services/flag_handle"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"time"
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

	timeStr := time.Now().Format("20060102150405")

	file.Filename = timeStr + "_" + services.GetRandomString(16) + fileExt

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

	picUrl, picPath, picSha := services.PushGitee(PushGiteeUrl, Base64)

	//删除临时图片
	os.Remove(filename)

	if picUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "上传失败，请检查token与其他配置参数是否正确",
			"data":  map[string]interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"data": map[string]interface{}{
			"url":  picUrl,
			"sha":  picSha,
			"path": picPath,
		},
	})
}

func ImageDel(c *gin.Context) {
	sha := c.PostForm("sha")
	_path := c.PostForm("path")

	DelGiteeUrl := "https://gitee.com/api/v5/repos/" + flag_handle.OWNER + "/" + flag_handle.REPO + "/contents/" + _path

	services.DelFile(DelGiteeUrl, sha)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"data":    map[string]interface{}{},
	})
}