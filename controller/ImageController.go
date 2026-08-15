package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"repo-image-hosting/services"
	"repo-image-hosting/services/connector"
	"strconv"
	"sync"
	"time"
)

var mutex = sync.Mutex{}

var (
	fileCacheLock sync.Mutex
	fileCache     []map[string]interface{}
	fileCacheAt   time.Time
	fileCacheAge  = 60 * time.Second
)

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
	invalidateFileCache()

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
	invalidateFileCache()

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"data":    map[string]interface{}{},
	})
}

// 缓存文件列表 避免每次刷新页面都请求一次远端仓库接口
func getFilesCached() ([]map[string]interface{}, error) {
	fileCacheLock.Lock()
	if fileCache != nil && time.Since(fileCacheAt) < fileCacheAge {
		files := fileCache
		fileCacheLock.Unlock()
		return files, nil
	}
	fileCacheLock.Unlock()

	// 网络请求在锁外执行 避免一个慢请求阻塞其他所有请求
	files, err := safeGetFiles()
	if err != nil {
		// 请求失败不缓存 避免把错误结果缓存起来
		return nil, err
	}

	fileCacheLock.Lock()
	defer fileCacheLock.Unlock()
	if files != nil {
		fileCache = files
		fileCacheAt = time.Now()
	}
	return files, nil
}

// 兜底 recover: 远端接口数据异常导致的 panic 不拖垮整个程序
func safeGetFiles() (files []map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("获取文件列表异常: %v", r)
			log.Println(err)
		}
	}()
	files, err = connector.RepoCreate().GetFiles()
	return
}

func invalidateFileCache() {
	fileCacheLock.Lock()
	defer fileCacheLock.Unlock()
	fileCache = nil
}

func Images(c *gin.Context) {
	images, err := getFilesCached()
	if err != nil {
		log.Println("Images 获取文件列表失败:", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"success": false,
			"message": err.Error(),
			"data": map[string]interface{}{
				"images": nil,
			},
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 30
	}
	if pageSize > 100 {
		pageSize = 100
	}

	total := len(images)

	var paged []map[string]interface{}
	end := 0
	if images != nil {
		start := (page - 1) * pageSize
		end = start + pageSize
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		paged = images[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"success": true,
		"data": map[string]interface{}{
			"images":   paged,
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"hasMore":  end < total,
		},
	})
}
