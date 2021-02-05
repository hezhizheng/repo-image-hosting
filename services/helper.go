package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitee-image-hosting/services/flag_handle"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"os"
	"time"
)

func ImagesToBase64(str_images string) []byte {
	//读原图片
	ff, _ := os.Open(str_images)
	defer ff.Close()
	sourcebuffer := make([]byte, 500000)
	n, _ := ff.Read(sourcebuffer)
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	return []byte(sourcestring)
}

func GetRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func PushGitee(url, content string) string {
	// 初始化请求与响应
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("POST")
	req.Header.SetBytesKV([]byte("Content-Type"), []byte("application/json"))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]string)
	args["access_token"] = flag_handle.TOKEN
	args["content"] = content
	args["message"] = "upload pic for gitee-image-hosting"

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	// 发起请求
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	var mapResult map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	s := mapResult["content"].(map[string]interface{})["download_url"].(string)

	return s
}

func GetGiteeFiles(url string) []map[string]interface{} {
	// 初始化请求与响应
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("GET")
	// 设置请求的目标网址
	req.SetRequestURI(url)

	// 发起请求
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	var mapResult []map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	return mapResult
}
