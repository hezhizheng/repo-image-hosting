package coding

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"log"
	"os"
	"repo-image-hosting/services"
	"repo-image-hosting/services/flag_handle"
	"strconv"
	"sync"
)

var mutex = sync.Mutex{}
var mutex1 = sync.Mutex{}
var mutex2 = sync.Mutex{}
var mutex3 = sync.Mutex{}

type CodingServe struct {
	services.RepoInterface
}

func (serve *CodingServe) Push(filename, content string) (string, string, string) {
	return Push(filename, content)
}

func (serve *CodingServe) GetFiles() []map[string]interface{} {
	return GetFiles()
}

func (serve *CodingServe) Del(filepath, sha string) string {
	return DelFile(filepath, sha)
}

func Push(filename, content string) (string, string, string) {
	// 加锁 ！！！
	mutex.Lock()
	defer mutex.Unlock()
	LastCommitSha := lastCommitSha()

	// log.Println("LastCommitSha", LastCommitSha)

	url := "https://e.coding.net/open-api"

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
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+flag_handle.TOKEN))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]interface{})
	args["Action"] = "CreateBinaryFiles"
	args["DepotId"] = flag_handle.DEPOTID
	args["UserId"] = getUserId()
	args["SrcRef"] = flag_handle.BRANCH
	args["DestRef"] = flag_handle.BRANCH
	args["Message"] = "upload pic for repo-image-hosting"
	args["LastCommitSha"] = LastCommitSha
	args["GitFiles"] = []map[string]interface{}{
		{
			"Path":    flag_handle.PATH + "/" + filename,
			"Content": content,
			"NewPath": "",
		},
	}

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

	// log.Println("Push_mapResult", mapResult)

	picUrl := ""
	picPath := ""
	picSha := ""

	// todo 这里存在还没提交就获取了文件的情况？
	files := GetFiles()
	flen := len(files)

	if flen > 0 { // 取最新一条
		last := files[flen-1]
		picUrl = last["download_url"].(string)
		picPath = last["Path"].(string)
		picSha = last["Sha"].(string)
	}

	return picUrl, picPath, picSha

}

func GetFiles() []map[string]interface{} {
	mutex1.Lock()
	defer mutex1.Unlock()

	// 初始化请求与响应
	url := "https://e.coding.net/open-api"

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
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+flag_handle.TOKEN))
	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]interface{})
	args["Action"] = "DescribeGitFiles"
	args["DepotId"] = flag_handle.DEPOTID
	args["Ref"] = flag_handle.BRANCH
	args["Path"] = flag_handle.PATH

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	// 发起请求
	if err := fasthttp.Do(req, resp); err != nil {
		log.Println(" 请求失败:", url, err.Error())
	}

	// 获取响应的数据实体
	body := resp.Body()

	//log.Println(string(body),url)

	var mapResult map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	items := mapResult["Response"].(map[string]interface{})["Items"].([]interface{})

	for _, v := range items { //hzz333
		v.(map[string]interface{})["download_url"] = "https://" + flag_handle.OWNER + ".coding.net" + "/p/" + flag_handle.PROJECT + "/d/" + flag_handle.REPO + "/git/raw/" + flag_handle.BRANCH + "/" + v.(map[string]interface{})["Path"].(string)
		v.(map[string]interface{})["name"] = v.(map[string]interface{})["Name"]
		v.(map[string]interface{})["sha"] = v.(map[string]interface{})["Sha"]
	}

	var mapResult2 []map[string]interface{}

	c, _ := json.Marshal(items)

	json.Unmarshal(c, &mapResult2)

	return mapResult2
}

func DelFile(filepath, sha string) string {
	// coding 不支持
	return ""
}

func getUserId() int {
	mutex2.Lock()
	defer mutex2.Unlock()

	cacheUId := cacheGetUserId()
	if cacheUId > 0 {
		return cacheUId
	}

	url := "https://e.coding.net/open-api"

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
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+flag_handle.TOKEN))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]interface{})
	args["Action"] = "DescribeCodingCurrentUser"

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

	id := mapResult["Response"].(map[string]interface{})["User"].(map[string]interface{})["Id"].(float64)

	stringId := strconv.FormatFloat(id, 'f', -1, 64)

	cacheUserId(stringId)

	// string转成int：
	intId, _ := strconv.Atoi(stringId)
	return intId
}

func lastCommitSha() string {
	mutex3.Lock()
	defer mutex3.Unlock()
	url := "https://e.coding.net/open-api"

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
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+flag_handle.TOKEN))

	// 设置请求的目标网址
	req.SetRequestURI(url)

	args := make(map[string]interface{})
	args["Action"] = "DescribeGitCommits"
	args["DepotId"] = flag_handle.DEPOTID
	args["Ref"] = flag_handle.BRANCH
	args["PageNumber"] = 1
	args["PageSize"] = 1

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
	Sha := mapResult["Response"].(map[string]interface{})["Commits"].([]interface{})[0].(map[string]interface{})["Sha"].(string)
	return Sha
}

func cacheUserId(str string) bool {
	f, err := os.Create("./" + _md5(flag_handle.TOKEN) + ".cache")
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		_, err2 := f.Write([]byte(str))
		if err2 != nil {
			return false
		}
		return true
	}
}

func cacheGetUserId() int {
	f, err := os.OpenFile("./"+_md5(flag_handle.TOKEN)+".cache", os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return 0
	} else {
		contentByte, err2 := ioutil.ReadAll(f)
		if err2 != nil {
			fmt.Println(err2.Error())
			return 0
		}
		intId, _ := strconv.Atoi(string(contentByte))
		return intId
	}
}

func _md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}
