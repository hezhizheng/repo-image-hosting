package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"repo-image-hosting/services"
	"repo-image-hosting/services/connector"
	"sync"
	"time"
)

var mutex = sync.Mutex{}

func ImgUpload(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()
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

	picUrl, picPath, picSha := connector.RepoCreate().Push(file.Filename, Base64)

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

	connector.RepoCreate().Del(_path, sha)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"data":    map[string]interface{}{},
	})
}
