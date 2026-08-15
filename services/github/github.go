package github

import (
	"encoding/json"
	"fmt"
	"log"
	"repo-image-hosting/services"
	"repo-image-hosting/services/flag_handle"
	"time"

	"github.com/valyala/fasthttp"
)

type GithubServe struct {
	services.RepoInterface
}

func (serve *GithubServe) Push(filename, content string) (string, string, string) {
	return Push(filename, content)
}

func (serve *GithubServe) GetFiles() ([]map[string]interface{}, error) {
	return GetFiles()
}

func (serve *GithubServe) Del(filepath, sha string) string {
	return DelFile(filepath, sha)
}

func Push(filename, content string) (string, string, string) {

	url := "https://api.github.com/repos/" +
		flag_handle.OWNER + "/" +
		flag_handle.REPO +
		"/contents/" +
		flag_handle.PATH +
		"/" + filename

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("PUT")
	req.Header.SetBytesKV([]byte("Content-Type"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+flag_handle.TOKEN))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]string)
	args["content"] = content
	args["message"] = "upload pic for repo-image-hosting"
	args["branch"] = flag_handle.BRANCH

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	// 发起请求 设置超时时间 避免请求一直挂起导致页面"数据加载中"卡死
	if err := fasthttp.DoTimeout(req, resp, 15*time.Second); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	var mapResult map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	d := ""
	p := ""
	s := ""

	_, ok := mapResult["content"]

	if ok {
		if mapResult["content"] != nil {
			path := mapResult["content"].(map[string]interface{})["path"].(string)
			d = "https://cdn.jsdmirror.com/gh/" + flag_handle.OWNER + "/" + flag_handle.REPO + "@" + flag_handle.BRANCH + "/" + path
			p = path
			s = mapResult["content"].(map[string]interface{})["sha"].(string)
		}
	}

	return d, p, s
}

func GetFiles() ([]map[string]interface{}, error) {
	// 初始化请求与响应
	url := "https://api.github.com/repos/" +
		flag_handle.OWNER + "/" +
		flag_handle.REPO +
		"/contents/" +
		flag_handle.PATH +
		"?ref=" + flag_handle.BRANCH

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("GET")
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+flag_handle.TOKEN))
	// 设置请求的目标网址
	req.SetRequestURI(url)

	// 发起请求 设置超时时间 避免请求一直挂起
	if err := fasthttp.DoTimeout(req, resp, 15*time.Second); err != nil {
		log.Println(" 请求失败:", url, err.Error())
		return nil, fmt.Errorf("github 接口请求失败: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("github 接口返回异常状态码: %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	// 获取响应的数据实体
	body := resp.Body()

	//log.Println(string(body),url)

	var mapResult []map[string]interface{}

	if err := json.Unmarshal(body, &mapResult); err != nil {
		fmt.Println("JsonToMapDemo err: ", err, string(body))
		return nil, fmt.Errorf("github 返回数据解析失败: %v", err)
	}

	for _, v := range mapResult {
		filePath, ok := v["path"].(string)
		if !ok {
			continue
		}
		v["download_url"] = "https://cdn.jsdmirror.com/gh/" + flag_handle.OWNER + "/" + flag_handle.REPO + "@" + flag_handle.BRANCH + "/" + filePath
	}

	mapResult = services.ArrayReverseAny(mapResult)

	return mapResult, nil
}

func DelFile(filepath, sha string) string {
	// 初始化请求与响应
	url := "https://api.github.com/repos/" +
		flag_handle.OWNER + "/" +
		flag_handle.REPO +
		"/contents/" + filepath

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	// 设置请求方法
	req.Header.SetMethod("DELETE")
	req.Header.SetBytesKV([]byte("Content-Type"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+flag_handle.TOKEN))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]string)
	args["sha"] = sha
	args["message"] = "delete pic for repo-image-hosting"
	args["branch"] = flag_handle.BRANCH

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	// 发起请求 设置超时时间 避免请求一直挂起
	if err := fasthttp.DoTimeout(req, resp, 15*time.Second); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	return string(body)
}
