package services

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"time"
)

func ImagesToBase64(str_images string) string {
	f, _ := os.Open(str_images)
	defer f.Close()
	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
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

func ArrayReverse(s []map[string]interface{}) []map[string]interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//func ArrayReverseAny[A ~[]map[any]interface{} | ~[]interface{}](s A) A {
//	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
//		s[i], s[j] = s[j], s[i]
//	}
//	return s
//}

func ArrayReverseAny[A ~[]interface{} | ~[]map[string]interface{}](s A) A {
	v := reflect.ValueOf(s) // 获取值的反射对象
	kind := v.Kind()        // 获取反射对象的类型

	if kind != reflect.Slice { // 如果不是切片类型则返回原值
		return s
	}

	length := v.Len()                                                // 获取切片长度
	reversed := reflect.MakeSlice(reflect.TypeOf(s), length, length) // 创建一个新的切片，用于存储反转后的元素

	for i, j := 0, length-1; i < length; i, j = i+1, j-1 {
		reversed.Index(i).Set(v.Index(j)) // 使用反射设置切片中的元素
	}

	return reversed.Interface().(A) // 将反转后的切片转换回 A 类型并返回
}
